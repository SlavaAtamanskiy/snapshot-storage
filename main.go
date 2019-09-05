package main

import (
        _ "fmt"
        "context"
        "log"
        "firebase.google.com/go"
        "google.golang.org/api/option"
        "net/http"
        "./routes"
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

  http.ListenAndServe(":8000", routes.Create())

}
