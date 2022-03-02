package decoder

import (
	"github.com/virtualtam/netscape-go/ast"
	"github.com/virtualtam/netscape-go/types"
)

// Decode walks a Netscape Bookmark AST and returns the corresponding document.
func Decode(f ast.File) (types.Document, error) {
	return decodeFile(f)
}
