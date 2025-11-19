// Copyright (c) VirtualTam
// SPDX-License-Identifier: MIT

package netscape_test

import (
	"os"
	"testing"

	"github.com/virtualtam/netscape-go/v2"
)

var (
	cases = []struct {
		name     string
		filepath string
	}{
		{
			name:     "flat: 100 bookmarks",
			filepath: "testdata/benchmark/flat100.htm",
		},
		{
			name:     "flat: 1000 bookmarks",
			filepath: "testdata/benchmark/flat1000.htm",
		},
		{
			name:     "flat: 10000 bookmarks",
			filepath: "testdata/benchmark/flat10000.htm",
		},
	}
)

func BenchmarkMarshal(b *testing.B) {
	b.ReportAllocs()

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			bytes, err := os.ReadFile(tc.filepath)
			if err != nil {
				b.Fatalf("failed to open file %q: %s", tc.filepath, err)
			}

			document, err := netscape.Unmarshal(bytes)
			if err != nil {
				b.Fatalf("failed to open file %q: %s", tc.filepath, err)
			}

			b.ResetTimer()

			for b.Loop() {
				_, err := netscape.Marshal(document)
				if err != nil {
					b.Fatalf("failed to marshal document: %s", err)
				}
			}
		})
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	b.ReportAllocs()

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			bytes, err := os.ReadFile(tc.filepath)
			if err != nil {
				b.Fatalf("failed to open file %q: %s", tc.filepath, err)
			}

			b.ResetTimer()

			for b.Loop() {
				_, err := netscape.Unmarshal(bytes)
				if err != nil {
					b.Fatalf("failed to open file %q: %s", tc.filepath, err)
				}
			}
		})
	}
}
