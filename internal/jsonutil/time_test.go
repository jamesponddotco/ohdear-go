package jsonutil_test

import (
	"testing"
	"time"

	"git.sr.ht/~jamesponddotco/ohdear-go/internal/jsonutil"
)

func TestTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		t             jsonutil.Time
		marshalJSON   string
		unmarshalJSON string
		expectErr     bool
	}{
		{
			name:          "Zero time",
			t:             jsonutil.Time{time.Time{}},
			marshalJSON:   "null",
			unmarshalJSON: "null",
			expectErr:     false,
		},
		{
			name:          "Non-zero time",
			t:             jsonutil.Time{time.Date(2023, 5, 14, 12, 0, 0, 0, time.UTC)},
			marshalJSON:   `"2023-05-14 12:00:00"`,
			unmarshalJSON: `"2023-05-14 12:00:00"`,
			expectErr:     false,
		},
		{
			name:          "Invalid time format",
			t:             jsonutil.Time{time.Time{}},
			marshalJSON:   "null",
			unmarshalJSON: `"invalid time"`,
			expectErr:     true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b, err := tt.t.MarshalJSON()
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)

				return
			}
			if string(b) != tt.marshalJSON {
				t.Errorf("MarshalJSON() = %s, want %s", b, tt.marshalJSON)
			}

			var tm jsonutil.Time

			err = tm.UnmarshalJSON([]byte(tt.unmarshalJSON))
			if (err != nil) != tt.expectErr {
				t.Errorf("UnmarshalJSON() error = %v, expectErr %v", err, tt.expectErr)
			}

			if !tt.expectErr && !tt.t.Equal(tm.Time) {
				t.Errorf("UnmarshalJSON() = %v, want %v", tm, tt.t)
			}
		})
	}
}
