package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/gotd/td/telegram"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiId, err := strconv.Atoi(os.Getenv("API_ID"))
	if err != nil {
		panic(err)
	}
	apiHash := os.Getenv("API_HASH")
	
	client := telegram.NewClient(int(apiId), apiHash, telegram.Options{
	})

	client.Run(context.Background(), func(ctx context.Context) error {
    		api := client.API()
        println(api)
		return nil
	})
}
