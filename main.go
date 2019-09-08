package main

import (
	"./routes"
	"context"
	"firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"net/http"
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
