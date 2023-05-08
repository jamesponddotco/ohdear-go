package ohdear

import "git.sr.ht/~jamesponddotco/xstd-go/xerrors"

const (
	// ErrNilContext is return when a nil context is passed to a function.
	ErrNilContext xerrors.Error = "context cannot be nil"

	// ErrInvalidSiteID is returned when the site ID passed to a function is zero.
	ErrInvalidSiteID xerrors.Error = "site ID cannot be zero"
)
