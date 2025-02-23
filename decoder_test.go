// Copyright (c) VirtualTam
// SPDX-License-Identifier: MIT

package netscape

import (
	"testing"
	"time"
)

func TestDecodeFile(t *testing.T) {
	cases := []struct {
		tname string
		file  FileNode
		want  Document
	}{
		{
			tname: "empty document",
		},
		{
			tname: "flat document",
			file: FileNode{
				Title: "Bookmarks",
				Root: FolderNode{
					Name:        "Test Folder",
					Description: "Add bookmarks to this folder",
					Bookmarks: []BookmarkNode{
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
			want: Document{
				Title: "Bookmarks",
				Root: Folder{
					Name:        "Test Folder",
					Description: "Add bookmarks to this folder",
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
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.tname, func(t *testing.T) {
			d := NewDecoder()

			got, err := d.decodeFile(tc.file)

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

func TestDecodeFolder(t *testing.T) {
	folderCreatedAt := time.Date(2022, time.March, 1, 17, 11, 13, 0, time.UTC)
	folderUpdatedAt := time.Date(2022, time.March, 1, 22, 9, 46, 0, time.UTC)

	cases := []struct {
		tname string
		input FolderNode
		want  Folder
	}{
		{
			tname: "empty folder",
			input: FolderNode{
				Name:        "Test Folder",
				Description: "Add bookmarks to this folder",
			},
			want: Folder{
				Name:        "Test Folder",
				Description: "Add bookmarks to this folder",
			},
		},
		{
			tname: "empty folder with creation date",
			input: FolderNode{
				Name: "Test Folder",
				Attributes: map[string]string{
					"ADD_DATE": "1646154673",
				},
			},
			want: Folder{
				CreatedAt: folderCreatedAt,
				UpdatedAt: folderCreatedAt,
				Name:      "Test Folder",
			},
		},
		{
			tname: "empty folder with creation and update dates, and extra attributes",
			input: FolderNode{
				Name:        "Test Folder",
				Description: "Add bookmarks to this folder",
				Attributes: map[string]string{
					"ADD_DATE":                "1646154673",
					"LAST_MODIFIED":           "1646172586",
					"PERSONAL_TOOLBAR_FOLDER": "true",
				},
			},
			want: Folder{
				CreatedAt:   folderCreatedAt,
				UpdatedAt:   folderUpdatedAt,
				Name:        "Test Folder",
				Description: "Add bookmarks to this folder",
				Attributes: map[string]string{
					"PERSONAL_TOOLBAR_FOLDER": "true",
				},
			},
		},
		{
			tname: "folder with bookmarks",
			input: FolderNode{
				Name:        "Test Folder",
				Description: "Add bookmarks to this folder",
				Bookmarks: []BookmarkNode{
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
			want: Folder{
				Name:        "Test Folder",
				Description: "Add bookmarks to this folder",
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
				},
			},
		},
		{
			tname: "folder with sub-folders and bookmarks",
			input: FolderNode{
				Name:        "Bookmarks",
				Description: "Root Folder",
				Bookmarks: []BookmarkNode{
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
				Subfolders: []FolderNode{
					{
						Name: "Empty",
					},
					{
						Name: "Personal Toolbar",
						Bookmarks: []BookmarkNode{
							{
								Href:  "https://personal.tld",
								Title: "Personal Domain",
							},
							{
								Description: "Weather Reports",
								Href:        "https://weather.tld",
								Title:       "Weather Reports",
							},
						},
					},
				},
			},
			want: Folder{
				Name:        "Bookmarks",
				Description: "Root Folder",
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
				},
				Subfolders: []Folder{
					{
						Name: "Empty",
					},
					{
						Name: "Personal Toolbar",
						Bookmarks: []Bookmark{
							{
								URL:   "https://personal.tld",
								Title: "Personal Domain",
							},
							{
								Description: "Weather Reports",
								URL:         "https://weather.tld",
								Title:       "Weather Reports",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.tname, func(t *testing.T) {
			d := NewDecoder()

			got, err := d.decodeFolder(tc.input)

			if err != nil {
				t.Fatalf("expected no error, got %q", err)
			}

			assertFoldersEqual(t, got, tc.want)
		})
	}
}

func TestDecodeBookmark(t *testing.T) {
	bookmarkCreatedAt := time.Date(2022, time.March, 1, 17, 11, 13, 0, time.UTC)
	bookmarkUpdatedAt := time.Date(2022, time.March, 1, 22, 9, 46, 0, time.UTC)

	cases := []struct {
		tname string
		input BookmarkNode
		want  Bookmark
	}{
		{
			tname: "bookmark with mandatory information only",
			input: BookmarkNode{
				Href:  "https://domain.tld",
				Title: "Test Domain",
			},
			want: Bookmark{
				Title: "Test Domain",
				URL:   "https://domain.tld",
			},
		},
		{
			tname: "bookmark with multi-line description",
			input: BookmarkNode{
				Description: "Nested lists:\n- list1\n  - item1.1\n  - item1.2\n  - item1.3\n- list2\n  - item2.1",
				Href:        "https://domain.tld",
				Title:       "Test Domain",
			},
			want: Bookmark{
				Title:       "Test Domain",
				URL:         "https://domain.tld",
				Description: "Nested lists:\n- list1\n  - item1.1\n  - item1.2\n  - item1.3\n- list2\n  - item2.1",
			},
		},
		{
			tname: "bookmark with description containing escaped HTML characters",
			input: BookmarkNode{
				Description: "&#34;Fran &amp; Freddie&#39;s Diner&#34; &lt;tasty@example.com&gt;",
				Href:        "https://domain.tld",
				Title:       "Test Domain",
			},
			want: Bookmark{
				Title:       "Test Domain",
				URL:         "https://domain.tld",
				Description: `"Fran & Freddie's Diner" <tasty@example.com>`,
			},
		},
		{
			tname: "bookmark with multi-line description containing escaped HTML characters",
			input: BookmarkNode{
				Description: `
&gt; The format of here-documents is:

` + "```shell" + `
[n]&lt;&lt;[-]word
        here-document
delimiter
` + "```" + `

&gt; If any part of word is quoted, the delimiter is the result of quote removal on word, and the lines in the here-document are not expanded.`,
				Href:  "https://domain.tld",
				Title: "Test Domain",
			},
			want: Bookmark{
				Title: "Test Domain",
				URL:   "https://domain.tld",
				Description: `
> The format of here-documents is:

` + "```shell" + `
[n]<<[-]word
        here-document
delimiter
` + "```" + `

> If any part of word is quoted, the delimiter is the result of quote removal on word, and the lines in the here-document are not expanded.`,
			},
		},
		{
			tname: "bookmark with creation date",
			input: BookmarkNode{
				Href:  "https://domain.tld",
				Title: "Test Domain",
				Attributes: map[string]string{
					"ADD_DATE": "1646154673",
				},
			},
			want: Bookmark{
				CreatedAt: bookmarkCreatedAt,
				UpdatedAt: bookmarkCreatedAt,
				Title:     "Test Domain",
				URL:       "https://domain.tld",
			},
		},
		{
			tname: "bookmark with creation and update date",
			input: BookmarkNode{
				Href:  "https://domain.tld",
				Title: "Test Domain",
				Attributes: map[string]string{
					"ADD_DATE":      "1646154673",
					"LAST_MODIFIED": "1646172586",
				},
			},
			want: Bookmark{
				CreatedAt: bookmarkCreatedAt,
				UpdatedAt: bookmarkUpdatedAt,
				Title:     "Test Domain",
				URL:       "https://domain.tld",
			},
		},
		{
			tname: "private bookmark",
			input: BookmarkNode{
				Href:  "https://domain.tld",
				Title: "Test Domain",
				Attributes: map[string]string{
					"PRIVATE": "1",
				},
			},
			want: Bookmark{
				Title:   "Test Domain",
				URL:     "https://domain.tld",
				Private: true,
			},
		},
		{
			tname: "bookmark with comma-separated tags and extra whitespace",
			input: BookmarkNode{
				Href:  "https://domain.tld",
				Title: "Test Domain",
				Attributes: map[string]string{
					"TAGS": "test, netscape,     bookmark",
				},
			},
			want: Bookmark{
				Title: "Test Domain",
				URL:   "https://domain.tld",
				Tags: []string{
					"bookmark",
					"netscape",
					"test",
				},
			},
		},
		{
			tname: "bookmark with extra attributes",
			input: BookmarkNode{
				Href:  "https://domain.tld",
				Title: "Test Domain",
				Attributes: map[string]string{
					"ICON_URI":     "https://domain.tld/favicon.ico",
					"LAST_CHARSET": "windows-1252",
					"PRIVATE":      "1",
				},
			},
			want: Bookmark{
				Title:   "Test Domain",
				URL:     "https://domain.tld",
				Private: true,
				Attributes: map[string]string{
					"ICON_URI":     "https://domain.tld/favicon.ico",
					"LAST_CHARSET": "windows-1252",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.tname, func(t *testing.T) {
			d := NewDecoder()

			got, err := d.decodeBookmark(tc.input)

			if err != nil {
				t.Errorf("expected no error, got %q", err)
			}

			assertBookmarksEqual(t, got, tc.want)
		})
	}
}

func TestDecodeDateTime(t *testing.T) {
	cases := []struct {
		tname string
		input string
		want  time.Time
	}{
		// UNIX time
		{
			// date +%s
			tname: "UNIX epoch",
			input: "1646154673",
			want:  time.Date(2022, time.March, 1, 17, 11, 13, 0, time.UTC),
		},
		{
			// date +%s%3N
			tname: "UNIX epoch (milliseconds)",
			input: "1646155662212",
			want:  time.Date(2022, time.March, 1, 17, 27, 42, 212000000, time.UTC),
		},
		{
			// date +%s%6N
			tname: "UNIX epoch (microseconds)",
			input: "1646156161974685",
			want:  time.Date(2022, time.March, 1, 17, 36, 01, 974685000, time.UTC),
		},
		{
			// date +%s%9N
			tname: "UNIX epoch (nanoseconds)",
			input: "1646156260253353101",
			want:  time.Date(2022, time.March, 1, 17, 37, 40, 253353101, time.UTC),
		},

		// String representations
		{
			// date --rfc-3339=seconds
			tname: "RFC3339",
			input: "2022-03-01T18:54:13+01:00",
			want:  time.Date(2022, time.March, 1, 17, 54, 13, 0, time.UTC),
		},
		{
			// date --rfc-3339=seconds
			tname: "RFC3339 (nanoseconds)",
			input: "2022-03-01T18:54:30.585063231+01:00",
			want:  time.Date(2022, time.March, 1, 17, 54, 30, 585063231, time.UTC),
		},
		{
			// date "+%d/%b/%Y:%H:%M:%S %z"
			tname: "Common Log",
			input: "10/Oct/2000:13:55:36 -0700",
			want:  time.Date(2000, 10, 10, 20, 55, 36, 0, time.UTC),
		},
	}

	for _, tc := range cases {
		t.Run(tc.tname, func(t *testing.T) {
			d := NewDecoder()

			got, err := d.decodeDate(tc.input)

			if err != nil {
				t.Fatalf("expected no error, got %q", err)
			}

			if got != tc.want {
				t.Errorf("want date/time %q, got %q", tc.want, got)
			}
		})
	}
}
