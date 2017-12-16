// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/miniflux/miniflux/helper"
	"github.com/miniflux/miniflux/logger"
)

const userAgent = "Miniflux <https://miniflux.net/>"
const requestTimeout = 300

// Client is a HTTP Client :)
type Client struct {
	url                string
	etagHeader         string
	lastModifiedHeader string
	username           string
	password           string
	Insecure           bool
}

// Get execute a GET HTTP request.
func (c *Client) Get() (*Response, error) {
	defer helper.ExecutionTime(time.Now(), fmt.Sprintf("[HttpClient:Get] url=%s", c.url))

	client := c.buildClient()
	resp, err := client.Do(c.buildRequest())
	if err != nil {
		return nil, err
	}

	response := &Response{
		Body:         resp.Body,
		StatusCode:   resp.StatusCode,
		EffectiveURL: resp.Request.URL.String(),
		LastModified: resp.Header.Get("Last-Modified"),
		ETag:         resp.Header.Get("ETag"),
		ContentType:  resp.Header.Get("Content-Type"),
	}

	logger.Debug("[HttpClient:Get]",
		"OriginalURL:", c.url,
		"StatusCode:", response.StatusCode,
		"ETag:", response.ETag,
		"LastModified:", response.LastModified,
		"EffectiveURL:", response.EffectiveURL,
	)

	return response, err
}

func (c *Client) buildRequest() *http.Request {
	link, _ := url.Parse(c.url)
	request := &http.Request{
		URL:    link,
		Method: http.MethodGet,
		Header: c.buildHeaders(),
	}

	if c.username != "" && c.password != "" {
		request.SetBasicAuth(c.username, c.password)
	}

	return request
}

func (c *Client) buildClient() http.Client {
	client := http.Client{Timeout: time.Duration(requestTimeout * time.Second)}
	if c.Insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return client
}

func (c *Client) buildHeaders() http.Header {
	headers := make(http.Header)
	headers.Add("User-Agent", userAgent)

	if c.etagHeader != "" {
		headers.Add("If-None-Match", c.etagHeader)
	}

	if c.lastModifiedHeader != "" {
		headers.Add("If-Modified-Since", c.lastModifiedHeader)
	}

	return headers
}

// NewClient returns a new HTTP client.
func NewClient(url string) *Client {
	return &Client{url: url, Insecure: false}
}

// NewClientWithCredentials returns a new HTTP client that require authentication.
func NewClientWithCredentials(url, username, password string) *Client {
	return &Client{url: url, Insecure: false, username: username, password: password}
}

// NewClientWithCacheHeaders returns a new HTTP client that send cache headers.
func NewClientWithCacheHeaders(url, etagHeader, lastModifiedHeader string) *Client {
	return &Client{url: url, etagHeader: etagHeader, lastModifiedHeader: lastModifiedHeader, Insecure: false}
}
