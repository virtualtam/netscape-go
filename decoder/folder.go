package decoder

import (
	"github.com/virtualtam/netscape-go"
	"github.com/virtualtam/netscape-go/ast"
)

func decodeFolder(f ast.Folder) (netscape.Folder, error) {
	folder := netscape.Folder{
		Name:        f.Name,
		Description: f.Description,
		Attributes:  map[string]string{},
	}

	for attr, value := range f.Attributes {
		switch attr {
		case createdAtAttr:
			createdAt, err := decodeDate(value)
			if err != nil {
				return netscape.Folder{}, err
			}
			folder.CreatedAt = createdAt
		case updatedAtAttr:
			updatedAt, err := decodeDate(value)
			if err != nil {
				return netscape.Folder{}, err
			}
			folder.UpdatedAt = updatedAt
		default:
			folder.Attributes[attr] = value
		}
	}

	for _, b := range f.Bookmarks {
		bookmark, err := decodeBookmark(b)
		if err != nil {
			return netscape.Folder{}, err
		}

		folder.Bookmarks = append(folder.Bookmarks, bookmark)
	}

	for _, sf := range f.Subfolders {
		subfolder, err := decodeFolder(sf)
		if err != nil {
			return netscape.Folder{}, err
		}

		folder.Subfolders = append(folder.Subfolders, subfolder)
	}

	return folder, nil
}
