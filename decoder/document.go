package decoder

import (
	"github.com/virtualtam/netscape-go"
	"github.com/virtualtam/netscape-go/ast"
)

func decodeFile(f ast.File) (netscape.Document, error) {
	document := netscape.Document{
		Title: f.Title,
	}

	root, err := decodeFolder(f.Root)
	if err != nil {
		return netscape.Document{}, err
	}
	document.Root = root

	return document, nil
}
