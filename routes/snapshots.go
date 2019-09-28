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
		utils.CustomError(response, "Bad request", 400)
		return
	}

	//getting reference
	docRef := client.Collection(CollectionName).Doc(docID)

	//checking if exists
	if _, err := docRef.Get(ctx); err != nil {
		utils.HandleError(response, err)
		return
	}

	//deleting
	if _, err := docRef.Delete(ctx); err != nil {
		utils.HandleError(response, err)
		return
	}

	//sending response
	utils.ResponseOk(response, nil)

}

func snapshotsCreateOne(response http.ResponseWriter, request *http.Request) {

	//check if there is body passed
	if request.Body == nil {
		utils.CustomError(response, "Bad request", 400)
		return
	}

	//generating template
	snap := new(types.DocumentCreate)

	//reading body from the request
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		utils.HandleError(response, err)
		return
	}

	//merging template data with data passed
	if err = json.Unmarshal(body, &snap); err != nil {
		utils.HandleError(response, err)
		return
	}

	//inserting document into db
	snap.DocumentID = utils.GenerateDocLink(CollectionName)
	snap.CreationDate = utils.GetLocalTime()
	client.Collection(CollectionName).Doc(snap.DocumentID).Set(ctx, utils.DecapitalizeStruct(*snap))

	//preparing response body json template
	var r = new(struct {
		Success bool `json:"success"`
		types.DocumentCreate
	})

	//copy data into response body template
	copier.Copy(r, snap)
	r.Success = true

	//convert response body into json
	jsonData, err := json.Marshal(r)
	if err != nil {
		utils.HandleError(response, err)
		return
	}

	//sending response
	utils.ResponseOk(response, jsonData)

}

func snapshotsUpdateOne(response http.ResponseWriter, request *http.Request) {

	//getting parameters (in our case unique id)
	query := request.URL.Query()
	docID := query.Get("document_id")

	//check if there are body and parameters passed
	if request.Body == nil || len(docID) == 0 {
		utils.CustomError(response, "Bad request", 400)
		return
	}

	//generating template
	snap := new(types.DocumentUpdate)

	//reading body from the request
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		utils.HandleError(response, err)
		return
	}

	//merging template data with data passed
	if err = json.Unmarshal(body, &snap); err != nil {
		utils.HandleError(response, err)
		return
	}
	
	//getting reference
	docRef := client.Collection(CollectionName).Doc(docID)
	
	//checking if exists
	if _, err = docRef.Get(ctx); err != nil {
		utils.HandleError(response, err)
		return
	}

	//updating
	if _, err = docRef.Update(ctx, []firestore.Update{
		{Path: "device", Value: snap.Device},
		{Path: "event", Value: snap.Event},
		{Path: "mimetype", Value: snap.Mimetype},
		{Path: "snapshot", Value: snap.Snapshot},
	}); err != nil {
		utils.HandleError(response, err)
		return
	}

	//preparing response body json
	var r = new(struct {
		Success    bool   `json:"success"`
		DocumentID string `json:"document_id"`
		types.DocumentUpdate
	})
	copier.Copy(r, snap)
	r.Success = true
	r.DocumentID = docID
	//r.CreationDate = doc.Data().creation_date

	jsonData, err := json.Marshal(r)
	if err != nil {
		utils.HandleError(response, err)
		return
	}

	//sending response
	utils.ResponseOk(response, jsonData)

}

/*
  function: snapshotsGet
  params:   document_id unique name of one specific document
	          limit       number of documents to limit
	          offset      number of documents to leave behind
	purpose: gets all the entries from Snapshots collection
	         or some of them using limit and offset params
					 or one specific entry by its document_id
*/
func snapshotsGet(response http.ResponseWriter, request *http.Request) {

	//getting parameters (in our case unique id)
	query := request.URL.Query()
	docID := query.Get("document_id")
	limit := query.Get("limit")
	offset := query.Get("offset")
	if len(docID) == 0 {
		lim, offs := utils.GetLimitOffset(limit, offset)
		snapshotsGetAll(response, request, lim, offs)
	} else {
		snapshotsGetOne(response, request, docID)
	}

}

func snapshotsGetAll(response http.ResponseWriter, request *http.Request, limit, offset int) {

	//getting struct model for data structuring and docs from the firebase
	model := new(struct {
		Success bool `json:"success"`
		types.All
	})

	var docs *firestore.DocumentIterator
	query := client.Collection(CollectionName).OrderBy("creation_date", firestore.Desc)

	if limit == 0 && offset == 0 {
		docs = query.Documents(ctx)
	} else if limit == 0 && offset != 0 {
		docs = query.Offset(offset).Documents(ctx)
	} else if limit != 0 && offset == 0 {
		docs = query.Limit(limit).Documents(ctx)
	} else {
		docs = query.Offset(offset).Limit(limit).Documents(ctx)
	}
	defer docs.Stop()
	counter := 0

	//filling an array of items in the model with docs from db
	for {
		doc, err := docs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			utils.HandleError(response, err)
			return
		}
		model.Items = append(model.Items, doc.Data())
		counter++
	}
	model.Count = counter
	model.Success = true

	//converting model with data to JSON
	jsonData, err := json.Marshal(model)
	if err != nil {
		utils.HandleError(response, err)
		return
	}

	//sending response
	utils.ResponseOk(response, jsonData)

}

func snapshotsGetOne(response http.ResponseWriter, request *http.Request, docID string) {

	//getting reference
	docRef := client.Collection(CollectionName).Doc(docID)

	//checking if exists
	doc, err := docRef.Get(ctx)
	if err != nil {
		utils.HandleError(response, err)
		return
	}

	//preparing response body json
	jsonData, err := json.Marshal(doc.Data())
	if err != nil {
		utils.HandleError(response, err)
		return
	}

	//sending response
	utils.ResponseOk(response, jsonData)

}
