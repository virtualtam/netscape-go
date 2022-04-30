package netscape

import (
	"time"
)

// A Document represents a collection of Netscape Bookmarks.
type Document struct {
	Title string `json:"title"`
	Root  Folder `json:"root"`
}

// A Folder represents a folder containing Netscape Bookmarks and child Folders.
type Folder struct {
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	Description string `json:"description,omitempty"`
	Name        string `json:"name"`

	Attributes map[string]string `json:"attributes,omitempty"`

	Bookmarks  []Bookmark `json:"bookmarks,omitempty"`
	Subfolders []Folder   `json:"subfolders,omitempty"`
}

// A Bookmark represents a Netscape Bookmark.
type Bookmark struct {
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	Title string `json:"title"`
	URL   string `json:"url"`

	Description string   `json:"description,omitempty"`
	Private     bool     `json:"private"`
	Tags        []string `json:"tags,omitempty"`

	Attributes map[string]string `json:"attributes,omitempty"`
}
