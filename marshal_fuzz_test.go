package netscape_test

import (
	"testing"

	"github.com/virtualtam/netscape-go/v2"
)

func Fuzz(f *testing.F) {
	f.Fuzz(func(t *testing.T, input []byte) {
		document, err := netscape.Unmarshal(input)
		if err != nil {
			return
		}

		_, err = netscape.Marshal(document)
		if err != nil {
			t.Error(err)
		}
	})
}
