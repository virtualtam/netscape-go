package decoder

import (
	"net/url"
	"sort"
	"strings"

	"github.com/virtualtam/netscape-go"
	"github.com/virtualtam/netscape-go/ast"
)

func decodeBookmark(b ast.Bookmark) (netscape.Bookmark, error) {
	bookmark := netscape.Bookmark{
		Title:       b.Title,
		Description: b.Description,
		Attributes:  map[string]string{},
	}

	url, err := url.Parse(b.Href)
	if err != nil {
		return netscape.Bookmark{}, err
	}
	bookmark.URL = *url

	for attr, value := range b.Attributes {
		switch attr {
		case createdAtAttr:
			createdAt, err := decodeDate(value)
			if err != nil {
				return netscape.Bookmark{}, err
			}
			bookmark.CreatedAt = createdAt
		case updatedAtAttr:
			updatedAt, err := decodeDate(value)
			if err != nil {
				return netscape.Bookmark{}, err
			}
			bookmark.UpdatedAt = updatedAt
		case privateAttr:
			if value == "1" {
				bookmark.Private = true
			}
		case tagsAttr:
			bookmark.Tags = decodeTags(b.Attributes)
		default:
			bookmark.Attributes[attr] = value
		}
	}

	return bookmark, nil
}

func decodeTags(attr map[string]string) []string {
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
