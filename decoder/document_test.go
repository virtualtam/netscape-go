package decoder

import (
	"net/url"
	"testing"

	"github.com/virtualtam/netscape-go"
	"github.com/virtualtam/netscape-go/ast"
)

func TestDecodeFile(t *testing.T) {
	cases := []struct {
		tname string
		file  ast.File
		want  netscape.Document
	}{
		{
			tname: "empty document",
		},
		{
			tname: "flat document",
			file: ast.File{
				Title: "Bookmarks",
				Root: ast.Folder{
					Name:        "Test Folder",
					Description: "Add bookmarks to this folder",
					Bookmarks: []ast.Bookmark{
						{
							Href:  "https://domain.tld",
							Title: "Test Domain",
						},
						{
							Description: "Second test",
							Href:        "https://test.domain.tld",
							Title:       "Test Domain II",
						},
					},
				},
			},
			want: netscape.Document{
				Title: "Bookmarks",
				Root: netscape.Folder{
					Name:        "Test Folder",
					Description: "Add bookmarks to this folder",
					Bookmarks: []netscape.Bookmark{
						{
							URL: url.URL{
								Scheme: "https",
								Host:   "domain.tld",
							},
							Title: "Test Domain",
						},
						{
							Description: "Second test",
							URL: url.URL{
								Scheme: "https",
								Host:   "test.domain.tld",
							},
							Title: "Test Domain II",
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.tname, func(t *testing.T) {
			got, err := decodeFile(tc.file)

			if err != nil {
				t.Errorf("expected no error, got %q", err)
			}

			if got.Title != tc.want.Title {
				t.Errorf("want title %q, got %q", tc.want.Title, got.Title)
			}

			assertFoldersEqual(t, got.Root, tc.want.Root)
		})
	}
}
