// Package urlutil provides utility functions for working with URLs.
package urlutil

import (
	"fmt"
	"net/url"
	"strings"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
)

// ErrInvalidURL is returned when the URL passed to a function is empty or cannot be parsed.
const ErrInvalidURL xerrors.Error = "invalid URL"

// Validate checks if the given URL is valid for use with the Oh Dear API.
func Validate(uri string) error {
	if uri == "" {
		return fmt.Errorf("%w: URL cannot be empty", ErrInvalidURL)
	}

	// Oh Dear requires URLs with HTTP or HTTPS prefixes.
	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		return fmt.Errorf("%w: URL must start with http:// or https://", ErrInvalidURL)
	}

	if _, err := url.ParseRequestURI(uri); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidURL, err)
	}

	return nil
}
