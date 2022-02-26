package netscape

import (
	"encoding/xml"
	"io"
)

// newDecoder initializes and returns a xml.Decoder with strict mode disabled,
// to handle Netscape Bookmark format quirks.
func newDecoder(reader io.Reader) *xml.Decoder {
	decoder := xml.NewDecoder(reader)

	decoder.Strict = false
	decoder.AutoClose = []string{
		"p",
	}

	return decoder
}
