package decoder

import (
	"net/url"
	"testing"
	"time"

	"github.com/virtualtam/netscape-go"
	"github.com/virtualtam/netscape-go/ast"
)

func TestDecodeBookmark(t *testing.T) {
	cases := []struct {
		tname string
		input ast.Bookmark
		want  netscape.Bookmark
	}{
		{
			tname: "bookmark with mandatory information only",
			input: ast.Bookmark{
				Href:  "https://domain.tld",
				Title: "Test Domain",
			},
			want: netscape.Bookmark{
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
			want: netscape.Bookmark{
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
			want: netscape.Bookmark{
				CreatedAt: time.Date(2022, time.March, 1, 17, 11, 13, 0, time.UTC),
				UpdatedAt: time.Date(2022, time.March, 1, 22, 9, 46, 0, time.UTC),
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
			want: netscape.Bookmark{
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
			want: netscape.Bookmark{
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
			want: netscape.Bookmark{
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
			got, err := decodeBookmark(tc.input)

			if err != nil {
				t.Errorf("expected no error, got %q", err)
			}

			if got.CreatedAt != tc.want.CreatedAt {
				t.Errorf("want creation date %q, got %q", tc.want.CreatedAt.String(), got.CreatedAt.String())
			}

			if got.UpdatedAt != tc.want.UpdatedAt {
				t.Errorf("want update date %q, got %q", tc.want.UpdatedAt.String(), got.UpdatedAt.String())
			}

			if got.Title != tc.want.Title {
				t.Errorf("want title %q, got %q", tc.want.Title, got.Title)
			}

			if got.URL != tc.want.URL {
				t.Errorf("want URL %q, got %q", tc.want.URL.String(), got.URL.String())
			}

			if got.Description != tc.want.Description {
				t.Errorf("want description %q, got %q", tc.want.Description, got.Description)
			}

			if got.Private != tc.want.Private {
				t.Errorf("want private %t, got %t", tc.want.Private, got.Private)
			}

			if len(got.Tags) != len(tc.want.Tags) {
				t.Errorf("want %d tags, got %d", len(tc.want.Tags), len(got.Tags))
				return
			}

			for index, wantTag := range tc.want.Tags {
				if got.Tags[index] != wantTag {
					t.Errorf("want tag %d value %q, got %q", index, wantTag, got.Tags[index])
				}
			}

			if len(got.Attributes) != len(tc.want.Attributes) {
				t.Errorf("want %d attributes, got %d", len(tc.want.Attributes), len(got.Attributes))
				return
			}

			for attr, wantValue := range tc.want.Attributes {
				if got.Attributes[attr] != wantValue {
					t.Errorf("want attribute %q value %q, got %q", attr, wantValue, got.Attributes[attr])
				}
			}
		})
	}
}
