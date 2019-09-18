package types

import "time"

type DocumentCreate struct {
	DocumentID   string `json: document_id`
	Device       string `json: device`
	Event        string `json: event`
	Mimetype     string `json: mimetype`
	Snapshot     string `json: snapshot`
	CreationDate time.Time
}

type DocumentUpdate struct {
	Device   string `json: device`
	Event    string `json: event`
	Mimetype string `json: mimetype`
	Snapshot string `json: snapshot`
}

type All struct {
	Count int           `json: count`
	Items []interface{} `json: items`
}
