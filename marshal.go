// Package netscape provides utilities to parse and export Web bookmarks using
// the Netscape Bookmark format.
package netscape

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/virtualtam/netscape-go/decoder"
	"github.com/virtualtam/netscape-go/encoder"
	"github.com/virtualtam/netscape-go/parser"
	"github.com/virtualtam/netscape-go/types"
)

// Marshal returns the Netscape Bookmark encoding of d.
func Marshal(d *types.Document) ([]byte, error) {
	var buf bytes.Buffer

	if err := encoder.NewEncoder(&buf).Encode(d); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

// Unmarshal unmarshals a []byte representation of a Netscape Bookmark
// file and returns the corresponding Document.
func Unmarshal(b []byte) (types.Document, error) {
	r := bytes.NewReader(b)
	return unmarshal(r)
}

// UnmarshalFile unmarshals a Netscape Bookmark file and returns the
// corresponding Document.
func UnmarshalFile(filePath string) (types.Document, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return types.Document{}, err
	}
	defer file.Close()

	return unmarshal(file)
}

// UnmarshalString unmarshals a string representation of a Netscape Bookmark
// file and returns the corresponding Document.
func UnmarshalString(data string) (types.Document, error) {
	r := strings.NewReader(data)
	return unmarshal(r)
}

func unmarshal(r io.ReadSeeker) (types.Document, error) {
	astFile, err := parser.Parse(r)
	if err != nil {
		return types.Document{}, err
	}

	return decoder.Decode(*astFile)
}