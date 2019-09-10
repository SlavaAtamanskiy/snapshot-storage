package types

import "time"

type Snapshot struct {
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
