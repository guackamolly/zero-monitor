package bootstrap

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/guackamolly/zero-monitor/internal/config"
	"github.com/guackamolly/zero-monitor/internal/env"
	"github.com/guackamolly/zero-monitor/internal/http"
)

const (
	defaultHttpPort         = "8080"
	defaultMessageQueuePort = "36113"

	defaultHttpHost         = "0.0.0.0"
	defaultMessageQueueHost = "0.0.0.0"
)

func Master() env.MasterEnv {
	configPath := must(config.Dir())

	// Generate a RSA private key with 2048 bits
	privateKey := must(rsa.GenerateKey(rand.Reader, 2048))

	// Create a PEM-encoded block for the private key
	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	}

	// Write the private key to a file
	privFile := must(os.Create(filepath.Join(configPath, "master.pem")))
	defer privFile.Close()
	must0(pem.Encode(privFile, &privBlock))

	boltDbFile := must(os.Create(filepath.Join(configPath, "master.db")))
	defer boltDbFile.Close()

	e := env.MasterEnv{
		ServerHost:                  defaultHttpHost,
		MessageQueueHost:            defaultMessageQueueHost,
		ServerPort:                  defaultHttpPort,
		MessageQueuePort:            defaultMessageQueuePort,
		MessageQueueTransportPemKey: privFile.Name(),
		BoltDBPath:                  boltDbFile.Name(),
	}

	// Save bootstrapped env.
	return must(e, env.Save(e))
}

// Bootstrapps node. If [inviteLink] is empty, it reads invite link from stdin.
func Node(inviteLink string) env.NodeEnv {
	if len(inviteLink) == 0 {
		// Wait for user input regarding the network invite link.
		println("Waiting for invite link... (press enter to resume)")
		fmt.Scanln(&inviteLink)
	}

	configPath := must(config.Dir())
	inviteCode := must(url.Parse(inviteLink)).Query().Get("join")

	// Query connection information and public key for key exchange.
	joinView := downloadUnmarshal[http.NetworkJoinView](inviteLink)
	pubKey := string(download(joinView.PublicKeyURL))
	endpoint := downloadUnmarshal[http.NetworkConnectionEndpointView](joinView.ConnectionEndpointURL)

	// Write pub key to config folder.
	pubFile := must(os.Create(filepath.Join(configPath, "node.pub")))
	defer pubFile.Close()
	io.WriteString(pubFile, pubKey)

	e := env.NodeEnv{
		MessageQueueHost:            endpoint.Host,
		MessageQueuePort:            fmt.Sprintf("%d", endpoint.Port),
		MessageQueueTransportPubKey: pubFile.Name(),
		MessageQueueInviteCode:      inviteCode,
	}

	// Save bootstrapped env.
	return must(e, env.Save(e))
}
