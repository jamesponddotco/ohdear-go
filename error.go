package ohdear

import "git.sr.ht/~jamesponddotco/xstd-go/xerrors"

const (
	// ErrNilContext is return when a nil context is passed to a function.
	ErrNilContext xerrors.Error = "context cannot be nil"

	// ErrNilSite is returned when a nil site is passed to a function.
	ErrNilSite xerrors.Error = "site cannot be nil"

	// ErrInvalidSiteID is returned when the site ID passed to a function is zero.
	ErrInvalidSiteID xerrors.Error = "site ID cannot be zero"

	// ErrInvalidURL is returned when the URL passed to a function is empty or cannot be parsed.
	ErrInvalidURL xerrors.Error = "invalid URL"

	// ErrInvalidTeamID is returned when the team ID passed to a function is zero.
	ErrInvalidTeamID xerrors.Error = "team ID cannot be zero"
)
