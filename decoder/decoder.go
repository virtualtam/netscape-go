// Package decoder implements decoding to convert AST nodes to Netscape Bookmark
// domain types.
package decoder

import (
	"sort"
	"strings"
	"time"

	"github.com/virtualtam/netscape-go/ast"
	"github.com/virtualtam/netscape-go/types"
)

// Decode walks a Netscape Bookmark AST and returns the corresponding document.
func Decode(f ast.File) (*types.Document, error) {
	var d Decoder
	return d.decodeFile(f)
}

// A Decoder walks a Netscape Bookmark AST and returns the corresponding
// document.
type Decoder struct {
	defaultTime time.Time
}

// NewDecoder initializes and returns a new Decoder.
func NewDecoder(defaultTime time.Time) *Decoder {
	return &Decoder{
		defaultTime: defaultTime,
	}
}

func (d *Decoder) decodeFile(f ast.File) (*types.Document, error) {
	document := types.Document{
		Title: f.Title,
	}

	root, err := d.decodeFolder(f.Root)
	if err != nil {
		return &types.Document{}, err
	}
	document.Root = root

	return &document, nil
}

func (d *Decoder) decodeFolder(f ast.Folder) (types.Folder, error) {
	folder := types.Folder{
		Name:        f.Name,
		Description: f.Description,
		Attributes:  map[string]string{},
	}

	for attr, value := range f.Attributes {
		switch attr {
		case createdAtAttr:
			createdAt, err := decodeDate(value)
			if err != nil {
				return types.Folder{}, err
			}
			folder.CreatedAt = &createdAt
		case updatedAtAttr:
			updatedAt, err := decodeDate(value)
			if err != nil {
				return types.Folder{}, err
			}
			folder.UpdatedAt = &updatedAt
		default:
			folder.Attributes[attr] = value
		}
	}

	for _, b := range f.Bookmarks {
		bookmark, err := d.decodeBookmark(b)
		if err != nil {
			return types.Folder{}, err
		}

		folder.Bookmarks = append(folder.Bookmarks, bookmark)
	}

	for _, sf := range f.Subfolders {
		subfolder, err := d.decodeFolder(sf)
		if err != nil {
			return types.Folder{}, err
		}

		folder.Subfolders = append(folder.Subfolders, subfolder)
	}

	return folder, nil
}

func (d Decoder) decodeBookmark(b ast.Bookmark) (types.Bookmark, error) {
	bookmark := types.Bookmark{
		Description: b.Description,
		Href:        b.Href,
		Title:       b.Title,
		Attributes:  map[string]string{},
	}

	for attr, value := range b.Attributes {
		switch attr {
		case createdAtAttr:
			createdAt, err := decodeDate(value)
			if err != nil {
				return types.Bookmark{}, err
			}
			bookmark.CreatedAt = &createdAt
		case updatedAtAttr:
			updatedAt, err := decodeDate(value)
			if err != nil {
				return types.Bookmark{}, err
			}
			bookmark.UpdatedAt = &updatedAt
		case privateAttr:
			if value == "1" {
				bookmark.Private = true
			}
		case tagsAttr:
			bookmark.Tags = d.decodeTags(b.Attributes)
		default:
			bookmark.Attributes[attr] = value
		}
	}

	return bookmark, nil
}

func (d *Decoder) decodeTags(attr map[string]string) []string {
	rawTags, ok := attr[tagsAttr]
	if !ok {
		return []string{}
	}

	tags := strings.Split(rawTags, ",")
	for index, tag := range tags {
		tags[index] = strings.TrimSpace(tag)
	}

	sort.Strings(tags)

	return tags
}
