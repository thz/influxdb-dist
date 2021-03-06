package influxql

import (
	"encoding/json"
	"errors"

	"github.com/influxdata/influxdb/models"
)

const (
	// WarningLevel is the message level for a warning.
	WarningLevel = "warning"
)

// TagSet is a fundamental concept within the query system. It represents a composite series,
// composed of multiple individual series that share a set of tag attributes.
type TagSet struct {
	Tags       map[string]string
	Filters    []Expr
	SeriesKeys []string
	Key        []byte
}

// AddFilter adds a series-level filter to the Tagset.
func (t *TagSet) AddFilter(key string, filter Expr) {
	t.SeriesKeys = append(t.SeriesKeys, key)
	t.Filters = append(t.Filters, filter)
}

// Message represents a user-facing message to be included with the result.
type Message struct {
	Level string `json:"level"`
	Text  string `json:"text"`
}

// Result represents a resultset returned from a single statement.
// Rows represents a list of rows that can be sorted consistently by name/tag.
type Result struct {
	// StatementID is just the statement's position in the query. It's used
	// to combine statement results if they're being buffered in memory.
	StatementID int `json:"-"`
	Series      models.Rows
	Messages    []*Message
	Err         error
}

// MarshalJSON encodes the result into JSON.
func (r *Result) MarshalJSON() ([]byte, error) {
	// Define a struct that outputs "error" as a string.
	var o struct {
		Series   []*models.Row `json:"series,omitempty"`
		Messages []*Message    `json:"messages,omitempty"`
		Err      string        `json:"error,omitempty"`
	}

	// Copy fields to output struct.
	o.Series = r.Series
	o.Messages = r.Messages
	if r.Err != nil {
		o.Err = r.Err.Error()
	}

	return json.Marshal(&o)
}

// UnmarshalJSON decodes the data into the Result struct
func (r *Result) UnmarshalJSON(b []byte) error {
	var o struct {
		Series   []*models.Row `json:"series,omitempty"`
		Messages []*Message    `json:"messages,omitempty"`
		Err      string        `json:"error,omitempty"`
	}

	err := json.Unmarshal(b, &o)
	if err != nil {
		return err
	}
	r.Series = o.Series
	r.Messages = o.Messages
	if o.Err != "" {
		r.Err = errors.New(o.Err)
	}
	return nil
}
