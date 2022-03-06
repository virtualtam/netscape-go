// Package encoder implements encoding and printing of domain types for Netscape
// Bookmark files.
package encoder

import (
	"bufio"
	"io"

	"github.com/virtualtam/netscape-go/types"
)

// An Encoder writes Netscape Bookmark data to an output stream.
type Encoder struct {
	p printer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		printer{
			Writer: bufio.NewWriter(w),
			indent: "    ",
		},
	}
}

// Encode writes the Netscape Bookmark encoding of d to the stream.
func (e *Encoder) Encode(d *types.Document) error {
	if err := e.p.marshalDocument(d); err != nil {
		return err
	}
	return e.p.Flush()
}
