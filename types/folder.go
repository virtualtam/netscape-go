package types

import "time"

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
