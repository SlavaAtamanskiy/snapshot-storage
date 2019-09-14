package routes

import (
	"../types"
	"../utils"
	"encoding/json"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"net/http"
)

const CollectionName = "Snapshots"

func snapshotsCreateOne(response http.ResponseWriter, request *http.Request) {

	//check if there is body passed
	if request.Body == nil {
		http.Error(response, "No body passed for request", 400)
		return
	}

	//generating template
	snap := new(types.Snapshot)
	snap.DocumentID = utils.GenerateDocLink(CollectionName)
	snap.CreationDate = utils.GetLocalTime()

	//reading body from the request
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	//merging template data with data passed
	err = json.Unmarshal(body, &snap)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	//inserting document into db
	client.Collection(CollectionName).Doc(snap.DocumentID).Set(cnt, snap)

	//preparing response body json
	jsonData, err := json.Marshal(snap)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	//sending response
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonData)

}

func snapshotsGetAll(response http.ResponseWriter, request *http.Request) {

	//getting struct model for data structuring and docs from the firebase
	model := new(types.All)
	docs := client.Collection(CollectionName).Documents(cnt)
	defer docs.Stop()
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
	//getting parameters (in our case unique id)
	//query := request.URL.Query()
	//docId := query.Get("document_id")

}

func snapshotsDeleteOne(response http.ResponseWriter, request *http.Request) {

	//getting parameters (in our case unique id)
	query := request.URL.Query()
	docId := query.Get("document_id")

	_, err := client.Collection(CollectionName).Doc(docId).Delete(cnt)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
	}

}
