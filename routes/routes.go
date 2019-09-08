package routes

import (
  "net/http"
  "context"
  "encoding/json"
  "io/ioutil"
  "time"
  "../types"
  "../utils"
)

func snapshotsCreateOne(response http.ResponseWriter, request *http.Request) {

  if request.Body == nil {
			http.Error(response, "Please send a request body", 400)
			return
	}

  snap := new(types.Snapshot)
  curTime := time.Now()

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

  snap.CreationDate = curTime.Local()
  doc := utils.Join(curTime.String(), "_", snap.Device, "_", snap.Event)

  client.Collection("Snapshots").Doc(doc).Set(context.Background(), snap)

  jsonData, err := json.Marshal(snap)
  if err != nil {
     http.Error(response, err.Error(), http.StatusInternalServerError)
     return
  }

  response.Header().Set("Content-Type","application/json")
  response.WriteHeader(http.StatusOK)
  response.Write(jsonData)

}

func snapshotsGetAll(response http.ResponseWriter, request *http.Request) {

    //docs := client.Collection("Snapshots").Limit(25).Documents(context.Background()).GetAll

}

func snapshotsGetOne(response http.ResponseWriter, request *http.Request) { }
