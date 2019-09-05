package routes

import (
  "fmt"
  "log"
  "net/http"
	"github.com/gorilla/mux"
  "encoding/json"
  "io/ioutil"
  "time"
  "../types"
)

func Create() *mux.Router {
     router := mux.NewRouter()
     router.HandleFunc("/", index)
	   router.HandleFunc("/snapshots", snapshotsCreateOne).Methods("POST")
	   router.HandleFunc("/snapshots", snapshotsGetAll).Methods("GET")
	   router.HandleFunc("/snapshots/{id}", snapshotsGetOne).Methods("GET")
     return router
}

func index(response http.ResponseWriter, request *http.Request) {
     fmt.Fprintf(response, "Home page")
}

func snapshotsCreateOne(response http.ResponseWriter, request *http.Request) {

  if request.Body == nil {
			http.Error(response, "Please send a request body", 400)
			return
	}

  snap := types.Snapshot {}

  body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		  log.Fatal(err)
	}

  err = json.Unmarshal(body, &snap)
  if err != nil {
     log.Fatal(err)
  }

  snap.CreationDate = time.Now().Local()

  jsonData, err := json.Marshal(snap)
  if err != nil {
    log.Fatal(err)
  }

  response.Header().Set("Content-Type","application/json")
  response.WriteHeader(http.StatusOK)
  response.Write(jsonData)

}

func snapshotsGetAll(response http.ResponseWriter, request *http.Request) { }
func snapshotsGetOne(response http.ResponseWriter, request *http.Request) { }
