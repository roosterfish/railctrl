package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
)

type Client struct {
	http *http.Client
}

func newClient(address string) *http.Client {
	var dial func(context.Context, string, string) (net.Conn, error)

	if address == "" {
		dial = func(ctx context.Context, _, _ string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", DaemonSocket)
		}
	} else {
		dial = func(ctx context.Context, _, _ string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp", address)
		}
	}

	return &http.Client{
		Transport: &http.Transport{
			DialContext: dial,
		},
	}
}

func (c *Client) request(ctx context.Context, method string, url string, in io.Reader, out any) error {
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		"http://railctl"+url,
		in,
	)
	if err != nil {
		return fmt.Errorf("failed creating HTTP request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed sending HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var apiErr struct {
			Error string
		}

		err := json.NewDecoder(resp.Body).Decode(&apiErr)
		if err != nil {
			return fmt.Errorf("failed decoding error response with code %d", resp.StatusCode)
		}

		return &APIError{
			Code:    resp.StatusCode,
			Message: apiErr.Error,
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return fmt.Errorf("failed decoding response: %w", err)
	}

	return nil
}

func (c *Client) GetRouteFind(ctx context.Context, start, end string) (*Route, error) {
	route := Route{}
	url := "/api/route/find?start=" + url.QueryEscape(start) + "&end=" + url.QueryEscape(end)

	err := c.request(ctx, http.MethodGet, url, nil, &route)
	if err != nil {
		return nil, err
	}

	return &route, nil
}

func NewClient(address string) *Client {
	return &Client{
		http: newClient(address),
	}
}
