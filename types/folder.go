package types

import "time"

// A Folder represents a folder containing Netscape Bookmarks and child Folders.
type Folder struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	Description string
	Name        string

	Attributes map[string]string

	Bookmarks  []Bookmark
	Subfolders []Folder
}
