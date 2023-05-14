package urlutil_test

import (
	"testing"

	"git.sr.ht/~jamesponddotco/ohdear-go/internal/urlutil"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		url       string
		wantError bool
	}{
		{
			name:      "Empty URL",
			url:       "",
			wantError: true,
		},
		{
			name:      "URL without HTTP or HTTPS prefix",
			url:       "www.example.com",
			wantError: true,
		},
		{
			name:      "URL with invalid format",
			url:       "http://192.168.0.%31/",
			wantError: true,
		},
		{
			name:      "Valid HTTP URL",
			url:       "http://www.example.com",
			wantError: false,
		},
		{
			name:      "Valid HTTPS URL",
			url:       "https://www.example.com",
			wantError: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := urlutil.Validate(tt.url)

			if (err != nil) != tt.wantError {
				t.Errorf("Validate(%q) error = %v, wantError %v", tt.url, err, tt.wantError)
			}
		})
	}
}
