package netscape_test

import (
	"fmt"

	"github.com/virtualtam/netscape-go"
	"github.com/virtualtam/netscape-go/types"
)

func ExampleMarshal() {
	document := types.Document{
		Title: "Bookmarks",
		Root: types.Folder{
			Name: "Bookmarks",
			Bookmarks: []types.Bookmark{
				{
					Href:  "https://domain.tld",
					Title: "Test Domain",
				},
				{
					Description: "Local\nLocal\nLocal",
					Href:        "https://local.domain.tld",
					Title:       "Local Test Domain",
				},
			},
			Subfolders: []types.Folder{
				{
					Name: "Sub",
					Bookmarks: []types.Bookmark{
						{
							Href:  "https://domain.tld",
							Title: "Test Domain",
							Attributes: map[string]string{
								"ATTR1": "v1",
								"ATTR2": "42",
							},
						},
						{
							Description: "Local\nLocal\nLocal",
							Href:        "https://local.domain.tld",
							Title:       "Local Test Domain",
						},
					},
				},
			},
		},
	}

	m, err := netscape.Marshal(&document)
	if err != nil {
		panic(err)
	}

	fmt.Print(string(m))

	// Output:
	// <!DOCTYPE NETSCAPE-Bookmark-file-1>
	// <!-- This is an automatically generated file.
	//      It will be read and overwritten.
	//      DO NOT EDIT! -->
	// <TITLE>Bookmarks</TITLE>
	// <H1>Bookmarks</H1>
	// <DL><p>
	//     <DT><A HREF="https://domain.tld" PRIVATE="0">Test Domain</A>
	//     <DT><A HREF="https://local.domain.tld" PRIVATE="0">Local Test Domain</A>
	//     <DD>Local
	// Local
	// Local
	//     <DT><H3>Sub</H3>
	//     <DL><p>
	//         <DT><A HREF="https://domain.tld" PRIVATE="0" ATTR1="v1" ATTR2="42">Test Domain</A>
	//         <DT><A HREF="https://local.domain.tld" PRIVATE="0">Local Test Domain</A>
	//         <DD>Local
	// Local
	// Local
	//     </DL><p>
	// </DL><p>
}
