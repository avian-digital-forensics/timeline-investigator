package internal

import (
	"strings"

	"github.com/google/uuid"
)

type Response struct {
	Took     int    `json:"took,omitempty"`
	TimedOut bool   `json:"timed_out,omitempty"`
	Shards   Shards `json:"_shards,omitempty"`
	Hits     Hits   `json:"hits,omitempty"`
	ScrollID string `json:"_scroll_id,omitempty"`
}

type Shards struct {
	Failed     float64 `json:"failed,omitempty"`
	Skipped    float64 `json:"skipped,omitempty"`
	Successful float64 `json:"successful,omitempty"`
	Total      float64 `json:"total,omitempty"`
}

type Hits struct {
	Hits     []Hit   `json:"hits,omitempty"`
	MaxScore float64 `json:"max_score,omitempty"`
	Total    float64 `json:"took,omitempty"`
}

type Hit struct {
	ID     string      `json:"_id,omitempty"`
	Index  string      `json:"_index,omitempty"`
	Type   string      `json:"_type,omitempty"`
	Score  float64     `json:"_score,omitempty"`
	Source interface{} `json:"_source,omitempty"`
}

type QueryRequest struct {
	Query Query `json:"query"`
}

type Query struct {
	Match    interface{} `json:"match,omitempty"`
	Wildcard interface{} `json:"wildcard,omitempty"`
	IDs      interface{} `json:"ids,omitempty"`
	Bool     Bool        `json:"bool,omitempty"`
}

type Bool struct {
	Must []Must `json:"must,omitempty"`
}

type Must struct {
	MatchPhrasePrefix interface{} `json:"match_phrase_prefix,omitempty"`
}

// NewID generates a new ID
func NewID() string { return strings.ReplaceAll(uuid.New().String(), "-", "") }
