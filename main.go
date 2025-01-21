package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
)

type Message struct {
	PhoneNumber string `json:"phoneNumber"`
	MessageText string `json:"message"`
}

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
		apiId,
		apiHash,
		gotgproto.ClientTypePhone(phoneNumber),
		&gotgproto.ClientOpts{
			Session: sessionMaker.SqlSession(sqlite.Open("echobot")),
		},
	)
  
	go func() {
		err = client.Run(context.Background(), func(ctx context.Context) error {
      // Клиент работает, можно выполнять действия
			log.Println("Telegram client is running")
			return nil
		})
		if err != nil {
			log.Fatalln("failed to start client:", err)
		}
	}()
  go http_sever(client)
	err = client.Idle()
	if err != nil {
		log.Fatalln("failed client idle:", err)
	}
	//	go send_message("77085690946", client, "sosal")
	//	go send_message("77085690946", client, "sosal2")
}

func http_sever(client *gotgproto.Client) {
	http.HandleFunc("/sosal", send_handler(client))
	http.ListenAndServe(":8080", nil)

}

func send_handler(client *gotgproto.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			fmt.Print("goida")
			decoder := json.NewDecoder(r.Body)
			var msg Message
			err := decoder.Decode(&msg)
			log.Println(r.Body)
			if err != nil {
				http.Error(w, "Reject", http.StatusBadRequest)
			}
			go send_message(msg.PhoneNumber, client, msg.MessageText)
			log.Println(msg.MessageText)
			log.Println(msg.PhoneNumber)
		default:
			http.Error(w, "Reject", http.StatusBadRequest)
		}
	}
}

func send_message(phone string, client *gotgproto.Client, msg string) {
	ctx := context.Background()
	contacts := []tg.InputPhoneContact{
		{
			ClientID:  1,
			FirstName: "",
			LastName:  "",
			Phone:     phone,
		},
	}

	imported, err := client.API().ContactsImportContacts(
		ctx,
		contacts,
	)
	if err != nil {
		panic(err)
	}

	if len(imported.Users) > 0 {
		user := imported.Users[0].(*tg.User)
		log.Printf("Chat ID (User ID) for phone %s: %d", phone, user.ID)
		sender := message.NewSender(client.API())
		peer := &tg.InputPeerUser{
			UserID:     user.ID,
			AccessHash: user.AccessHash,
		}
		_, err = sender.To(peer).Text(ctx, msg)
		if err != nil {
			log.Fatalln("Error sending message:", err)
		}
	}
}
