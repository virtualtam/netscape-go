package netscape

// A File Node represents a Netscape Bookmark file.
type File struct {
	Title string
	Root  Folder
}

// A Folder Node represents a bookmark (sub-)folder that may contain Bookmarks
// and child Folders.
type Folder struct {
	Name      string
	Bookmarks []Bookmark
}

// A Bookmark Node represents a Netscape bookmark.
type Bookmark struct {
	Href        string
	Title       string
	Description string
}
