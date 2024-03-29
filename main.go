package main

import (
	"context"
	"log"
	"net/http"

	"./routes"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func main() {

	op := option.WithCredentialsFile("./keys/accountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, op)
	if err != nil {
		log.Fatal(err)
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	http.ListenAndServe(":8000", routes.Create(client))

}
