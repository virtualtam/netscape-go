package netscape

import (
	"errors"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		tname   string
		input   string
		want    File
		wantErr error
	}{
		// nominal cases
		{
			tname: "flat document",
			input: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
     It will be read and overwritten.
     DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
<DT><A HREF="https://domain.tld">Test Domain</A>
<DT><A HREF="https://desc.domain.tld">Test Domain (with description)</A>
<DD>Look! A short description for this bookmark.
<DT><A HREF="https://emptydesc.domain.tld">Test Domain (with empty description)</A>
<DD>
</DL><p>
`,
			want: File{
				Title: "Bookmarks",
				Root: Folder{
					Name: "Bookmarks",
					Bookmarks: []Bookmark{
						{
							Href:  "https://domain.tld",
							Title: "Test Domain",
						},
						{
							Href:        "https://desc.domain.tld",
							Title:       "Test Domain (with description)",
							Description: "Look! A short description for this bookmark.",
						},
						{
							Href:  "https://emptydesc.domain.tld",
							Title: "Test Domain (with empty description)",
						},
					},
				},
			},
		},
		{
			tname: "bookmark with attributes",
			input: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><A HREF="https://domain.tld" ADD_DATE="151637044" PRIVATE="1" TAGS="test tags">Test Domain</A>
</DL><p>
`,
			want: File{
				Title: "Bookmarks",
				Root: Folder{
					Name: "Bookmarks",
					Bookmarks: []Bookmark{
						{
							Href:  "https://domain.tld",
							Title: "Test Domain",
							Attributes: map[string]string{
								"ADD_DATE": "151637044",
								"PRIVATE":  "1",
								"TAGS":     "test tags",
							},
						},
					},
				},
			},
		},

		// error cases
		{
			tname:   "empty document",
			wantErr: ErrDoctypeMissing,
		},
		{
			tname:   "missing DOCTYPE",
			input:   `<!-- No DOCTYPE\n  -->\n`,
			wantErr: ErrDoctypeMissing,
		},
		{
			tname:   "invalid DOCTYPE",
			input:   `<!DOCTYPE dummy SYSTEM "dummy.dtd">`,
			wantErr: ErrDoctypeInvalid,
		},
		{
			tname: "incomplete TITLE",
			input: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<TITLE>Bookmarks`,
			wantErr: ErrTokenUnexpected,
		},
		{
			tname: "incomplete H1",
			input: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<TITLE>Bookmarks</TITLE>
<H1>Bookma`,
			wantErr: ErrTokenUnexpected,
		},
	}

	for _, tc := range cases {
		t.Run(tc.tname, func(t *testing.T) {
			r := strings.NewReader(tc.input)

			got, err := Parse(r)

			if tc.wantErr != nil {
				if err == nil {
					t.Error("expected an error, got none")
				} else if !errors.Is(err, tc.wantErr) {
					t.Errorf("want error %q, got %q", tc.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("expected no error, got %q", err)
			}

			if got.Title != tc.want.Title {
				t.Errorf("want title %q, got %q", tc.want.Title, got.Title)
			}

			if got.Root.Name != tc.want.Root.Name {
				t.Errorf("want root folder name %q, got %q", tc.want.Root.Name, got.Root.Name)
			}

			if len(got.Root.Bookmarks) != len(tc.want.Root.Bookmarks) {
				t.Errorf("want %d bookmarks in the root folder, got %d", len(tc.want.Root.Bookmarks), len(got.Root.Bookmarks))
				return
			}

			for index, wantBookmark := range tc.want.Root.Bookmarks {
				if got.Root.Bookmarks[index].Description != wantBookmark.Description {
					t.Errorf("want bookmark %d description %q, got %q", index, wantBookmark.Description, got.Root.Bookmarks[index].Description)
				}

				if got.Root.Bookmarks[index].Href != wantBookmark.Href {
					t.Errorf("want bookmark %d href %q, got %q", index, wantBookmark.Href, got.Root.Bookmarks[index].Href)
				}

				if got.Root.Bookmarks[index].Title != wantBookmark.Title {
					t.Errorf("want bookmark %d title %q, got %q", index, wantBookmark.Title, got.Root.Bookmarks[index].Title)
				}

				if len(got.Root.Bookmarks[index].Attributes) != len(wantBookmark.Attributes) {
					t.Errorf("want %d attributes for bookmark %d, got %d", len(wantBookmark.Attributes), index, len(got.Root.Bookmarks[index].Attributes))
				}
			}
		})
	}
}
