package netscape

import (
	"time"
)

// A Document represents a collection of Netscape Bookmarks.
type Document struct {
	Title string `json:"title"`
	Root  Folder `json:"root"`
}

// Flatten returns a flat version of this Document, with all Bookmarks attached
// to the Root Folder.
func (d *Document) Flatten() *Document {
	flattenedRoot := d.Root.flatten()

	return &Document{
		Title: d.Title,
		Root:  *flattenedRoot,
	}
}

// A Folder represents a folder containing Netscape Bookmarks and child Folders.
type Folder struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Description string `json:"description,omitempty"`
	Name        string `json:"name"`

	Attributes map[string]string `json:"attributes,omitempty"`

	Bookmarks  []Bookmark `json:"bookmarks,omitempty"`
	Subfolders []Folder   `json:"subfolders,omitempty"`
}

func (f *Folder) flatten() *Folder {
	flattened := &Folder{
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
		Description: f.Description,
		Name:        f.Name,
		Attributes:  f.Attributes,
	}

	flattened.Bookmarks = append(flattened.Bookmarks, f.Bookmarks...)

	for _, subfolder := range f.Subfolders {
		flattenedSubfolder := subfolder.flatten()
		flattened.Bookmarks = append(flattened.Bookmarks, flattenedSubfolder.Bookmarks...)
	}

	return flattened
}

// A Bookmark represents a Netscape Bookmark.
type Bookmark struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title string `json:"title"`
	URL   string `json:"url"`

	Description string   `json:"description,omitempty"`
	Private     bool     `json:"private"`
	Tags        []string `json:"tags,omitempty"`

	Attributes map[string]string `json:"attributes,omitempty"`
}
