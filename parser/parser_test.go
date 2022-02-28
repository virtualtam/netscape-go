package parser

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/virtualtam/netscape-go/ast"
)

func TestParse(t *testing.T) {
	cases := []struct {
		tname   string
		input   string
		want    ast.File
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
			want: ast.File{
				Title: "Bookmarks",
				Root: ast.Folder{
					Name: "Bookmarks",
					Bookmarks: []ast.Bookmark{
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
			want: ast.File{
				Title: "Bookmarks",
				Root: ast.Folder{
					Name: "Bookmarks",
					Bookmarks: []ast.Bookmark{
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
		{
			tname: "bookmark with empty description",
			input: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
<DT><A HREF="https://domain.tld">Test Domain</A>
<DD>
<DT><A HREF="https://domain.tld">Test Domain</A>
<DD>
<DT><A HREF="https://emptydesc.domain.tld">Test Domain (with empty description)</A>
<DD>
</DL><p>
`,
			want: ast.File{
				Title: "Bookmarks",
				Root: ast.Folder{
					Name: "Bookmarks",
					Bookmarks: []ast.Bookmark{
						{
							Href:  "https://domain.tld",
							Title: "Test Domain",
						},
						{
							Href:  "https://domain.tld",
							Title: "Test Domain",
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
			tname: "bookmark with multi-line description",
			input: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
<DT><A HREF="https://domain.tld">Test Domain</A>
<DD>Description:

- item 1
    - item 1.1
    - item 1.2
- item 2
- item 3
</DL><p>
`,
			want: ast.File{
				Title: "Bookmarks",
				Root: ast.Folder{
					Name: "Bookmarks",
					Bookmarks: []ast.Bookmark{
						{
							Description: "Description:\n\n- item 1\n    - item 1.1\n    - item 1.2\n- item 2\n- item 3",
							Href:        "https://domain.tld",
							Title:       "Test Domain",
						},
					},
				},
			},
		},
		{
			tname: "bookmark with description containing HTML markup",
			input: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
<DT><A HREF="https://domain.tld">Test Domain</A>
<DD>Markup:
<a href="http://localhost:8080"><img src="http://localhost:8080/splash.png"/></a>
</DL><p>
`,
			want: ast.File{
				Title: "Bookmarks",
				Root: ast.Folder{
					Name: "Bookmarks",
					Bookmarks: []ast.Bookmark{
						{
							Description: `Markup:
<a href="http://localhost:8080"><img src="http://localhost:8080/splash.png"/></a>`,
							Href:  "https://domain.tld",
							Title: "Test Domain",
						},
					},
				},
			},
		},
		{
			tname: "nested folders",
			input: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<TITLE>Bookmarks</TITLE>
<H1>Level 0</H1>
<DL><p>
    <DT><H3>Level 1A</H3>
	<DL><p>
        <DT><H3>Level 2A</H3>
	    <DL><p>
	    </DL><p>
	</DL><p>
    <DT><H3>Level 1B</H3>
	<DL><p>
	</DL><p>
</DL><p>
`,
			want: ast.File{
				Title: "Bookmarks",
				Root: ast.Folder{
					Name: "Level 0",
					Subfolders: []ast.Folder{
						{
							Name: "Level 1A",
							Subfolders: []ast.Folder{
								{Name: "Level 2A"},
							},
						},
						{Name: "Level 1B"},
					},
				},
			},
		},
		{
			tname: "nested folders with bookmarks",
			input: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<TITLE>Bookmarks</TITLE>
<H1>Level 0</H1>
<DL><p>
    <DT><A HREF="https://l0.domain.tld">Level 0</A>
    <DT><A HREF="https://l0.domain.tld">Level 0</A>
    <DT><H3>Level 1A</H3>
	<DL><p>
		<DT><A HREF="https://l1a.domain.tld">Level 1A</A>
        <DT><H3>Level 2A</H3>
	    <DL><p>
		    <DT><A HREF="https://l2a.domain.tld">Level 2A</A>
	    </DL><p>
	</DL><p>
    <DT><H3>Level 1B</H3>
	<DL><p>
	</DL><p>
</DL><p>
`,
			want: ast.File{
				Title: "Bookmarks",
				Root: ast.Folder{
					Name: "Level 0",
					Bookmarks: []ast.Bookmark{
						{Href: "https://l0.domain.tld", Title: "Level 0"},
						{Href: "https://l0.domain.tld", Title: "Level 0"},
					},
					Subfolders: []ast.Folder{
						{
							Name: "Level 1A",
							Bookmarks: []ast.Bookmark{
								{Href: "https://l1a.domain.tld", Title: "Level 1A"},
							},
							Subfolders: []ast.Folder{
								{
									Name: "Level 2A",
									Bookmarks: []ast.Bookmark{
										{Href: "https://l2a.domain.tld", Title: "Level 2A"},
									},
								},
							},
						},
						{Name: "Level 1B"},
					},
				},
			},
		},
		{
			tname: "nested folder with attributes",
			input: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><H3 ADD_DATE="1460294955" LAST_MODIFIED="1460294956" PERSONAL_TOOLBAR_FOLDER="true">Personal toolbar</H3>
	<DD>Add bookmarks to this folder to see them displayed on the Bookmarks Toolbar
	<DL><p>
	</DL><p>
</DL><p>
`,
			want: ast.File{
				Title: "Bookmarks",
				Root: ast.Folder{
					Name: "Bookmarks",
					Subfolders: []ast.Folder{
						{
							Name:        "Personal toolbar",
							Description: "Add bookmarks to this folder to see them displayed on the Bookmarks Toolbar",
							Attributes: map[string]string{
								"ADD_DATE":                "1460294955",
								"LAST_MODIFIED":           "1460294956",
								"PERSONAL_TOOLBAR_FOLDER": "true",
							},
						},
					},
				},
			},
		},
		{
			tname: "nested folder with description and bookmarks",
			input: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<TITLE>Bookmarks</TITLE>
<H1>Level 0</H1>
<DL><p>
    <DT><H3>Level 1</H3>
    <DD>Folder with description
    <DL><p>
        <DT><A HREF="https://domain.tld">Test Domain</A>
    </DL><p>
</DL><p>
`,
			want: ast.File{
				Title: "Bookmarks",
				Root: ast.Folder{
					Name: "Level 0",
					Subfolders: []ast.Folder{
						{
							Name:        "Level 1",
							Description: "Folder with description",
							Bookmarks: []ast.Bookmark{
								{
									Href:  "https://domain.tld",
									Title: "Test Domain",
								},
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
				return
			}

			if got.Title != tc.want.Title {
				t.Errorf("want title %q, got %q", tc.want.Title, got.Title)
			}

			assertFoldersEqual(t, got.Root, tc.want.Root)
		})
	}
}

func TestParseFile(t *testing.T) {
	cases := []struct {
		tname         string
		inputFilename string
		want          ast.File
	}{
		{
			tname:         "Netscape (basic)",
			inputFilename: "netscape_basic.htm",
			want: ast.File{
				Title: "Bookmarks",
				Root: ast.Folder{
					Name: "Bookmarks",
					Bookmarks: []ast.Bookmark{
						{
							Description: "Super-secret stuff you're not supposed to know about",
							Href:        "https://private.tld",
							Title:       "Secret stuff",
							Attributes: map[string]string{
								"ADD_DATE": "10/Oct/2000:13:55:36 +0300",
								"PRIVATE":  "1",
								"TAGS":     "private secret",
							},
						},
						{
							Href:  "http://public.tld",
							Title: "Public stuff",
							Attributes: map[string]string{
								"ADD_DATE": "1456433748",
								"PRIVATE":  "0",
								"TAGS":     "public hello world",
							},
						},
					},
				},
			},
		},

		{
			tname:         "Netscape (extended markup)",
			inputFilename: "netscape_extended.htm",
			want: ast.File{
				Title: "My local links",
				Root: ast.Folder{
					Name: "Shaarli export of all bookmarks on Sat, 06 Jun 20 15:50:59 +0200",
					Bookmarks: []ast.Bookmark{
						{
							Description: `For 10 years, a rogue fishing vessel and its crew plundered the world’s oceans, escaping repeated attempts of capture. Then a dramatic pursuit finally netted the one that got away.
<a href="http://localhost.localdomain:8083/Shaarli/?JVvqCA"><img src="http://localhost.localdomain:8083/Shaarli/cache/thumb/290ccda0deea6083ee613d358446103e/c975558ad43acdbd982ffafd8c01163d6c9ec5ca125901.jpg"/></a>`,
							Href:  "https://www.bbc.com/future/article/20190213-the-dramatic-hunt-for-the-fish-pirates-exploiting-our-seas",
							Title: "The hunt for the fish pirates who exploit the sea - BBC Future",
							Attributes: map[string]string{
								"ADD_DATE": "1591451445",
								"PRIVATE":  "1",
								"TAGS":     "story,oceans",
							},
						},
					},
				},
			},
		},

		{
			tname:         "Netscape (multiline descriptions)",
			inputFilename: "netscape_multiline.htm",
			want: ast.File{
				Title: "Bookmarks",
				Root: ast.Folder{
					Name: "Bookmarks",
					Bookmarks: []ast.Bookmark{
						{
							Description: "List:\n- item1\n- item2\n- item3",
							Href:        "http://multi.li.ne/1",
							Title:       "Multiline desc",
							Attributes: map[string]string{
								"ADD_DATE": "1456433741",
								"PRIVATE":  "0",
								"TAGS":     "multi",
							},
						},
						{
							Description: "Nested lists:\n- list1\n  - item1.1\n  - item1.2\n  - item1.3\n- list2\n  - item2.1",
							Href:        "http://multi.li.ne/2",
							Title:       "Multiline desc",
							Attributes: map[string]string{
								"ADD_DATE": "1456433742",
								"PRIVATE":  "0",
								"TAGS":     "multi",
							},
						},
						{
							Description: "List:\n- item1\n- item2\n\nParagraph number one.\n\nParagraph\nnumber\ntwo.",
							Href:        "http://multi.li.ne/3",
							Title:       "Multiline desc",
							Attributes: map[string]string{
								"ADD_DATE": "1456433747",
								"PRIVATE":  "0",
								"TAGS":     "multi",
							},
						},
					},
				},
			},
		},

		{
			tname:         "Netscape (nested)",
			inputFilename: "netscape_nested.htm",
			want: ast.File{
				Title: "Bookmarks",
				Root: ast.Folder{
					Name: "Bookmarks",
					Bookmarks: []ast.Bookmark{
						{
							Href:  "http://nest.ed/1",
							Title: "Nested 1",
							Attributes: map[string]string{
								"ADD_DATE": "1456433741",
								"PRIVATE":  "0",
								"TAGS":     "tag1,tag2, multi word",
							},
						},
						{
							Href:  "http://nest.ed/2",
							Title: "Nested 2",
							Attributes: map[string]string{
								"ADD_DATE": "1456733741",
								"PRIVATE":  "0",
								"TAGS":     "tag4",
							},
						},
					},
					Subfolders: []ast.Folder{
						{
							Name: "Folder1, the first,folder to encounter",
							Attributes: map[string]string{
								"ADD_DATE":      "1456433722",
								"LAST_MODIFIED": "1456433739",
							},
							Bookmarks: []ast.Bookmark{
								{
									Href:  "http://nest.ed/1-1",
									Title: "Nested 1-1",
									Attributes: map[string]string{
										"ADD_DATE": "1456433742",
										"PRIVATE":  "0",
										"TAGS":     "tag1,tag2,multi word",
									},
								},
								{
									Href:  "http://nest.ed/1-2",
									Title: "Nested 1-2",
									Attributes: map[string]string{
										"ADD_DATE": "1456433747",
										"PRIVATE":  "0",
										"TAGS":     "tag3,tag4, leaf multi word",
									},
								},
							},
						},
						{
							Name:        "Folder2",
							Description: "This second folder contains wonderful links!",
							Attributes: map[string]string{
								"ADD_DATE": "1456433722",
							},
							Bookmarks: []ast.Bookmark{
								{
									Description: "First link of the second section",
									Href:        "http://nest.ed/2-1",
									Title:       "Nested 2-1",
									Attributes: map[string]string{
										"ADD_DATE": "1454433742",
										"PRIVATE":  "0",
									},
								},
								{
									Description: "Second link of the second section",
									Href:        "http://nest.ed/2-2",
									Title:       "Nested 2-2",
									Attributes: map[string]string{
										"ADD_DATE": "1453233747",
										"PRIVATE":  "0",
									},
								},
							},
						},
						{
							Name: "Folder3",
							Subfolders: []ast.Folder{
								{
									Name: "Folder3-1",
									Bookmarks: []ast.Bookmark{
										{
											Href:  "http://nest.ed/3-1",
											Title: "Nested 3-1",
											Attributes: map[string]string{
												"ADD_DATE": "1454433742",
												"PRIVATE":  "0",
												"TAGS":     "tag3",
											},
										},
										{
											Href:  "http://nest.ed/3-2",
											Title: "Nested 3-2",
											Attributes: map[string]string{
												"ADD_DATE": "1453233747",
												"PRIVATE":  "0",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.tname, func(t *testing.T) {
			got, err := ParseFile(filepath.Join("testdata", tc.inputFilename))

			if err != nil {
				t.Errorf("expected no error, got %q", err)
				return
			}

			if got.Title != tc.want.Title {
				t.Errorf("want title %q, got %q", tc.want.Title, got.Title)
			}

			assertFoldersEqual(t, got.Root, tc.want.Root)
		})
	}
}
func assertFoldersEqual(t *testing.T, got ast.Folder, want ast.Folder) {
	t.Helper()

	if got.Name != want.Name {
		t.Errorf("want folder name %q, got %q", want.Name, got.Name)
	}

	if got.Description != want.Description {
		t.Errorf("want folder description %q, got %q", want.Description, got.Description)
	}

	if len(got.Attributes) != len(want.Attributes) {
		t.Errorf("want %d attributes for folder %q, got %d", len(want.Attributes), want.Name, len(got.Attributes))
		return
	}

	for name, wantValue := range want.Attributes {
		if got.Attributes[name] != wantValue {
			t.Errorf("want folder %q attribute %q value %q, got %q", want.Name, name, wantValue, got.Attributes[name])
		}
	}

	if len(got.Bookmarks) != len(want.Bookmarks) {
		t.Errorf("want %d bookmarks in folder %q, got %d", len(want.Bookmarks), want.Name, len(got.Bookmarks))
		return
	}

	for index, wantBookmark := range want.Bookmarks {
		if got.Bookmarks[index].Description != wantBookmark.Description {
			t.Errorf("want bookmark %d description %q, got %q", index, wantBookmark.Description, got.Bookmarks[index].Description)
		}

		if got.Bookmarks[index].Href != wantBookmark.Href {
			t.Errorf("want bookmark %d href %q, got %q", index, wantBookmark.Href, got.Bookmarks[index].Href)
		}

		if got.Bookmarks[index].Title != wantBookmark.Title {
			t.Errorf("want bookmark %d title %q, got %q", index, wantBookmark.Title, got.Bookmarks[index].Title)
		}

		if len(got.Bookmarks[index].Attributes) != len(wantBookmark.Attributes) {
			t.Errorf("want %d attributes for bookmark %d, got %d", len(wantBookmark.Attributes), index, len(got.Bookmarks[index].Attributes))
		}

		for name, wantValue := range wantBookmark.Attributes {
			if got.Bookmarks[index].Attributes[name] != wantValue {
				t.Errorf("want bookmark %d attribute %q value %q, got %q", index, name, wantValue, got.Bookmarks[index].Attributes[name])
			}
		}
	}

	if len(got.Subfolders) != len(want.Subfolders) {
		t.Errorf("want %d subfolders for folder %q, got %d", len(want.Subfolders), want.Name, len(got.Subfolders))
		return
	}

	for index, wantSubfolder := range want.Subfolders {
		assertFoldersEqual(t, got.Subfolders[index], wantSubfolder)
	}
}
