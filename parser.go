package netscape

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
)

const (
	NetscapeBookmarkDoctype string = "NETSCAPE-Bookmark-file-1"
)

// Parse reads a Netscape Bookmark document and processes token by token to
// build and return the corresponding AST.
func Parse(reader io.Reader) (*File, error) {
	p := newParser(reader)
	return p.parse()
}

type parser struct {
	decoder *xml.Decoder
	file    *File
}

func newParser(reader io.Reader) *parser {
	decoder := xml.NewDecoder(reader)
	file := &File{}

	return &parser{
		decoder: decoder,
		file:    file,
	}
}

func (p *parser) parse() (*File, error) {
	if err := p.verifyDoctype(); err != nil {
		return &File{}, err
	}

	for {
		tok, err := p.decoder.Token()
		if tok == nil || errors.Is(err, io.EOF) {
			break
		}

		switch tokType := tok.(type) {
		case xml.StartElement:
			switch tokType.Name.Local {
			case "TITLE":
				if err := p.parseTitle(&tokType); err != nil {
					return &File{}, err
				}
			}
		}
	}

	return p.file, nil
}

func (p *parser) parseTitle(start *xml.StartElement) error {
	var title struct {
		Value string `xml:",chardata"`
	}

	if err := p.decoder.DecodeElement(&title, start); err != nil {
		return ErrTokenUnexpected
	}

	p.file.Title = title.Value

	return nil
}

func (p *parser) verifyDoctype() error {
	tok, err := p.decoder.Token()

	if tok == nil || errors.Is(err, io.EOF) {
		return ErrDoctypeMissing
	}

	switch tokType := tok.(type) {
	case xml.Directive:
		if string(tokType) != fmt.Sprintf("DOCTYPE %s", NetscapeBookmarkDoctype) {
			return ErrDoctypeInvalid
		}
	default:
		return ErrDoctypeMissing
	}

	return nil
}
