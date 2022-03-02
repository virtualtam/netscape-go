package decoder

import (
	"github.com/virtualtam/netscape-go/ast"
	"github.com/virtualtam/netscape-go/types"
)

func decodeFolder(f ast.Folder) (types.Folder, error) {
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
			folder.CreatedAt = createdAt
		case updatedAtAttr:
			updatedAt, err := decodeDate(value)
			if err != nil {
				return types.Folder{}, err
			}
			folder.UpdatedAt = updatedAt
		default:
			folder.Attributes[attr] = value
		}
	}

	for _, b := range f.Bookmarks {
		bookmark, err := decodeBookmark(b)
		if err != nil {
			return types.Folder{}, err
		}

		folder.Bookmarks = append(folder.Bookmarks, bookmark)
	}

	for _, sf := range f.Subfolders {
		subfolder, err := decodeFolder(sf)
		if err != nil {
			return types.Folder{}, err
		}

		folder.Subfolders = append(folder.Subfolders, subfolder)
	}

	return folder, nil
}
