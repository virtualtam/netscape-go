package netscape

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	NetscapeBookmarkDoctype string = "NETSCAPE-Bookmark-file-1"
)

var (
	ErrDoctypeMissing error = errors.New("missing DOCTYPE")
	ErrDoctypeInvalid error = errors.New("invalid DOCTYPE")

	ErrRootFolderMissing      error = errors.New("missing root folder (<H1> tag)")
	ErrParentFolderMissing    error = errors.New("missing parent folder")
	ErrFolderTitleEmpty       error = errors.New("empty folder title")
	ErrFolderStructureInvalid error = errors.New("invalid folder structure")
)

var (
	// UTF-8 Byte Order Mark.
	utf8bom = []byte{0xef, 0xbb, 0xbf}
)

// A ParseError is returned when we fail to parse a Netscape Bookmark token or XML
// element.
type ParseError struct {
	// Custom message for this ParseError.
	Msg string

	// Position in the input where the error was raised.
	Pos int64

	// Initial error raised while parsing the input.
	Err error
}

func newParseError(msg string, pos int64, inner error) error {
	return &ParseError{
		Msg: msg,
		Pos: pos,
		Err: inner,
	}
}

// Error returns the string representation for this error.
func (e *ParseError) Error() string {
	return fmt.Sprintf("%s at position %d: %s", e.Msg, e.Pos, e.Err)
}

// Is compares this Error with a target error to satisfy an equality check.
func (e *ParseError) Is(target error) bool {
	t, ok := target.(*ParseError)
	if !ok {
		return false
	}

	return e.Msg == t.Msg && e.Pos == t.Pos
}

// Unwrap returns the inner error wrapped by this Error.
func (e *ParseError) Unwrap() error {
	return e.Err
}

// Parse reads a Netscape Bookmark document and processes it token by token to
// build and return the corresponding AST.
func Parse(readseeker io.ReadSeeker) (*FileNode, error) {
	p := newParser(readseeker)
	return p.parse()
}

type parser struct {
	readseeker  io.ReadSeeker
	decoder     *xml.Decoder
	tokenOffset int64

	file            *FileNode
	currentDepth    int
	currentFolder   *FolderNode
	currentBookmark *BookmarkNode
}

// newXMLDecoder initializes and returns a xml.Decoder with strict mode disabled,
// to handle Netscape Bookmark format quirks.
func newXMLDecoder(reader io.Reader) *xml.Decoder {
	decoder := xml.NewDecoder(reader)

	decoder.Strict = false
	decoder.AutoClose = []string{
		"p",
	}

	return decoder
}

func newParser(readseeker io.ReadSeeker) *parser {
	decoder := newXMLDecoder(readseeker)
	file := &FileNode{}

	return &parser{
		readseeker: readseeker,
		decoder:    decoder,
		file:       file,
	}
}

func (p *parser) parse() (*FileNode, error) {
	if err := p.verifyDoctype(); err != nil {
		return &FileNode{}, err
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
					return &FileNode{}, err
				}
			case "H1":
				folder, err := p.parseFolder(&tokType)
				if err != nil {
					return &FileNode{}, err
				}

				p.file.Root = folder
				p.currentFolder = &p.file.Root
			case "DL", "DT":
				if p.currentFolder == nil {
					// The document must have a <H1>...</H1> root folder.
					return &FileNode{}, newParseError("failed to parse bookmarks", p.tokenOffset, ErrRootFolderMissing)
				}

				if err := p.parseBookmarks(); err != nil {
					return &FileNode{}, err
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
		return newParseError("failed to parse title", p.tokenOffset, err)
	}

	p.tokenOffset = p.decoder.InputOffset()
	p.file.Title = title.Value

	return nil
}

func (p *parser) parseFolder(start *xml.StartElement) (FolderNode, error) {
	var elt struct {
		Name string `xml:",chardata"`
	}

	if err := p.decoder.DecodeElement(&elt, start); err != nil {
		return FolderNode{}, newParseError("failed to parse folder", p.tokenOffset, err)
	}

	if elt.Name == "" {
		return FolderNode{}, newParseError("failed to parse folder", p.tokenOffset, ErrFolderTitleEmpty)
	}

	p.tokenOffset = p.decoder.InputOffset()

	folder := FolderNode{
		Name:       elt.Name,
		Attributes: map[string]string{},
	}

	for _, attr := range start.Attr {
		folder.Attributes[attr.Name.Local] = attr.Value
	}

	return folder, nil
}

func (p *parser) parseBookmarks() error {
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

				if p.currentFolder == nil {
					return newParseError("failed to parse bookmarks", p.tokenOffset, ErrFolderStructureInvalid)
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

				if p.currentFolder == nil {
					return newParseError("failed to parse H3 subfolder", p.tokenOffset, ErrFolderStructureInvalid)
				}

				folder.Parent = p.currentFolder
				p.currentFolder.Subfolders = append(p.currentFolder.Subfolders, folder)
				p.currentFolder = &p.currentFolder.Subfolders[len(p.currentFolder.Subfolders)-1]
				p.currentDepth++

				lastElementType = "H3"
			}
		case xml.EndElement:
			switch tokType.Name.Local {
			case "DL":
				if p.currentDepth < 0 {
					return newParseError("failed to parse bookmarks", p.tokenOffset, ErrFolderStructureInvalid)
				}

				p.currentDepth--
				p.currentFolder = p.currentFolder.Parent
			}
		}
	}

	return nil
}

func (p *parser) parseBookmark(start *xml.StartElement) (BookmarkNode, error) {
	var link struct {
		Title string `xml:",chardata"`
	}

	if err := p.decoder.DecodeElement(&link, start); err != nil {
		return BookmarkNode{}, newParseError("failed to parse bookmark", p.tokenOffset, err)
	}

	p.tokenOffset = p.decoder.InputOffset()

	bookmark := BookmarkNode{
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
		if err != nil {
			return "", newParseError("failed to parse description", p.tokenOffset, err)
		}

		p.tokenOffset = p.decoder.InputOffset()

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

		p.tokenOffset = p.decoder.InputOffset()

		switch tokType := tok.(type) {
		case xml.CharData:
			if bytes.Equal(tokType, utf8bom) {
				continue
			}
			return newParseError("unexpected character data", p.tokenOffset, ErrDoctypeInvalid)

		case xml.Directive:
			if string(tokType) != fmt.Sprintf("DOCTYPE %s", NetscapeBookmarkDoctype) {
				return newParseError(fmt.Sprintf("unknown DOCTYPE %s", string(tokType)), p.tokenOffset, ErrDoctypeInvalid)
			}
			return nil

		default:
			return ErrDoctypeMissing
		}
	}
}
