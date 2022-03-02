package decoder

import (
	"github.com/virtualtam/netscape-go/ast"
	"github.com/virtualtam/netscape-go/types"
)

func decodeFile(f ast.File) (types.Document, error) {
	document := types.Document{
		Title: f.Title,
	}

	root, err := decodeFolder(f.Root)
	if err != nil {
		return types.Document{}, err
	}
	document.Root = root

	return document, nil
}
