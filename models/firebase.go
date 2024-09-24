package models

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var FirebaseDB *db.Client

func ConnectRealtimeDatabase() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("error loading .env file")
	}

	ctx := context.Background()

	// configure database URL
	conf := &firebase.Config{
		DatabaseURL: os.Getenv("FIREBASE_DB"),
	}

	// fetch service account key
	opt := option.WithCredentialsFile("golang-dev-3ffde-firebase-adminsdk-rwxba-5f885d8abb.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("error in initializing firebase app: ", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("error in creating firebase DB client: ", err)
	}

	FirebaseDB = client
}
