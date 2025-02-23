// Copyright (c) VirtualTam
// SPDX-License-Identifier: MIT

package netscape

import (
	"testing"
	"time"
)

func TestDocumentFlatten(t *testing.T) {
	cases := []struct {
		tname    string
		document Document
		want     Document
	}{
		{
			tname: "empty",
		},
		{
			tname: "flat",
			document: Document{
				Title: "Already Flat",
				Root: Folder{
					Name: "Already Flat",
					Bookmarks: []Bookmark{
						{
							Title: "Flat 1",
							URL:   "https://flat1.domain.tld",
						},
						{
							Title: "Flat 2",
							URL:   "https://flat2.domain.tld",
						},
					},
				},
			},
			want: Document{
				Title: "Already Flat",
				Root: Folder{
					Name: "Already Flat",
					Bookmarks: []Bookmark{
						{
							Title: "Flat 1",
							URL:   "https://flat1.domain.tld",
						},
						{
							Title: "Flat 2",
							URL:   "https://flat2.domain.tld",
						},
					},
				},
			},
		},
		{
			tname: "nested",
			document: Document{
				Title: "Nested",
				Root: Folder{
					Name: "Nested",
					Bookmarks: []Bookmark{
						{
							Title: "Nested 1",
							URL:   "https://n1.domain.tld",
						},
						{
							Title: "Nested 2",
							URL:   "https://n2.domain.tld",
						},
					},
					Subfolders: []Folder{
						{
							Name: "Subfolder A",
							Bookmarks: []Bookmark{
								{
									Title: "Nested A1",
									URL:   "https://na1.domain.tld",
								},
								{
									Title: "Nested A2",
									URL:   "https://na2.domain.tld",
								},
							},
						},
						{
							Name: "Subfolder B",
							Bookmarks: []Bookmark{
								{
									Title: "Nested B1",
									URL:   "https://nb1.domain.tld",
								},
							},
							Subfolders: []Folder{
								{
									Name: "Subfolder B-1",
									Bookmarks: []Bookmark{
										{
											Title: "Nested B-1.1",
											URL:   "https://nb1-1.domain.tld",
										},
										{
											Title: "Nested B-1.2",
											URL:   "https://nb1-2.domain.tld",
										},
									},
								},
								{
									Name: "Subfolder B-2",
									Bookmarks: []Bookmark{
										{
											Title: "Nested B-2.1",
											URL:   "https://nb2-1.domain.tld",
										},
										{
											Title: "Nested B-2.2",
											URL:   "https://nb2-2.domain.tld",
										},
									},
								},
							},
						},
					},
				},
			},
			want: Document{
				Title: "Nested",
				Root: Folder{
					Name: "Nested",
					Bookmarks: []Bookmark{
						{
							Title: "Nested 1",
							URL:   "https://n1.domain.tld",
						},
						{
							Title: "Nested 2",
							URL:   "https://n2.domain.tld",
						},
						{
							Title: "Nested A1",
							URL:   "https://na1.domain.tld",
						},
						{
							Title: "Nested A2",
							URL:   "https://na2.domain.tld",
						},
						{
							Title: "Nested B1",
							URL:   "https://nb1.domain.tld",
						},
						{
							Title: "Nested B-1.1",
							URL:   "https://nb1-1.domain.tld",
						},
						{
							Title: "Nested B-1.2",
							URL:   "https://nb1-2.domain.tld",
						},
						{
							Title: "Nested B-2.1",
							URL:   "https://nb2-1.domain.tld",
						},
						{
							Title: "Nested B-2.2",
							URL:   "https://nb2-2.domain.tld",
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.tname, func(t *testing.T) {
			got := tc.document.Flatten()

			if got.Title != tc.want.Title {
				t.Errorf("want title %q, got %q", tc.want.Title, got.Title)
			}

			if len(got.Root.Subfolders) > 0 {
				t.Errorf("want no subfolders, got %d", len(got.Root.Subfolders))
			}

			assertFoldersEqual(t, got.Root, tc.want.Root)
		})
	}
}

func assertFoldersEqual(t *testing.T, got Folder, want Folder) {
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

func assertBookmarksEqual(t *testing.T, got Bookmark, want Bookmark) {
	assertDatesEqual(t, "creation", got.CreatedAt, want.CreatedAt)
	assertDatesEqual(t, "update", got.UpdatedAt, want.UpdatedAt)

	if got.Title != want.Title {
		t.Errorf("want title %q, got %q", want.Title, got.Title)
	}

	if got.URL != want.URL {
		t.Errorf("want URL string %q, got %q", want.URL, got.URL)
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

func assertDatesEqual(t *testing.T, name string, got, want time.Time) {
	t.Helper()

	if !got.Equal(want) {
		t.Errorf("want %s date %q, got %q", name, want.String(), got.String())
	}
}
