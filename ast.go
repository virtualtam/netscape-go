// Copyright (c) VirtualTam
// SPDX-License-Identifier: MIT

package netscape

// A FileNode represents a Netscape Bookmark file.
type FileNode struct {
	Title string
	Root  FolderNode
}

// A FolderNode represents a bookmark (sub-)folder that may contain Bookmarks
// and child Folders.
type FolderNode struct {
	Parent *FolderNode

	Name        string
	Description string
	Attributes  map[string]string

	Bookmarks  []BookmarkNode
	Subfolders []FolderNode
}

// A BookmarkNode represents a Netscape bookmark.
type BookmarkNode struct {
	Href        string
	Title       string
	Description string
	Attributes  map[string]string
}
