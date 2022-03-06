package netscape_test

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/virtualtam/netscape-go"
)

func ExampleUnmarshal() {
	blob := `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><H3>Linux Distributions</H3>
	<DL><p>
		<DT><A HREF="https://archlinux.org">Arch Linux</A>
	    <DT><A HREF="https://debian.org">Debian</A>
	</DL><p>
    <DT><H3>Programming Languages</H3>
	<DL><p>
		<DT><A HREF="https://go.dev">Go</A>
		<DT><A HREF="https://www.rust-lang.org/">Rust</A>
	</DL><p>
    <DT><H3>Secret stuff</H3>
	<DL><p>
		<DT><A HREF="https://https://en.wikipedia.org/wiki/Caesar_cipher" PRIVATE="1">Caesar cipher</A>
		<DT><A HREF="https://en.wikipedia.org/wiki/Vigen%C3%A8re_cipher" PRIVATE="1">Vigenère cipher</A>
	</DL><p>
</DL><p>
`

	document, err := netscape.Unmarshal([]byte(blob))
	if err != nil {
		fmt.Println("failed to unmarshal file:", err)
		os.Exit(1)
	}

	jsonData, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		fmt.Println("failed to marshal data as JSON:", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))

	// Output:
	// {
	//   "title": "Bookmarks",
	//   "root": {
	//     "name": "Bookmarks",
	//     "subfolders": [
	//       {
	//         "name": "Linux Distributions",
	//         "bookmarks": [
	//           {
	//             "title": "Arch Linux",
	//             "url": "https://archlinux.org",
	//             "private": false
	//           },
	//           {
	//             "title": "Debian",
	//             "url": "https://debian.org",
	//             "private": false
	//           }
	//         ]
	//       },
	//       {
	//         "name": "Programming Languages",
	//         "bookmarks": [
	//           {
	//             "title": "Go",
	//             "url": "https://go.dev",
	//             "private": false
	//           },
	//           {
	//             "title": "Rust",
	//             "url": "https://www.rust-lang.org/",
	//             "private": false
	//           }
	//         ]
	//       },
	//       {
	//         "name": "Secret stuff",
	//         "bookmarks": [
	//           {
	//             "title": "Caesar cipher",
	//             "url": "https://https://en.wikipedia.org/wiki/Caesar_cipher",
	//             "private": true
	//           },
	//           {
	//             "title": "Vigenère cipher",
	//             "url": "https://en.wikipedia.org/wiki/Vigen%C3%A8re_cipher",
	//             "private": true
	//           }
	//         ]
	//       }
	//     ]
	//   }
	// }
}
