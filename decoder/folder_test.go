package decoder

import (
	"net/url"
	"testing"
	"time"

	"github.com/virtualtam/netscape-go"
	"github.com/virtualtam/netscape-go/ast"
)

func TestDecodeFolder(t *testing.T) {
	cases := []struct {
		tname string
		input ast.Folder
		want  netscape.Folder
	}{
		{
			tname: "empty folder",
			input: ast.Folder{
				Name:        "Test Folder",
				Description: "Add bookmarks to this folder",
			},
			want: netscape.Folder{
				Name:        "Test Folder",
				Description: "Add bookmarks to this folder",
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
			want: netscape.Folder{
				CreatedAt:   time.Date(2022, time.March, 1, 17, 11, 13, 0, time.UTC),
				UpdatedAt:   time.Date(2022, time.March, 1, 22, 9, 46, 0, time.UTC),
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
			want: netscape.Folder{
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
			want: netscape.Folder{
				Name:        "Bookmarks",
				Description: "Root Folder",
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
				Subfolders: []netscape.Folder{
					{
						Name: "Empty",
					},
					{
						Name: "Personal Toolbar",
						Bookmarks: []netscape.Bookmark{
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
			got, err := decodeFolder(tc.input)

			if err != nil {
				t.Errorf("expected no error, got %q", err)
				return
			}

			assertFoldersEqual(t, got, tc.want)
		})
	}
}

func assertFoldersEqual(t *testing.T, got netscape.Folder, want netscape.Folder) {
	t.Helper()

	if got.CreatedAt != want.CreatedAt {
		t.Errorf("want creation date %q, got %q", want.CreatedAt.String(), got.CreatedAt.String())
	}

	if got.UpdatedAt != want.UpdatedAt {
		t.Errorf("want update date %q, got %q", want.UpdatedAt.String(), got.UpdatedAt.String())
	}

	if got.Description != want.Description {
		t.Errorf("want description %q, got %q", want.Description, got.Description)
	}

	if got.Name != want.Name {
		t.Errorf("want name %q, got %q", want.Name, got.Name)
	}

	assertAttributesEqual(t, got.Attributes, want.Attributes)

	if len(got.Bookmarks) != len(want.Bookmarks) {
		t.Errorf("want %d bookmarks, got %d", len(want.Bookmarks), len(got.Bookmarks))
		return
	}

	for index, wantBookmark := range want.Bookmarks {
		assertBookmarksEqual(t, got.Bookmarks[index], wantBookmark)
	}

	if len(got.Subfolders) != len(want.Subfolders) {
		t.Errorf("want %d subfolders, got %d", len(want.Subfolders), len(got.Subfolders))
		return
	}

	for index, wantSubfolder := range want.Subfolders {
		assertFoldersEqual(t, got.Subfolders[index], wantSubfolder)
	}
}