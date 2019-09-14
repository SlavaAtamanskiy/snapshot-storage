package types

import "time"

type Snapshot struct {
	DocumentID   string `json: document_id`
	Device       string `json: device`
	Event        string `json: event`
	Mimetype     string `json: mimetype`
	Snapshot     string `json: snapshot`
	CreationDate time.Time
}

type All struct {
	Count int           `json: count`
	Items []interface{} `json: items`
}
