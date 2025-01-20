package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
  "context"
	"github.com/gotd/td/telegram"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiId, err := strconv.Atoi(os.Getenv("APP_ID"))
	if err != nil {
		panic(err)
	}
	apiHash := os.Getenv("API_HASH")
	sessionStorage := &telegram.FileSessionStorage{
		Path: filepath.Join(sessionDir, "session.json"),
	}
	client := telegram.NewClient(int(apiId), apiHash, telegram.Options{
		SessionStorage: sessionStorage,
	})

  client.Run(ctx context.Context,func {})
}
