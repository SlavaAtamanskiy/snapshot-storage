package routes

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/gorilla/mux"
)

var client *firestore.Client
var cnt context.Context

func Create(c *firestore.Client) *mux.Router {

	client = c
	cnt = context.Background()

	router := mux.NewRouter()

	router.HandleFunc("/snapshots", snapshotsGetAll).Methods("GET")
	router.HandleFunc("/snapshots/{id}", snapshotsGetOne).Methods("GET")
	router.HandleFunc("/snapshots", snapshotsCreateOne).Methods("POST")
	router.HandleFunc("/snapshots", snapshotsDeleteOne).Methods("DELETE")

	return router

}
