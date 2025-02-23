// Copyright (c) VirtualTam
// SPDX-License-Identifier: MIT

package netscape

import (
	"testing"
	"time"
)

func TestMarshal(t *testing.T) {
	folderCreatedAt := time.Date(2021, time.June, 1, 17, 11, 13, 0, time.UTC)
	folderUpdatedAt := time.Date(2021, time.August, 1, 22, 9, 46, 0, time.UTC)
	bookmarkCreatedAt := time.Date(2022, time.January, 1, 17, 11, 13, 0, time.UTC)
	bookmarkUpdatedAt := time.Date(2022, time.March, 1, 22, 9, 46, 0, time.UTC)

	cases := []struct {
		tname    string
		document Document
		want     string
	}{
		{
			tname: "empty document",
			want: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
     It will be read and overwritten.
     DO NOT EDIT! -->
<TITLE></TITLE>
<H1></H1>
<DL><p>
</DL><p>
`,
		},

		{
			tname: "document with bookmarks",
			document: Document{
				Title: "Bookmarks",
				Root: Folder{
					Name: "Bookmarks",
					Bookmarks: []Bookmark{
						{
							URL:   "https://domain.tld",
							Title: "Test Domain",
						},
						{
							Description: "Second test",
							URL:         "https://test.domain.tld",
							Title:       "Test Domain II",
						},
						{
							Description: `
> The format of here-documents is:

` + "```shell" + `
[n]<<[-]word
        here-document
delimiter
` + "```" + `

> If any part of word is quoted, the delimiter is the result of quote removal on word, and the lines in the here-document are not expanded.`,
							URL:   "https://markdown.xml",
							Title: "Markdown Description",
						},
					},
				},
			},
			want: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
     It will be read and overwritten.
     DO NOT EDIT! -->
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><A HREF="https://domain.tld" PRIVATE="0">Test Domain</A>
    <DT><A HREF="https://test.domain.tld" PRIVATE="0">Test Domain II</A>
    <DD>Second test
    <DT><A HREF="https://markdown.xml" PRIVATE="0">Markdown Description</A>
    <DD>
&gt; The format of here-documents is:

` + "```shell" + `
[n]&lt;&lt;[-]word
        here-document
delimiter
` + "```" + `

&gt; If any part of word is quoted, the delimiter is the result of quote removal on word, and the lines in the here-document are not expanded.
</DL><p>
`,
		},

		{
			tname: "document with private bookmarks and dates",
			document: Document{
				Title: "Bookmarks",
				Root: Folder{
					Name: "Bookmarks",
					Subfolders: []Folder{
						{
							Name:        "Favorites",
							Description: "Add bookmarks here",
							CreatedAt:   folderCreatedAt,
							UpdatedAt:   folderUpdatedAt,
							Bookmarks: []Bookmark{
								{
									CreatedAt: bookmarkCreatedAt,
									URL:       "https://domain.tld",
									Title:     "Test Domain",
									Private:   true,
								},
								{
									CreatedAt:   bookmarkCreatedAt,
									UpdatedAt:   bookmarkUpdatedAt,
									Description: "Second test",
									URL:         "https://test.domain.tld",
									Title:       "Test Domain II",
									Private:     true,
								},
							},
						},
					},
				},
			},
			want: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
     It will be read and overwritten.
     DO NOT EDIT! -->
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><H3 ADD_DATE="1622567473" LAST_MODIFIED="1627855786">Favorites</H3>
    <DD>Add bookmarks here
    <DL><p>
        <DT><A HREF="https://domain.tld" ADD_DATE="1641057073" PRIVATE="1">Test Domain</A>
        <DT><A HREF="https://test.domain.tld" ADD_DATE="1641057073" LAST_MODIFIED="1646172586" PRIVATE="1">Test Domain II</A>
        <DD>Second test
    </DL><p>
</DL><p>
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.tname, func(t *testing.T) {
			m, err := Marshal(&tc.document)

			if err != nil {
				t.Fatalf("expected no error, got %q", err)
			}

			got := string(m)

			if got != tc.want {
				t.Errorf("\nwant:\n%s\n\ngot:\n%s", tc.want, got)
			}
		})
	}
}
