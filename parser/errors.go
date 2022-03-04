package parser

import (
	"errors"
	"fmt"
)

var (
	ErrDoctypeMissing error = errors.New("missing DOCTYPE")
	ErrDoctypeInvalid error = errors.New("invalid DOCTYPE")
)

// An Error is returned when we fail to parse a Netscape Bookmark token or XML
// element.
type Error struct {
	Msg string
	Pos int64
	Err error
}

// Error returns the string representation for this error.
func (e *Error) Error() string {
	return fmt.Sprintf("%s at position %d: %s", e.Msg, e.Pos, e.Err)
}

// Is compares this Error with a target error to satisfy an equality check.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}

	return e.Msg == t.Msg && e.Pos == t.Pos
}

// Unwrap returns the inner error wrapped by this Error.
func (e *Error) Unwrap() error {
	return e.Err
}

func wrapWithError(msg string, pos int64, inner error) error {
	return &Error{
		Msg: msg,
		Pos: pos,
		Err: inner,
	}
}
