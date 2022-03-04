package parser

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/virtualtam/netscape-go/ast"
)

const (
	NetscapeBookmarkDoctype string = "NETSCAPE-Bookmark-file-1"
)

var (
	// Byte order mark
	utf8bom = []byte{0xef, 0xbb, 0xbf}
)

// Parse reads a Netscape Bookmark document and processes it token by token to
// build and return the corresponding AST.
func Parse(readseeker io.ReadSeeker) (*ast.File, error) {
	p := newParser(readseeker)
	return p.parse()
}

type parser struct {
	readseeker      io.ReadSeeker
	decoder         *xml.Decoder
	file            *ast.File
	currentFolder   *ast.Folder
	currentBookmark *ast.Bookmark
}

func newParser(readseeker io.ReadSeeker) *parser {
	decoder := newDecoder(readseeker)
	file := &ast.File{}

	return &parser{
		readseeker: readseeker,
		decoder:    decoder,
		file:       file,
	}
}

func (p *parser) parse() (*ast.File, error) {
	if err := p.verifyDoctype(); err != nil {
		return &ast.File{}, err
	}

	for {
		tok, err := p.decoder.Token()
		if tok == nil || errors.Is(err, io.EOF) {
			break
		}

		switch tokType := tok.(type) {
		case xml.StartElement:
			switch tokType.Name.Local {
			case "TITLE", "Title":
				if err := p.parseTitle(&tokType); err != nil {
					return &ast.File{}, err
				}
			case "H1":
				folder, err := p.parseFolder(&tokType)
				if err != nil {
					return &ast.File{}, err
				}

				p.file.Root = folder
				p.currentFolder = &p.file.Root
			case "DL", "DT":
				if err := p.parseBookmarks(&tokType); err != nil {
					return &ast.File{}, err
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

func (p *parser) parseFolder(start *xml.StartElement) (ast.Folder, error) {
	var elt struct {
		Name string `xml:",chardata"`
	}

	if err := p.decoder.DecodeElement(&elt, start); err != nil {
		return ast.Folder{}, ErrTokenUnexpected
	}

	folder := ast.Folder{
		Name:       elt.Name,
		Attributes: map[string]string{},
	}

	for _, attr := range start.Attr {
		folder.Attributes[attr.Name.Local] = attr.Value
	}

	return folder, nil
}

func (p *parser) parseBookmarks(start *xml.StartElement) error {
	var lastElementType string

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
					return err
				}
				p.currentFolder.Bookmarks = append(p.currentFolder.Bookmarks, bookmark)
				p.currentBookmark = &p.currentFolder.Bookmarks[len(p.currentFolder.Bookmarks)-1]
				lastElementType = "A"
			case "DD":
				description, err := p.parseDescription()
				if err != nil {
					return err
				}

				switch lastElementType {
				case "A":
					p.currentBookmark.Description = description
				case "H3":
					p.currentFolder.Description = description
				}
			case "H3":
				folder, err := p.parseFolder(&tokType)
				if err != nil {
					return err
				}

				folder.Parent = p.currentFolder
				p.currentFolder.Subfolders = append(p.currentFolder.Subfolders, folder)
				p.currentFolder = &p.currentFolder.Subfolders[len(p.currentFolder.Subfolders)-1]
				lastElementType = "H3"
			}
		case xml.EndElement:
			if tokType.Name.Local == "DL" {
				p.currentFolder = p.currentFolder.Parent
			}
		}
	}

	return nil
}

func (p *parser) parseBookmark(start *xml.StartElement) (ast.Bookmark, error) {
	var link struct {
		Title string `xml:",chardata"`
	}

	if err := p.decoder.DecodeElement(&link, start); err != nil {
		return ast.Bookmark{}, ErrTokenUnexpected
	}

	bookmark := ast.Bookmark{
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

// parseDescription returns a string containing all data following a <DD>
// element, and preceding either a <DT> or </DL> element.
//
// Leading and trailing whitespace is trimmed from the returned string.
//
// A description may contain text and HTML elements.
func (p *parser) parseDescription() (string, error) {
	startOffset := p.decoder.InputOffset()
	endOffset := startOffset

	var readseekOffset int64

	// As the description may contain either text or HTML elements, we do not
	// directly process the stream of XML tokens, and instead look for the start
	// and end offsets of the description data in the underlying io.ReadSeeker.
loop:
	for {
		tok, err := p.decoder.Token()
		if tok == nil || errors.Is(err, io.EOF) {
			return "", ErrEOFUnexpected
		}

		switch tokType := tok.(type) {
		case xml.CharData:
			endOffset = p.decoder.InputOffset()
		case xml.StartElement:
			if tokType.Name.Local == "DL" || tokType.Name.Local == "DT" {
				break loop
			}
		case xml.EndElement:
			if tokType.Name.Local == "DD" || tokType.Name.Local == "DL" {
				break loop
			}

			endOffset = p.decoder.InputOffset()
		}
	}

	readseekOffset, err := p.readseeker.Seek(0, io.SeekCurrent)
	if err != nil {
		return "", err
	}
	// read raw data between start and end offsets
	dataLen := int(endOffset - startOffset)

	data := make([]byte, dataLen)
	_, err = p.readseeker.Seek(startOffset, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("description: failed to reset readseeker position to %d: %w", startOffset, err)
	}

	nRead, err := p.readseeker.Read(data)
	if err != nil {
		return "", fmt.Errorf("description: failed to read data in range [%d:%d]: %w", startOffset, endOffset, err)
	}

	if nRead != dataLen {
		return "", fmt.Errorf("description: expected to read %d bytes, read %d", dataLen, nRead)
	}

	// reset the io.ReadSeeker position
	_, err = p.readseeker.Seek(readseekOffset, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("description: failed to reset readseeker position to %d: %w", startOffset, err)
	}

	// sanitize data
	description := strings.TrimSpace(string(data))
	return description, nil
}

func (p *parser) verifyDoctype() error {
	for {
		tok, err := p.decoder.Token()

		if tok == nil || errors.Is(err, io.EOF) {
			return ErrDoctypeMissing
		}

		switch tokType := tok.(type) {
		case xml.CharData:
			if bytes.Equal(tokType, utf8bom) {
				continue
			}
			return ErrTokenUnexpected

		case xml.Directive:
			if string(tokType) != fmt.Sprintf("DOCTYPE %s", NetscapeBookmarkDoctype) {
				return ErrDoctypeInvalid
			}
			return nil

		default:
			return ErrDoctypeMissing
		}
	}
}
