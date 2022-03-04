package decoder

import (
	"net/url"
	"testing"
	"time"

	"github.com/virtualtam/netscape-go/ast"
	"github.com/virtualtam/netscape-go/types"
)

func TestDecodeFile(t *testing.T) {
	cases := []struct {
		tname string
		file  ast.File
		want  types.Document
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
			want: types.Document{
				Title: "Bookmarks",
				Root: types.Folder{
					Name:        "Test Folder",
					Description: "Add bookmarks to this folder",
					Bookmarks: []types.Bookmark{
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
			var d Decoder
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
		input ast.Folder
		want  types.Folder
	}{
		{
			tname: "empty folder",
			input: ast.Folder{
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
			input: ast.Folder{
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
			input: ast.Folder{
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
			input: ast.Folder{
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
			want: types.Folder{
				Name:        "Test Folder",
				Description: "Add bookmarks to this folder",
				Bookmarks: []types.Bookmark{
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
		{
			tname: "folder with sub-folders and bookmarks",
			input: ast.Folder{
				Name:        "Bookmarks",
				Description: "Root Folder",
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
				Subfolders: []ast.Folder{
					{
						Name: "Empty",
					},
					{
						Name: "Personal Toolbar",
						Bookmarks: []ast.Bookmark{
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
				Subfolders: []types.Folder{
					{
						Name: "Empty",
					},
					{
						Name: "Personal Toolbar",
						Bookmarks: []types.Bookmark{
							{
								URL: url.URL{
									Scheme: "https",
									Host:   "personal.tld",
								},
								Title: "Personal Domain",
							},
							{
								Description: "Weather Reports",
								URL: url.URL{
									Scheme: "https",
									Host:   "weather.tld",
								},
								Title: "Weather Reports",
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.tname, func(t *testing.T) {
			var d Decoder
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
		input ast.Bookmark
		want  types.Bookmark
	}{
		{
			tname: "bookmark with mandatory information only",
			input: ast.Bookmark{
				Href:  "https://domain.tld",
				Title: "Test Domain",
			},
			want: types.Bookmark{
				Title: "Test Domain",
				URL: url.URL{
					Scheme: "https",
					Host:   "domain.tld",
				},
			},
		},
		{
			tname: "bookmark with multi-line description",
			input: ast.Bookmark{
				Description: "Nested lists:\n- list1\n  - item1.1\n  - item1.2\n  - item1.3\n- list2\n  - item2.1",
				Href:        "https://domain.tld",
				Title:       "Test Domain",
			},
			want: types.Bookmark{
				Title: "Test Domain",
				URL: url.URL{
					Scheme: "https",
					Host:   "domain.tld",
				},
				Description: "Nested lists:\n- list1\n  - item1.1\n  - item1.2\n  - item1.3\n- list2\n  - item2.1",
			},
		},
		{
			tname: "bookmark with creation and update date",
			input: ast.Bookmark{
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
				URL: url.URL{
					Scheme: "https",
					Host:   "domain.tld",
				},
			},
		},
		{
			tname: "private bookmark",
			input: ast.Bookmark{
				Href:  "https://domain.tld",
				Title: "Test Domain",
				Attributes: map[string]string{
					"PRIVATE": "1",
				},
			},
			want: types.Bookmark{
				Title: "Test Domain",
				URL: url.URL{
					Scheme: "https",
					Host:   "domain.tld",
				},
				Private: true,
			},
		},
		{
			tname: "bookmark with comma-separated tags and extra whitespace",
			input: ast.Bookmark{
				Href:  "https://domain.tld",
				Title: "Test Domain",
				Attributes: map[string]string{
					"TAGS": "test, netscape,     bookmark",
				},
			},
			want: types.Bookmark{
				Title: "Test Domain",
				URL: url.URL{
					Scheme: "https",
					Host:   "domain.tld",
				},
				Tags: []string{
					"bookmark",
					"netscape",
					"test",
				},
			},
		},
		{
			tname: "bookmark with extra attributes",
			input: ast.Bookmark{
				Href:  "https://domain.tld",
				Title: "Test Domain",
				Attributes: map[string]string{
					"ICON_URI":     "https://domain.tld/favicon.ico",
					"LAST_CHARSET": "windows-1252",
					"PRIVATE":      "1",
				},
			},
			want: types.Bookmark{
				Title: "Test Domain",
				URL: url.URL{
					Scheme: "https",
					Host:   "domain.tld",
				},
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
			var d Decoder
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

	if got.URL != want.URL {
		t.Errorf("want URL %q, got %q", want.URL.String(), got.URL.String())
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
