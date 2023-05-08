package jsonutil

import (
	"fmt"
	"strings"
	"time"
)

// Time is a wrapper around time.Time that allows us to marshal/unmarshal Time
// the format used by the Oh Dear API.
type Time struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) { //nolint:unparam // required by the interface
	if t.Time.IsZero() {
		return []byte("null"), nil
	}

	return []byte(t.Time.Format(`"2006-01-02 15:04:05"`)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.Time = time.Time{}

		return nil
	}

	s := strings.Trim(string(data), "\"")

	tt, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	t.Time = tt

	return nil
}
