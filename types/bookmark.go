package types

import (
	"net/url"
	"time"
)

// A Bookmark represents a Netscape Bookmark.
type Bookmark struct {
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	Title string  `json:"title"`
	URL   url.URL `json:"-"`

	Description string   `json:"description,omitempty"`
	Private     bool     `json:"private"`
	Tags        []string `json:"tags,omitempty"`

	Attributes map[string]string `json:"attributes,omitempty"`
}
