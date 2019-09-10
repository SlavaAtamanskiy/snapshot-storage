package routes

import (
	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"context"
)

var client *firestore.Client
var cnt     context.Context

func Create(c *firestore.Client) *mux.Router {

	client = c
  cnt    = context.Background()

	router := mux.NewRouter()
	router.HandleFunc("/snapshots", snapshotsCreateOne).Methods("POST")
	router.HandleFunc("/snapshots", snapshotsGetAll).Methods("GET")
	router.HandleFunc("/snapshots/{id}", snapshotsGetOne).Methods("GET")

	return router

}
