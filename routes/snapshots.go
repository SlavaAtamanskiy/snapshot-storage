package routes

import (
	"../types"
	"../utils"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"github.com/jinzhu/copier"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"net/http"
)

const CollectionName = "Snapshots"

func snapshotsDeleteOne(response http.ResponseWriter, request *http.Request) {

	//getting parameters (in our case unique id)
	query := request.URL.Query()
	docID := query.Get("document_id")
	if len(docID) == 0 {
		http.Error(response, "Bad request", 400)
		return
	}

	_, err := client.Collection(CollectionName).Doc(docID).Delete(cnt)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	var snap struct {
		Success bool
	}
	snap.Success = true

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

func snapshotsCreateOne(response http.ResponseWriter, request *http.Request) {

	//check if there is body passed
	if request.Body == nil {
		http.Error(response, "No body passed for request", 400)
		return
	}

	//generating template
	snap := new(types.DocumentCreate)

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
	snap.DocumentID = utils.GenerateDocLink(CollectionName)
	snap.CreationDate = utils.GetLocalTime()
	client.Collection(CollectionName).Doc(snap.DocumentID).Set(cnt, snap)

	//preparing response body json
	var r = new(struct {
		Success bool
		types.DocumentCreate
	})
	copier.Copy(r, snap)
	r.Success = true

	jsonData, err := json.Marshal(r)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	//sending response
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonData)

}

func snapshotsUpdateOne(response http.ResponseWriter, request *http.Request) {

	//getting parameters (in our case unique id)
	query := request.URL.Query()
	docID := query.Get("document_id")

	//check if there is body passed
	if request.Body == nil {
		http.Error(response, "No body passed for request", 400)
		return
	}

	//generating template
	snap := new(types.DocumentUpdate)

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

	//getting reference
	docRef := client.Collection(CollectionName).Doc(docID)
	//checking if exists
	_, err = docRef.Get(cnt)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	//updating
	_, err = docRef.Update(cnt, []firestore.Update{
		{Path: "Device", Value: snap.Device},
		{Path: "Event", Value: snap.Event},
		{Path: "Mimetype", Value: snap.Mimetype},
		{Path: "Snapshot", Value: snap.Snapshot},
	})
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	//preparing response body json
	var r = new(struct {
		Success    bool
		DocumentID string
		types.DocumentUpdate
	})
	copier.Copy(r, snap)
	r.Success = true
	r.DocumentID = docID

	jsonData, err := json.Marshal(r)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	//sending response
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonData)

}

func snapshotsGet(response http.ResponseWriter, request *http.Request) {

	//getting parameters (in our case unique id)
	query := request.URL.Query()
	docID := query.Get("document_id")
	if len(docID) == 0 {
		snapshotsGetAll(response, request)
	} else {
		snapshotsGetOne(response, request, docID)
	}

}

func snapshotsGetAll(response http.ResponseWriter, request *http.Request) {

	//getting struct model for data structuring and docs from the firebase
	model := new(struct {
		Success bool
		types.All
	})

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
	model.Success = true

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

func snapshotsGetOne(response http.ResponseWriter, request *http.Request, documentID string) {

	docs := client.Collection(CollectionName).Where("DocumentID", "==", documentID).Documents(cnt)
	defer docs.Stop()

	doc, err := docs.Next()
	if err == iterator.Done {
		http.Error(response, "No data found", 404)
		return
	}
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	//preparing response body json
	jsonData, err := json.Marshal(doc.Data())
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	//sending response
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonData)

}
