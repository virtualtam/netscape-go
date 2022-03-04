package types

import (
	"net/url"
	"time"
)

// A Bookmark represents a Netscape Bookmark.
type Bookmark struct {
	CreatedAt *time.Time
	UpdatedAt *time.Time

	Title string
	URL   url.URL

	Description string
	Private     bool
	Tags        []string

	Attributes map[string]string
}
