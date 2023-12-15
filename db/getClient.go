package db

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"github.com/jschoedt/go-firestorm"
	"google.golang.org/api/option"
)

func getClient() (*firestorm.FSClient, context.Context) {
	opt := option.WithCredentialsFile("firebase-credentials.json")
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln("Error occured while trying to connect to the DB: ", err)
	}

	client, _ := app.Firestore(ctx)
	return firestorm.New(client, "ID", ""), ctx
}
