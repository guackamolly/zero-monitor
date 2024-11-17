package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"

	"github.com/guackamolly/zero-monitor/internal/config"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

func main() {
	logging.AddLogger(logging.NewConsoleLogger())

	configPath, err := config.Dir()
	if err != nil {
		logging.LogFatal("couldn't lookup config dir, %v", err)
	}

	// Generate a RSA private key with 2048 bits
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logging.LogFatal("couldn't generate private key, %v", err)
	}

	// Create a PEM-encoded block for the private key
	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	}

	// Write the private key to a file
	privFile, err := os.Create(filepath.Join(configPath, "mq.pem"))
	if err != nil {
		logging.LogFatal("couldn't create private key file, %v", err)
	}
	defer privFile.Close()
	err = pem.Encode(privFile, &privBlock)
	if err != nil {
		logging.LogFatal("couldn't write private key to file, %v", err)
	}

	// Extract the public key from the private key
	publicKey := &privateKey.PublicKey

	// Create a PEM-encoded block for the public key
	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		logging.LogFatal("couldn't extract public key from private key, %v", err)
	}
	pubBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubBytes,
	}

	// Write the public key to a file
	pubFile, err := os.Create(filepath.Join(configPath, "mq.pub"))
	if err != nil {
		logging.LogFatal("couldn't create public key file, %v", err)
	}
	defer pubFile.Close()
	err = pem.Encode(pubFile, &pubBlock)
	if err != nil {
		logging.LogFatal("couldn't write public key to file, %v", err)
	}

	logging.LogInfo("generated private and public keys for encrypting message queue transport")
	logging.LogInfo(privFile.Name())
	logging.LogInfo(pubFile.Name())
}
