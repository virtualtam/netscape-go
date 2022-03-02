package decoder

import (
	"github.com/virtualtam/netscape-go"
	"github.com/virtualtam/netscape-go/ast"
)

// Decode walks a Netscape Bookmark AST and returns the corresponding document.
func Decode(f ast.File) (netscape.Document, error) {
	return decodeFile(f)
}
