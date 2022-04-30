package decoder

import (
	"testing"
	"time"

	"github.com/virtualtam/netscape-go/ast"
	"github.com/virtualtam/netscape-go/types"
)

func TestDecodeFile(t *testing.T) {
	cases := []struct {
		tname string
		file  ast.FileNode
		want  types.Document
	}{
		{
			tname: "empty document",
		},
		{
			tname: "flat document",
			file: ast.FileNode{
				Title: "Bookmarks",
				Root: ast.FolderNode{
					Name:        "Test Folder",
					Description: "Add bookmarks to this folder",
					Bookmarks: []ast.BookmarkNode{
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
			want: types.Document{
				Title: "Bookmarks",
				Root: types.Folder{
					Name:        "Test Folder",
					Description: "Add bookmarks to this folder",
					Bookmarks: []types.Bookmark{
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
		input ast.FolderNode
		want  types.Folder
	}{
		{
			tname: "empty folder",
			input: ast.FolderNode{
				Name:        "Test Folder",
				Description: "Add bookmarks to this folder",
			},
			want: types.Folder{
				Name:        "Test Folder",
				Description: "Add bookmarks to this folder",
			},
		},
		{
			tname: "empty folder with creation date",
			input: ast.FolderNode{
				Name: "Test Folder",
				Attributes: map[string]string{
					"ADD_DATE": "1646154673",
				},
			},
			want: types.Folder{
				CreatedAt: &folderCreatedAt,
				Name:      "Test Folder",
			},
		},
		{
			tname: "empty folder with creation and update dates, and extra attributes",
			input: ast.FolderNode{
				Name:        "Test Folder",
				Description: "Add bookmarks to this folder",
				Attributes: map[string]string{
					"ADD_DATE":                "1646154673",
					"LAST_MODIFIED":           "1646172586",
					"PERSONAL_TOOLBAR_FOLDER": "true",
				},
			},
			want: types.Folder{
				CreatedAt:   &folderCreatedAt,
				UpdatedAt:   &folderUpdatedAt,
				Name:        "Test Folder",
				Description: "Add bookmarks to this folder",
				Attributes: map[string]string{
					"PERSONAL_TOOLBAR_FOLDER": "true",
				},
			},
		},
		{
			tname: "folder with bookmarks",
			input: ast.FolderNode{
				Name:        "Test Folder",
				Description: "Add bookmarks to this folder",
				Bookmarks: []ast.BookmarkNode{
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
			want: types.Folder{
				Name:        "Test Folder",
				Description: "Add bookmarks to this folder",
				Bookmarks: []types.Bookmark{
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
		{
			tname: "folder with sub-folders and bookmarks",
			input: ast.FolderNode{
				Name:        "Bookmarks",
				Description: "Root Folder",
				Bookmarks: []ast.BookmarkNode{
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
				Subfolders: []ast.FolderNode{
					{
						Name: "Empty",
					},
					{
						Name: "Personal Toolbar",
						Bookmarks: []ast.BookmarkNode{
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
			want: types.Folder{
				Name:        "Bookmarks",
				Description: "Root Folder",
				Bookmarks: []types.Bookmark{
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
				Subfolders: []types.Folder{
					{
						Name: "Empty",
					},
					{
						Name: "Personal Toolbar",
						Bookmarks: []types.Bookmark{
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

func assertFoldersEqual(t *testing.T, got types.Folder, want types.Folder) {
	t.Helper()

	assertDatesEqual(t, "creation", got.CreatedAt, want.CreatedAt)
	assertDatesEqual(t, "update", got.UpdatedAt, want.UpdatedAt)

	if got.Description != want.Description {
		t.Errorf("want description %q, got %q", want.Description, got.Description)
	}

	if got.Name != want.Name {
		t.Errorf("want name %q, got %q", want.Name, got.Name)
	}

	assertAttributesEqual(t, got.Attributes, want.Attributes)

	if len(got.Bookmarks) != len(want.Bookmarks) {
		t.Fatalf("want %d bookmarks, got %d", len(want.Bookmarks), len(got.Bookmarks))
	}

	for index, wantBookmark := range want.Bookmarks {
		assertBookmarksEqual(t, got.Bookmarks[index], wantBookmark)
	}

	if len(got.Subfolders) != len(want.Subfolders) {
		t.Fatalf("want %d subfolders, got %d", len(want.Subfolders), len(got.Subfolders))
	}

	for index, wantSubfolder := range want.Subfolders {
		assertFoldersEqual(t, got.Subfolders[index], wantSubfolder)
	}
}

func TestDecodeBookmark(t *testing.T) {
	bookmarkCreatedAt := time.Date(2022, time.March, 1, 17, 11, 13, 0, time.UTC)
	bookmarkUpdatedAt := time.Date(2022, time.March, 1, 22, 9, 46, 0, time.UTC)

	cases := []struct {
		tname string
		input ast.BookmarkNode
		want  types.Bookmark
	}{
		{
			tname: "bookmark with mandatory information only",
			input: ast.BookmarkNode{
				Href:  "https://domain.tld",
				Title: "Test Domain",
			},
			want: types.Bookmark{
				Title: "Test Domain",
				Href:  "https://domain.tld",
			},
		},
		{
			tname: "bookmark with multi-line description",
			input: ast.BookmarkNode{
				Description: "Nested lists:\n- list1\n  - item1.1\n  - item1.2\n  - item1.3\n- list2\n  - item2.1",
				Href:        "https://domain.tld",
				Title:       "Test Domain",
			},
			want: types.Bookmark{
				Title:       "Test Domain",
				Href:        "https://domain.tld",
				Description: "Nested lists:\n- list1\n  - item1.1\n  - item1.2\n  - item1.3\n- list2\n  - item2.1",
			},
		},
		{
			tname: "bookmark with creation and update date",
			input: ast.BookmarkNode{
				Href:  "https://domain.tld",
				Title: "Test Domain",
				Attributes: map[string]string{
					"ADD_DATE":      "1646154673",
					"LAST_MODIFIED": "1646172586",
				},
			},
			want: types.Bookmark{
				CreatedAt: &bookmarkCreatedAt,
				UpdatedAt: &bookmarkUpdatedAt,
				Title:     "Test Domain",
				Href:      "https://domain.tld",
			},
		},
		{
			tname: "private bookmark",
			input: ast.BookmarkNode{
				Href:  "https://domain.tld",
				Title: "Test Domain",
				Attributes: map[string]string{
					"PRIVATE": "1",
				},
			},
			want: types.Bookmark{
				Title:   "Test Domain",
				Href:    "https://domain.tld",
				Private: true,
			},
		},
		{
			tname: "bookmark with comma-separated tags and extra whitespace",
			input: ast.BookmarkNode{
				Href:  "https://domain.tld",
				Title: "Test Domain",
				Attributes: map[string]string{
					"TAGS": "test, netscape,     bookmark",
				},
			},
			want: types.Bookmark{
				Title: "Test Domain",
				Href:  "https://domain.tld",
				Tags: []string{
					"bookmark",
					"netscape",
					"test",
				},
			},
		},
		{
			tname: "bookmark with extra attributes",
			input: ast.BookmarkNode{
				Href:  "https://domain.tld",
				Title: "Test Domain",
				Attributes: map[string]string{
					"ICON_URI":     "https://domain.tld/favicon.ico",
					"LAST_CHARSET": "windows-1252",
					"PRIVATE":      "1",
				},
			},
			want: types.Bookmark{
				Title:   "Test Domain",
				Href:    "https://domain.tld",
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

func assertBookmarksEqual(t *testing.T, got types.Bookmark, want types.Bookmark) {
	assertDatesEqual(t, "creation", got.CreatedAt, want.CreatedAt)
	assertDatesEqual(t, "update", got.UpdatedAt, want.UpdatedAt)

	if got.Title != want.Title {
		t.Errorf("want title %q, got %q", want.Title, got.Title)
	}

	if got.Href != want.Href {
		t.Errorf("want URL string %q, got %q", want.Href, got.Href)
	}

	if got.Description != want.Description {
		t.Errorf("want description %q, got %q", want.Description, got.Description)
	}

	if got.Private != want.Private {
		t.Errorf("want private %t, got %t", want.Private, got.Private)
	}

	if len(got.Tags) != len(want.Tags) {
		t.Fatalf("want %d tags, got %d", len(want.Tags), len(got.Tags))
	}

	for index, wantTag := range want.Tags {
		if got.Tags[index] != wantTag {
			t.Errorf("want tag %d value %q, got %q", index, wantTag, got.Tags[index])
		}
	}

	assertAttributesEqual(t, got.Attributes, want.Attributes)
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

func assertAttributesEqual(t *testing.T, got map[string]string, want map[string]string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("want %d attributes, got %d", len(want), len(got))
	}

	for attr, wantValue := range want {
		if got[attr] != wantValue {
			t.Errorf("want attribute %q value %q, got %q", attr, wantValue, got[attr])
		}
	}
}

func assertDatesEqual(t *testing.T, name string, got *time.Time, want *time.Time) {
	t.Helper()

	if want == nil {
		if got != nil {
			t.Errorf("want %s datetime nil, got %q", name, got.String())
		}
		return
	}

	if got == nil {
		t.Errorf("want %s datetime %q, got nil", name, want.String())
		return
	}

	if got.String() != want.String() {
		t.Errorf("want %s date %q, got %q", name, want.String(), got.String())
	}
}
