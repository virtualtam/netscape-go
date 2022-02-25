package netscape

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
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
			case "H1":
				folder, err := p.parseFolder(&tokType)
				if err != nil {
					return &File{}, err
				}

				p.file.Root = folder
			case "DL":
				bookmarks, err := p.parseBookmarks(&tokType)
				if err != nil {
					return &File{}, err
				}

				p.file.Root.Bookmarks = bookmarks
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

func (p *parser) parseFolder(start *xml.StartElement) (Folder, error) {
	var folder struct {
		Name string `xml:",chardata"`
	}

	if err := p.decoder.DecodeElement(&folder, start); err != nil {
		return Folder{}, ErrTokenUnexpected
	}

	return Folder{Name: folder.Name}, nil
}

func (p *parser) parseBookmarks(start *xml.StartElement) ([]Bookmark, error) {
	bookmarks := []Bookmark{}
	currentBookmarkIndex := -1

	for {
		tok, err := p.decoder.Token()
		if tok == nil || errors.Is(err, io.EOF) {
			break
		}

		switch tokType := tok.(type) {
		case xml.StartElement:
			switch tokType.Name.Local {
			case "A":
				bookmark, err := p.parseBookmark(&tokType)
				if err != nil {
					return []Bookmark{}, err
				}
				bookmarks = append(bookmarks, bookmark)
				currentBookmarkIndex++
			case "DD":
				description, err := p.parseBookmarkDescription()
				if err != nil {
					return []Bookmark{}, err
				}
				bookmarks[currentBookmarkIndex].Description = description
			}
		case xml.EndElement:
			if tokType.Name.Local == "DL" {
				return bookmarks, nil
			}
		}
	}

	return bookmarks, nil
}

func (p *parser) parseBookmark(start *xml.StartElement) (Bookmark, error) {
	var link struct {
		Title string `xml:",chardata"`
	}

	if err := p.decoder.DecodeElement(&link, start); err != nil {
		return Bookmark{}, ErrTokenUnexpected
	}

	bookmark := Bookmark{
		Title:      link.Title,
		Attributes: map[string]string{},
	}

	for _, attr := range start.Attr {
		if attr.Name.Local == "HREF" {
			bookmark.Href = attr.Value
			continue
		}

		bookmark.Attributes[attr.Name.Local] = attr.Value
	}

	return bookmark, nil
}

func (p *parser) parseBookmarkDescription() (string, error) {
	tok, err := p.decoder.Token()
	if tok == nil || errors.Is(err, io.EOF) {
		return "", ErrEOFUnexpected
	}

	switch tokType := tok.(type) {
	case xml.CharData:
		description := string(tokType)
		description = strings.TrimSpace(description)
		return description, nil
	}

	return "", ErrTokenUnexpected
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
