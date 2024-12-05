package bootstrap

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/guackamolly/zero-monitor/internal/config"
	"github.com/guackamolly/zero-monitor/internal/env"
	"github.com/guackamolly/zero-monitor/internal/http"
)

func Node() env.NodeEnv {
	// Wait for user input regarding the network invite link.
	var inviteLink string
	println("Waiting for invite link... (press enter to resume)")
	fmt.Scanln(&inviteLink)

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
