package decoder

import (
	"testing"
	"time"
)

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
	}

	for _, tc := range cases {
		t.Run(tc.tname, func(t *testing.T) {
			got, err := decodeDate(tc.input)

			if err != nil {
				t.Errorf("expected no error, got %q", err)
				return
			}

			if got != tc.want {
				t.Errorf("want date/time %q, got %q", tc.want, got)
			}
		})
	}
}
