package ohdear

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"git.sr.ht/~jamesponddotco/httpx-go"
	"git.sr.ht/~jamesponddotco/ohdear-go/internal/endpoint"
	"git.sr.ht/~jamesponddotco/ohdear-go/internal/jsonutil"
)

// SitesService handles communication with the /sites endpoint of Oh Dear's API.
type SitesService service

type Sites struct {
	Data []Site `json:"data"`
	Pagination
}

type Site struct {
	CreatedAt                            jsonutil.Time `json:"created_at,omitempty"`
	UpdatedAt                            jsonutil.Time `json:"updated_at,omitempty"`
	LatestRunDate                        jsonutil.Time `json:"latest_run_date,omitempty"`
	GroupName                            *string       `json:"group_name,omitempty"`
	HTTPClientHeaders                    *string       `json:"http_client_headers,omitempty"`
	MarkedForDeletionAt                  *string       `json:"marked_for_deletion_at,omitempty"`
	BrokenLinksWhitelistedURLs           *string       `json:"broken_links_whitelisted_urls,omitempty"`
	Notes                                *string       `json:"notes,omitempty"`
	FriendlyName                         *string       `json:"friendly_name,omitempty"`
	Label                                string        `json:"label,omitempty"`
	SortURL                              string        `json:"sort_url,omitempty"`
	URL                                  string        `json:"url,omitempty"`
	SummarizedCheckResult                string        `json:"summarized_check_result,omitempty"`
	Checks                               []Check       `json:"checks,omitempty"`
	Tags                                 []string      `json:"tags,omitempty"`
	UptimeCheckPayload                   []string      `json:"uptime_check_payload,omitempty"`
	ID                                   int           `json:"id,omitempty"`
	TeamID                               int           `json:"team_id,omitempty"`
	UsesHTTPS                            bool          `json:"uses_https,omitempty"`
	BrokenLinksCheckIncludeExternalLinks bool          `json:"broken_links_check_include_external_links,omitempty"`
}

type Check struct {
	LatestRunEndedAt jsonutil.Time `json:"latest_run_ended_at,omitempty"`
	Type             string        `json:"type,omitempty"`
	Label            string        `json:"label,omitempty"`
	LatestRunResult  string        `json:"latest_run_result,omitempty"`
	Summary          string        `json:"summary,omitempty"`
	ID               int           `json:"id,omitempty"`
	Enabled          bool          `json:"enabled,omitempty"`
}

// List returns a list of all sites in your account.
//
// [API Reference].
//
// [API Reference]: https://ohdear.app/docs/integrations/the-oh-dear-api#get-all-sites-in-your-account
func (s *SitesService) List(ctx context.Context, page uint) (*Sites, *Pagination, *Response, error) {
	if ctx == nil {
		return nil, nil, nil, ErrNilContext
	}

	path := DefaultBaseURL + endpoint.Sites

	if page > 1 {
		path += "?page=[number]=" + strconv.Itoa(int(page))
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, nil, nil, err
	}

	ret, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, nil, nil, err
	}

	var sites Sites
	if err := json.Unmarshal(ret.Body, &sites); err != nil {
		return nil, nil, nil, fmt.Errorf("could not unmarshal sites: %w", err)
	}

	return &sites, &sites.Pagination, ret, nil
}

// Get returns a single site by ID.
//
// [API Reference].
//
// [API Reference]: https://ohdear.app/docs/integrations/the-oh-dear-api#get-a-specific-site-via-the-api
func (s *SitesService) Get(ctx context.Context, id uint) (*Site, *Response, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	if id == 0 {
		return nil, nil, ErrInvalidSiteID
	}

	path := DefaultBaseURL + endpoint.Sites + "/" + strconv.Itoa(int(id))

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, http.NoBody)
	if err != nil {
		return nil, nil, err
	}

	ret, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	var site Site
	if err := json.Unmarshal(ret.Body, &site); err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal site: %w", err)
	}

	return &site, ret, nil
}

// Add adds a new site to your account.
//
// [API Reference].
//
// [API Reference]: https://ohdear.app/docs/integrations/the-oh-dear-api#add-a-site-through-the-api
func (s *SitesService) Add(ctx context.Context, site *Site) (*Site, *Response, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	if site == nil {
		return nil, nil, ErrNilSite
	}

	if site.URL == "" {
		return nil, nil, fmt.Errorf("%w: URL cannot be empty", ErrInvalidURL)
	}

	// Oh Dear requires URLs with HTTP or HTTPS prefixes.
	if !strings.HasPrefix(site.URL, "http://") && !strings.HasPrefix(site.URL, "https://") {
		return nil, nil, fmt.Errorf("%w: URL must start with http:// or https://", ErrInvalidURL)
	}

	if _, err := url.ParseRequestURI(site.URL); err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrInvalidURL, err)
	}

	if site.TeamID == 0 {
		return nil, nil, ErrInvalidTeamID
	}

	payload, err := httpx.WriteJSON(site)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	path := DefaultBaseURL + endpoint.Sites

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, payload)
	if err != nil {
		return nil, nil, err
	}

	ret, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	var addedSite Site
	if err := json.Unmarshal(ret.Body, &addedSite); err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal added site: %w", err)
	}

	return &addedSite, ret, nil
}
