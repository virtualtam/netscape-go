package encoder

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/virtualtam/netscape-go/types"
)

type printer struct {
	*bufio.Writer
	depth  int
	indent string
}

func (p *printer) writeString(s string) (int, error) {
	for i := 0; i < p.depth; i++ {
		if n, err := p.WriteString(p.indent); err != nil {
			return n, err
		}
	}

	return p.WriteString(s)
}

func (p *printer) marshalDocument(d *types.Document) error {
	const header = `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
     It will be read and overwritten.
     DO NOT EDIT! -->`

	_, err := p.WriteString(fmt.Sprintf("%s\n<TITLE>%s</TITLE>\n", header, d.Title))
	if err != nil {
		return err
	}

	if err := p.marshalFolder(&d.Root, true); err != nil {
		return err
	}

	return nil
}

type netscapeH3 struct {
	XMLName xml.Name `xml:"H3"`

	Name string `xml:",chardata"`

	AddDate      string `xml:"ADD_DATE,attr,omitempty"`
	LastModified string `xml:"LAST_MODIFIED,attr,omitempty"`

	Attrs []xml.Attr `xml:",attr,omitempty"`
}

func newNetscapeH3(f *types.Folder) *netscapeH3 {
	h3 := netscapeH3{
		Name: f.Name,
	}

	if f.CreatedAt != nil {
		h3.AddDate = fmt.Sprintf("%d", f.CreatedAt.Unix())
	}

	if f.UpdatedAt != nil {
		h3.LastModified = fmt.Sprintf("%d", f.UpdatedAt.Unix())
	}

	for k, v := range f.Attributes {
		attr := xml.Attr{Name: xml.Name{Local: k}, Value: v}
		h3.Attrs = append(h3.Attrs, attr)
	}

	return &h3
}

func (p *printer) marshalFolderHeader(f *types.Folder) error {
	h3 := newNetscapeH3(f)

	m, err := xml.Marshal(h3)
	if err != nil {
		return err
	}

	_, err = p.writeString(fmt.Sprintf("<DT>%s\n", string(m)))
	if err != nil {
		return err
	}

	return nil
}

func (p *printer) marshalFolder(f *types.Folder, isRoot bool) error {
	if !isRoot {
		if err := p.marshalFolderHeader(f); err != nil {
			return err
		}
	} else {
		_, err := p.writeString(fmt.Sprintf("<H1>%s</H1>\n", f.Name))
		if err != nil {
			return err
		}
	}

	if f.Description != "" {
		_, err := p.writeString(fmt.Sprintf("<DD>%s\n", f.Description))
		if err != nil {
			return err
		}
	}
	_, err := p.writeString("<DL><p>\n")
	if err != nil {
		return err
	}

	p.depth++

	for _, b := range f.Bookmarks {
		if err := p.marshalBookmark(&b); err != nil {
			return err
		}
	}

	for _, sf := range f.Subfolders {
		if err := p.marshalFolder(&sf, false); err != nil {
			return err
		}
	}

	p.depth--

	_, err = p.writeString("</DL><p>\n")
	if err != nil {
		return err
	}

	return nil
}

type netscapeA struct {
	XMLName xml.Name `xml:"A"`

	Title string `xml:",chardata"`
	Href  string `xml:"HREF,attr"`

	AddDate      string `xml:"ADD_DATE,attr,omitempty"`
	LastModified string `xml:"LAST_MODIFIED,attr,omitempty"`

	Private int    `xml:"PRIVATE,attr"`
	Tags    string `xml:"TAGS,attr,omitempty"`

	Attrs []xml.Attr `xml:",attr,omitempty"`
}

func newNetscapeA(b *types.Bookmark) *netscapeA {
	a := netscapeA{
		Href:  b.Href,
		Title: b.Title,
		Tags:  strings.Join(b.Tags, ","),
		Attrs: []xml.Attr{},
	}

	if b.CreatedAt != nil {
		a.AddDate = fmt.Sprintf("%d", b.CreatedAt.Unix())
	}

	if b.UpdatedAt != nil {
		a.LastModified = fmt.Sprintf("%d", b.UpdatedAt.Unix())
	}

	if b.Private {
		a.Private = 1
	}

	for k, v := range b.Attributes {
		attr := xml.Attr{Name: xml.Name{Local: k}, Value: v}
		a.Attrs = append(a.Attrs, attr)
	}

	return &a
}

func (p *printer) marshalBookmark(b *types.Bookmark) error {
	a := newNetscapeA(b)

	m, err := xml.Marshal(a)
	if err != nil {
		return err
	}

	_, err = p.writeString(fmt.Sprintf("<DT>%s\n", string(m)))
	if err != nil {
		return err
	}

	if b.Description != "" {
		_, err = p.writeString(fmt.Sprintf("<DD>%s\n", b.Description))
		if err != nil {
			return err
		}
	}

	return nil
}