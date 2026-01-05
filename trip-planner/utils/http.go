package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// More reliable HTTP client for external APIs (OpenAI/Google/Weather)
var DefaultClient = &http.Client{
	Timeout: 120 * time.Second, // was 25s
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   20 * time.Second,
		ResponseHeaderTimeout: 120 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
	},
}

func GetJSON(ctx context.Context, url string, out any, headers map[string]string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GET %s failed: %d - %s", url, resp.StatusCode, string(b))
	}

	dec := json.NewDecoder(resp.Body)
	dec.DisallowUnknownFields() // optional: remove if you expect unknown fields
	return dec.Decode(out)
}

func PostJSON(ctx context.Context, url string, payload any, out any, headers map[string]string) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("POST %s failed: %d - %s", url, resp.StatusCode, string(body))
	}

	dec := json.NewDecoder(resp.Body)
	return dec.Decode(out)
}
