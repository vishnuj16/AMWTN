// Package shared provides a minimal Upstash Redis REST client and small
// helpers shared across every Vercel Go serverless function in /api.
// It intentionally avoids external dependencies (net/http + encoding/json
// from the standard library only) to keep cold starts fast and the
// deployment footprint tiny.
package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

var httpClient = &http.Client{Timeout: 8 * time.Second}

// redisResult mirrors the shape Upstash's REST API returns:
// {"result": <value>} on success or {"error": "..."} on failure.
type redisResult struct {
	Result json.RawMessage `json:"result"`
	Error  string          `json:"error"`
}

// cmd sends a single Redis command to Upstash using its REST pipeline
// endpoint (POST body is a JSON array: ["SET", "key", "value"]).
func cmd(args ...string) (json.RawMessage, error) {
	url := os.Getenv("UPSTASH_REDIS_REST_URL")
	token := os.Getenv("UPSTASH_REDIS_REST_TOKEN")
	if url == "" || token == "" {
		return nil, errors.New("UPSTASH_REDIS_REST_URL / UPSTASH_REDIS_REST_TOKEN are not set")
	}

	body, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out redisResult
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if out.Error != "" {
		return nil, fmt.Errorf("redis error: %s", out.Error)
	}
	return out.Result, nil
}

// Get returns the raw string value for key, and false if the key is missing.
func Get(key string) (string, bool, error) {
	res, err := cmd("GET", key)
	if err != nil {
		return "", false, err
	}
	if res == nil || string(res) == "null" {
		return "", false, nil
	}
	var s string
	if err := json.Unmarshal(res, &s); err != nil {
		return "", false, err
	}
	return s, true, nil
}

// Set stores a string value for key with no expiry.
func Set(key, value string) error {
	_, err := cmd("SET", key, value)
	return err
}

// SetEx stores a string value for key that expires after ttlSeconds.
func SetEx(key, value string, ttlSeconds int) error {
	_, err := cmd("SET", key, value, "EX", fmt.Sprintf("%d", ttlSeconds))
	return err
}

// Del removes a key.
func Del(key string) error {
	_, err := cmd("DEL", key)
	return err
}

// Exists reports whether a key is present.
func Exists(key string) (bool, error) {
	res, err := cmd("EXISTS", key)
	if err != nil {
		return false, err
	}
	var n int
	if err := json.Unmarshal(res, &n); err != nil {
		return false, err
	}
	return n == 1, nil
}

// Incr atomically increments key and returns the new value.
func Incr(key string) (int64, error) {
	res, err := cmd("INCR", key)
	if err != nil {
		return 0, err
	}
	var n int64
	if err := json.Unmarshal(res, &n); err != nil {
		return 0, err
	}
	return n, nil
}
