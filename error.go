package gopi

import "fmt"

///////////////////////////////////////////////////////////////////////////////
// Types

type Error uint

///////////////////////////////////////////////////////////////////////////////
// Globals

const (
	ErrNone Error = iota
	ErrBadParameter
	ErrNotImplemented
	ErrUnexpectedResponse
	ErrHelp
	ErrInternalAppError
	ErrDuplicateEntry
	ErrOutOfOrder
)

///////////////////////////////////////////////////////////////////////////////
// Implementation

func (e Error) Error() string {
	switch e {
	case ErrNone:
		return "No Error"
	case ErrBadParameter:
		return "Bad Parameter"
	case ErrNotImplemented:
		return "Not Implemented"
	case ErrHelp:
		return "Help Requested"
	case ErrUnexpectedResponse:
		return "Unexpected Response"
	case ErrInternalAppError:
		return "Internal Application Error"
	case ErrDuplicateEntry:
		return "Duplicate Entry"
	case ErrOutOfOrder:
		return "Out of Order"
	default:
		return "[?? Invalid Error]"
	}
}

func (e Error) WithPrefix(prefix interface{}) error {
	return fmt.Errorf("%v: %w", prefix, e)
}
