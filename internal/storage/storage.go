package storage

import "errors"

const (
	ErrNotFoundFmt = "record not found with id %s"
)

var ErrNoRecords = errors.New("no records")

type Word struct {
	ID       string                   `json:"id,omitempty"`
	Word     string                   `json:"word,omitempty"`
	Sections []map[string]interface{} `json:"sections,omitempty"`
}
