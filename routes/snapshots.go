package routes

import (
	"../types"
	"encoding/json"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func snapshotsCreateOne(response http.ResponseWriter, request *http.Request) {

	if request.Body == nil {
		http.Error(response, "Please send a request body", 400)
		return
	}

	snap := new(types.Snapshot)
	curTime := time.Now()
	snap.CreationDate = curTime.Local()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &snap)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	doc := strings.Join([]string{curTime.String(), snap.Device, snap.Event}, " ")
	client.Collection("Snapshots").Doc(doc).Set(cnt, snap)

	jsonData, err := json.Marshal(snap)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonData)

}

func snapshotsGetAll(response http.ResponseWriter, request *http.Request) {

	//getting struct model for data structuring and docs from the firebase
	model := new(types.All)
	docs := client.Collection("Snapshots").Documents(cnt)
	counter := 0

	//filling an array of items in the model with docs from db
	for {
		doc, err := docs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}
		model.Items = append(model.Items, doc.Data())
		counter++
	}
	model.Count = counter

	//converting model with data to JSON
	jsonString, err := json.Marshal(model)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	//sending a response
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonString)

}

func snapshotsGetOne(response http.ResponseWriter, request *http.Request) {

}
