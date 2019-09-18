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

	router.HandleFunc("/snapshots", snapshotsGet).Methods("GET")
	router.HandleFunc("/snapshots", snapshotsCreateOne).Methods("POST")
	router.HandleFunc("/snapshots", snapshotsUpdateOne).Methods("PUT")
	router.HandleFunc("/snapshots", snapshotsDeleteOne).Methods("DELETE")

	return router

}
