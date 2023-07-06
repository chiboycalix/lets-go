package data

import (
	"time"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // omit from JSON
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"` // omit empty values
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}
