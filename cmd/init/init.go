package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"

	"github.com/guackamolly/zero-monitor/internal/config"
	_http "github.com/guackamolly/zero-monitor/internal/http"
	"github.com/joho/godotenv"
)

const (
	BootstrapMaster Action = iota
	BootstrapNode
)

type Action int

type NodeEnv struct {
	MessageQueueHost            string `env:"mq_sub_host"`
	MessageQueuePort            int    `env:"mq_sub_port"`
	MessageQueueTransportPubKey string `env:"mq_transport_pub_key"`
	MessageQueueInviteCode      string `env:"mq_invite_code"`
}

type MasterEnv struct {
	ServerHost                  string `env:"server_host"`
	ServerPort                  int    `env:"server_port"`
	MessageQueueHost            string `env:"mq_sub_host"`
	MessageQueuePort            int    `env:"mq_sub_port"`
	MessageQueueTransportPubKey string `env:"mq_transport_pub_key"`
	MessageQueueTransportPemKey string `env:"mq_transport_pem_key"`
	BoltDBPath                  string `env:"bolt_db_path"`
}

var action = BootstrapMaster

var inviteLink *url.URL

func init() {
	flag.Func("node", "configures the environment for starting a node", func(s string) error {
		if len(s) == 0 {
			println("Waiting for invite link... (press enter to resume)")
			fmt.Scanln(&s)
		}

		url, err := url.Parse(s)
		if err != nil {
			return err
		}

		inviteLink = url
		action = BootstrapNode
		return nil
	})

	flag.Parse()
}

func main() {
	switch action {
	case BootstrapMaster:
		bootstrapMaster()
	case BootstrapNode:
		bootstrapNode()
	}
}

func bootstrapMaster() {
	println("Bootstrapping master configuration")
	configPath := must(config.Dir())

	// Generate a RSA private key with 2048 bits
	privateKey := must(rsa.GenerateKey(rand.Reader, 2048))

	// Create a PEM-encoded block for the private key
	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	}

	// Extract the public key from the private key
	publicKey := &privateKey.PublicKey

	// Create a PEM-encoded block for the public key
	pubBytes := must(x509.MarshalPKIXPublicKey(publicKey))
	pubBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubBytes,
	}

	// Write the public key to a file
	pubFile := must(os.Create(filepath.Join(configPath, "master.pub")))
	defer pubFile.Close()
	must0(pem.Encode(pubFile, &pubBlock))

	// Write the private key to a file
	privFile := must(os.Create(filepath.Join(configPath, "master.pem")))
	defer privFile.Close()
	must0(pem.Encode(privFile, &privBlock))

	boltDbFile := must(os.Create(filepath.Join(configPath, "master.db")))
	defer boltDbFile.Close()

	env := MasterEnv{
		ServerHost:                  "0.0.0.0",
		MessageQueueHost:            "0.0.0.0",
		ServerPort:                  8080,
		MessageQueuePort:            36113,
		MessageQueueTransportPubKey: pubFile.Name(),
		MessageQueueTransportPemKey: privFile.Name(),
		BoltDBPath:                  boltDbFile.Name(),
	}

	envpath := fmt.Sprintf("%s/master.env", must(config.Dir()))
	writeEnv(env, envpath)

	println("> Generated private key on: %s", privFile.Name())
	println("> Generated public key on: %s", pubFile.Name())
	println("> Generated .env on: %s", envpath)
}

func bootstrapNode() {
	println("Bootstrapping node configuration using invite link: %s", inviteLink)
	configPath := must(config.Dir())
	inviteCode := inviteLink.Query().Get("join")

	v := downloadUnmarshal[_http.NetworkJoinView](inviteLink.String())
	pubKey := string(download(v.PublicKeyURL))
	endpoint := downloadUnmarshal[_http.NetworkConnectionEndpointView](v.ConnectionEndpointURL)

	pubFile := must(os.Create(filepath.Join(configPath, "node.pub")))
	defer pubFile.Close()
	io.WriteString(pubFile, pubKey)

	env := NodeEnv{
		MessageQueueHost:            endpoint.Host,
		MessageQueuePort:            endpoint.Port,
		MessageQueueTransportPubKey: pubFile.Name(),
		MessageQueueInviteCode:      inviteCode,
	}

	envpath := fmt.Sprintf("%s/node.env", configPath)
	writeEnv(env, envpath)

	println("> Extracted invite code: %s", inviteCode)
	println("> Saved public key on: %s", pubFile.Name())
	println("> Generated .env on: %s", envpath)
}

// helper functions
func writeEnv(env any, path string) {
	v := reflect.ValueOf(env)
	t := v.Type()
	m := map[string]string{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		m[f.Tag.Get("env")] = fmt.Sprintf("%v", v.Field(i))
	}

	godotenv.Write(m, path)
}

func download(url string) []byte {
	println("> GET %s", url)

	resp := must(http.Get(url))
	if sc := resp.StatusCode; sc != 200 {
		panic(fmt.Sprintf("sc: %d", resp.StatusCode))
	}

	bs := must(io.ReadAll(resp.Body))

	return bs
}

func downloadUnmarshal[T any](url string) T {
	var v T
	must0(json.Unmarshal(download(url), &v))

	return v
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

func must0(err error) {
	if err != nil {
		panic(err)
	}
}

func panic(v any) {
	log.Fatal(v)
}

func println(f any, v ...any) {
	if _, ok := f.(string); !ok || len(v) == 0 {
		fmt.Println(f)
		return
	}

	fmt.Printf("%s\n", fmt.Sprintf(f.(string), v))
}
