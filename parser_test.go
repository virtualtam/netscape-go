package netscape

import (
	"errors"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		tname   string
		input   string
		want    File
		wantErr error
	}{
		// nominal cases
		{
			tname: "valid document",
			input: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
     It will be read and overwritten.
     DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
`,
			want: File{
				Title: "Bookmarks",
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
			}

			if got.Title != tc.want.Title {
				t.Errorf("want title %q, got %q", tc.want.Title, got.Title)
			}
		})
	}
}
