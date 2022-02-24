package netscape

import "errors"

var (
	ErrDoctypeMissing error = errors.New("missing DOCTYPE")
	ErrDoctypeInvalid error = errors.New("invalid DOCTYPE")
)
