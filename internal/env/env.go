package env

import (
	"log"
	"os"
	"path/filepath"

	"github.com/guackamolly/zero-monitor/internal/config"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err == nil {
		return
	}

	d, err := config.Dir()
	if err != nil {
		log.Fatal("couldn't lookup a .env file!")
	}

	exe, err := os.Executable()
	if err != nil {
		log.Fatal("couldn't lookup a .env file!")
	}

	envp := filepath.Join(d, filepath.Base(exe)+".env")
	if err = godotenv.Load(envp); err == nil {
		return
	}

	log.Fatalf("couldn't load .env or %s", envp)
}
