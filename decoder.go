// Copyright (c) VirtualTam
// SPDX-License-Identifier: MIT

package netscape

import (
	"html"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	createdAtAttr string = "ADD_DATE"
	updatedAtAttr string = "LAST_MODIFIED"
	privateAttr   string = "PRIVATE"
	tagsAttr      string = "TAGS"
)

// Decode walks a Netscape Bookmark AST and returns the corresponding document.
func Decode(f FileNode) (*Document, error) {
	d := NewDecoder()
	return d.decodeFile(f)
}

// A Decoder walks a Netscape Bookmark AST and returns the corresponding
// document.
type Decoder struct {
	now     time.Time
	maxTime time.Time
}

// NewDecoder initializes and returns a new Decoder.
func NewDecoder() *Decoder {
	now := time.Now().UTC()
	rangeYears := 30
	maxTime := now.AddDate(rangeYears, 0, 0)

	return &Decoder{
		now:     now,
		maxTime: maxTime,
	}
}

func (d *Decoder) decodeFile(f FileNode) (*Document, error) {
	document := Document{
		Title: f.Title,
	}

	root, err := d.decodeFolder(f.Root)
	if err != nil {
		return &Document{}, err
	}
	document.Root = root

	return &document, nil
}

func (d *Decoder) decodeFolder(f FolderNode) (Folder, error) {
	folder := Folder{
		Name:        f.Name,
		Description: f.Description,
		Attributes:  map[string]string{},
	}

	for attr, value := range f.Attributes {
		switch attr {
		case createdAtAttr:
			createdAt, err := d.decodeDate(value)
			if err != nil {
				return Folder{}, err
			}
			folder.CreatedAt = createdAt
			if folder.UpdatedAt.IsZero() {
				folder.UpdatedAt = createdAt
			}
		case updatedAtAttr:
			updatedAt, err := d.decodeDate(value)
			if err != nil {
				return Folder{}, err
			}
			folder.UpdatedAt = updatedAt
		default:
			folder.Attributes[attr] = value
		}
	}

	for _, b := range f.Bookmarks {
		bookmark, err := d.decodeBookmark(b)
		if err != nil {
			return Folder{}, err
		}

		folder.Bookmarks = append(folder.Bookmarks, bookmark)
	}

	for _, sf := range f.Subfolders {
		subfolder, err := d.decodeFolder(sf)
		if err != nil {
			return Folder{}, err
		}

		folder.Subfolders = append(folder.Subfolders, subfolder)
	}

	return folder, nil
}

func (d *Decoder) decodeBookmark(b BookmarkNode) (Bookmark, error) {
	bookmark := Bookmark{
		Description: html.UnescapeString(b.Description),
		URL:         b.Href,
		Title:       b.Title,
		Attributes:  map[string]string{},
	}

	for attr, value := range b.Attributes {
		switch attr {
		case createdAtAttr:
			createdAt, err := d.decodeDate(value)
			if err != nil {
				return Bookmark{}, err
			}
			bookmark.CreatedAt = createdAt
			if bookmark.UpdatedAt.IsZero() {
				bookmark.UpdatedAt = createdAt
			}
		case updatedAtAttr:
			updatedAt, err := d.decodeDate(value)
			if err != nil {
				return Bookmark{}, err
			}
			bookmark.UpdatedAt = updatedAt
		case privateAttr:
			if value == "1" {
				bookmark.Private = true
			}
		case tagsAttr:
			bookmark.Tags = d.decodeTags(b.Attributes)
		default:
			bookmark.Attributes[attr] = value
		}
	}

	return bookmark, nil
}

func (d *Decoder) decodeTags(attr map[string]string) []string {
	rawTags, ok := attr[tagsAttr]
	if !ok {
		return []string{}
	}

	tags := strings.Split(rawTags, ",")
	for index, tag := range tags {
		tags[index] = strings.TrimSpace(tag)
	}

	sort.Strings(tags)

	return tags
}

func (d *Decoder) decodeDate(input string) (time.Time, error) {
	// First, attempt to parse the date as a UNIX epoch, which is the most
	// commonly used format.
	unixTime, err := strconv.ParseInt(input, 10, 64)
	if err == nil {
		return d.decodeUnixDate(unixTime), nil
	}

	// Attempt to parse the date as RFC3339
	date, err := d.decodeRFC3339Date(input)
	if err == nil {
		return date, nil
	}

	// Attempt to parse the date as Comon Log
	date, err = d.decodeCommonLogDate(input)
	if err == nil {
		return date, nil
	}

	return time.Time{}, err
}

const (
	commonLogLayout string = "02/Jan/2006:15:04:05 -0700"
)

// decodeCommonLogDate returns the time.Time corresponding to a date formatted
// in the Common Log format.
func (d *Decoder) decodeCommonLogDate(input string) (time.Time, error) {
	date, err := time.Parse(commonLogLayout, input)
	if err == nil {
		return date.UTC(), nil
	}

	return time.Time{}, err
}

// decodeRFC3339Date returns the time.Time corresponding to a RFC3339
// representation.
func (d *Decoder) decodeRFC3339Date(input string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, input)
	if err == nil {
		return date.UTC(), nil
	}

	date, err = time.Parse(time.RFC3339Nano, input)
	if err == nil {
		return date.UTC(), nil
	}

	return time.Time{}, err
}

// decodeUnixDate returns the time.Time corresponding to a UNIX timestamp.
//
// Dates are usually specified in seconds, but some browsers and bookmarking
// services may use milliseconds, microseconds or nanoseconds.
//
// To address these cases, we ensure the resulting time.Time is comprised in a
// reasonable interval (ie not further than N years in the future).
func (d *Decoder) decodeUnixDate(unixTime int64) time.Time {
	date := time.Unix(unixTime, 0).UTC()

	if date.After(d.maxTime) {
		date = time.UnixMilli(unixTime).UTC()
	}

	if date.After(d.maxTime) {
		date = time.UnixMicro(unixTime).UTC()
	}

	if date.After(d.maxTime) {
		date = time.Unix(0, unixTime).UTC()
	}

	return date
}
