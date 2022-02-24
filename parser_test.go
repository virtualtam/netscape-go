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
			input: `<!DOCTYPE NETSCAPE-Bookmark-file-1>\n`,
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
	}

	for _, tc := range cases {
		t.Run(tc.tname, func(t *testing.T) {
			r := strings.NewReader(tc.input)

			_, err := Parse(r)

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
		})
	}
}
