package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/gotd/td/tg"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
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
	phoneNumber := os.Getenv("PHONE_NUMBER")
	client, err := gotgproto.NewClient(
		// Get AppID from https://my.telegram.org/apps
		apiId,
		// Get ApiHash from https://my.telegram.org/apps
		apiHash,
		// ClientType, as we defined above
		gotgproto.ClientTypePhone(phoneNumber),
		// Optional parameters of client
		&gotgproto.ClientOpts{
			Session: sessionMaker.SqlSession(sqlite.Open("echobot")),
		},
	)


	client.Run(context.Background(), func(ctx context.Context) error {
		contacts := []tg.InputPhoneContact{
			{
				ClientID:  1,
				FirstName: "",
				LastName:  "",
				Phone:     "77085690946",
			},
		}

		imported, err := client.API().ContactsImportContacts(
			ctx,
			contacts,
		)
		if err != nil {
			return err
		}

		if len(imported.Users) > 0 {
			user := imported.Users[0].(*tg.User)
			log.Printf("Chat ID (User ID) for phone %s: %d", "+1234567890", user.ID)
		}

    return nil
	})
	if err != nil {
		log.Fatalln("failed to start client:", err)
	}
	client.Idle()
}
