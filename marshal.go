// Copyright (c) VirtualTam
// SPDX-License-Identifier: MIT

// Package netscape provides utilities to parse and export Web bookmarks using
// the Netscape Bookmark format.
package netscape

import (
	"bytes"
	"io"
	"os"
	"strings"
)

// Marshal returns the Netscape Bookmark encoding of d.
func Marshal(d *Document) ([]byte, error) {
	var buf bytes.Buffer

	if err := NewEncoder(&buf).Encode(d); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

// Unmarshal unmarshals a []byte representation of a Netscape Bookmark
// file and returns the corresponding Document.
func Unmarshal(b []byte) (*Document, error) {
	r := bytes.NewReader(b)
	return unmarshal(r)
}

// UnmarshalFile unmarshals a Netscape Bookmark file and returns the
// corresponding Document.
func UnmarshalFile(filePath string) (*Document, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return &Document{}, err
	}
	defer file.Close()

	return unmarshal(file)
}

// UnmarshalString unmarshals a string representation of a Netscape Bookmark
// file and returns the corresponding Document.
func UnmarshalString(data string) (*Document, error) {
	r := strings.NewReader(data)
	return unmarshal(r)
}

func unmarshal(r io.ReadSeeker) (*Document, error) {
	astFile, err := Parse(r)
	if err != nil {
		return &Document{}, err
	}

	return Decode(*astFile)
}
