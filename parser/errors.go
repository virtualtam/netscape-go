package parser

import "errors"

var (
	ErrDoctypeMissing  error = errors.New("missing DOCTYPE")
	ErrDoctypeInvalid  error = errors.New("invalid DOCTYPE")
	ErrEOFUnexpected   error = errors.New("unexpected end-of-file")
	ErrTokenUnexpected error = errors.New("unexpected token")
)
