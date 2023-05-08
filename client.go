package ohdear

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"git.sr.ht/~jamesponddotco/httpx-go"
	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
)

// ErrConfigRequired is returned when a Client is created without a Config.
const ErrConfigRequired xerrors.Error = "config cannot be empty"

// DefaultBaseURL is the default endpoint for the Oh Dear API.
const DefaultBaseURL string = "https://ohdear.app/api"

type (
	// Service is a common struct that can be reused instead of allocating a new
	// one for each service on the heap.
	service struct {
		client *Client
	}

	// Client is a client for the Help Scout Docs API.
	Client struct {
		// httpc is the underlying HTTP client used by the API client.
		httpc *httpx.Client

		// cfg specifies the configuration used by the API client.
		cfg *Config

		// Service fields.
		Sites *SitesService

		// common service fields shared by all services.
		common service
	}
)

// NewClient returns a new client for the Oh Dear API.
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, ErrConfigRequired
	}

	cfg.init()

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	c := &Client{
		httpc: httpx.NewClient(),
		cfg:   cfg,
	}

	c.httpc.UserAgent = cfg.Application.UserAgent()
	c.httpc.Logger = cfg.Logger
	c.httpc.Debug = cfg.Debug

	c.common.client = c
	c.Sites = (*SitesService)(&c.common)

	return c, nil
}

// Do performs an HTTP request using the underlying HTTP client.
func (c *Client) Do(ctx context.Context, req *http.Request) (*Response, error) {
	ret, err := c.httpc.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	defer func() {
		if err = httpx.DrainResponseBody(ret); err != nil {
			log.Fatal(err)
		}
	}()

	if c.cfg.Debug {
		var dump []byte

		dump, err = httputil.DumpResponse(ret, true)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		c.cfg.Logger.Printf("\n%s", dump)
	}

	body, err := io.ReadAll(ret.Body)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	response := &Response{
		Header: ret.Header.Clone(),
		Body:   body,
		Status: ret.StatusCode,
	}

	return response, nil
}

// NewRequest is a convenience function for creating an HTTP request.
func (c *Client) NewRequest(
	ctx context.Context,
	method, uri string,
	body io.Reader,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	ua := c.cfg.Application.UserAgent()
	if ua != nil {
		req.Header.Set("User-Agent", ua.String())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.Key)

	if c.cfg.Debug {
		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		c.cfg.Logger.Printf("\n%s", string(dump))
	}

	return req, nil
}

// Links contains the URLs for pagination navigation.
type Links struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
}

// Meta contains pagination metadata for the API response.
type Meta struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	Pages       int `json:"total"`
}

// Pagination represents a paginated response from the Oh Dear API.
type Pagination struct {
	Links Links `json:"links"`
	Meta  Meta  `json:"meta,omitempty"`
}

// HasNextPage checks if there is a next page in the paginated response.
func (p *Pagination) HasNextPage() bool {
	return p.Links.Next != ""
}

// HasPrevPage checks if there is a previous page in the paginated response.
func (p *Pagination) HasPrevPage() bool {
	return p.Links.Prev != ""
}

// NextPage returns the URL for the next page in the paginated response.
func (p *Pagination) NextPage() (link string, number uint) {
	if p.HasNextPage() {
		return p.Links.Next, uint(p.Meta.CurrentPage + 1)
	}

	return "", 0
}

// PrevPage returns the URL for the previous page in the paginated response.
func (p *Pagination) PrevPage() (link string, number uint) {
	if p.HasPrevPage() {
		return p.Links.Prev, uint(p.Meta.CurrentPage - 1)
	}

	return "", 0
}

// Response represents a response from the Oh Dear API.
type Response struct {
	// Header contains the response headers.
	Header http.Header

	// Body contains the response body as a byte slice.
	Body []byte

	// Status is the HTTP status code of the response.
	Status int
}

// IsSuccessful checks if the response status code is within the successful range.
func (r *Response) IsSuccessful() bool {
	return r.Status >= http.StatusOK && r.Status < http.StatusMultipleChoices || r.Status == http.StatusNoContent
}
