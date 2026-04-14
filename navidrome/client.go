// Package navidrome provides a client for the Subsonic API used by Navidrome.
// Navidrome is a music streaming server that implements the Subsonic API,
// allowing clients to browse, search, and stream music from a personal collection.
//
// This package provides functionality to:
//   - Authenticate using the Subsonic salt+token method (MD5)
//   - Fetch playlists by name
//   - Get song IDs from playlists
//   - Stream audio in original format (FLAC, etc.)
//
// For more information about the Subsonic API, see: http://www.subsonic.org/pages/api.jsp
package navidrome

import (
	crand "crypto/rand"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"ticktask/config"
)

// Subsonic API constants.
const subsonicClient = "ticktask"   // Client identifier sent with API requests
const subsonicVersion = "1.16.0"    // Subsonic API version

// RandomTrackFromPlaylist fetches a random track from the specified playlist.
// Steps:
//  1. Resolves playlist name to ID
//  2. Gets all song IDs in the playlist
//  3. Picks a random song
//  4. Streams and returns the audio data in original format (FLAC)
//
// Returns the raw audio bytes or an error if any step fails.
func RandomTrackFromPlaylist(cfg *config.NavidromeMusic, playlistName string) ([]byte, error) {
	c := &Client{cfg: cfg}
	pid, err := c.playlistIDByName(playlistName)
	if err != nil {
		return nil, err
	}
	ids, err := c.songIDsInPlaylist(pid)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("navidrome: playlist %q has no songs", playlistName)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := ids[rng.Intn(len(ids))]
	return c.streamTrack(id)
}

// Client performs Subsonic REST API calls against a Navidrome server.
// Handles authentication and response parsing for playlist and streaming operations.
type Client struct {
	cfg *config.NavidromeMusic
}

// base returns the normalized base URL for API requests.
// Trims whitespace and trailing slashes.
func (c *Client) base() (string, error) {
	b := strings.TrimSpace(c.cfg.BaseURL)
	b = strings.TrimRight(b, "/")
	if b == "" {
		return "", fmt.Errorf("navidrome: base_url is empty")
	}
	return b, nil
}

// authValues generates the authentication query parameters for Subsonic API.
// Uses the salt+token authentication method:
//   - salt: Random 16-character hex string
//   - token: MD5(password + salt)
//
// This method avoids sending the password in plaintext over the network.
func (c *Client) authValues() (url.Values, error) {
	u := strings.TrimSpace(c.cfg.Username)
	if u == "" {
		return nil, fmt.Errorf("navidrome: username is empty")
	}
	p := c.cfg.Password
	if p == "" {
		return nil, fmt.Errorf("navidrome: password is empty (set in config or %s)", config.NavidromePasswordEnv)
	}
	salt := randomSalt()
	token := md5hex(p + salt)
	v := url.Values{}
	v.Set("u", u)        // username
	v.Set("s", salt)     // salt
	v.Set("t", token)    // token = md5(password + salt)
	v.Set("v", subsonicVersion)
	v.Set("c", subsonicClient)
	v.Set("f", "json")   // response format
	return v, nil
}

// md5hex computes the MD5 hash of a string and returns it as a hex string.
func md5hex(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

// randomSalt generates a random 16-character hex string for authentication.
// Falls back to timestamp-based salt if crypto/rand fails.
func randomSalt() string {
	b := make([]byte, 8)
	if _, err := crand.Read(b); err != nil {
		return fmt.Sprintf("%016x", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

// getJSON makes an authenticated GET request to a Subsonic API endpoint.
// Merges extra query parameters with authentication parameters.
// Returns the parsed payload data (excluding metadata fields).
func (c *Client) getJSON(path string, extra url.Values) (json.RawMessage, error) {
	base, err := c.base()
	if err != nil {
		return nil, err
	}
	auth, err := c.authValues()
	if err != nil {
		return nil, err
	}
	// Merge extra parameters
	for k, vs := range extra {
		for _, v := range vs {
			auth.Add(k, v)
		}
	}
	reqURL := base + "/rest/" + path + "?" + auth.Encode()
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("navidrome: %s HTTP %d: %s", path, resp.StatusCode, truncate(body, 200))
	}
	return parseSubsonicPayload(body)
}

// parseSubsonicPayload extracts the data payload from a Subsonic API response.
// Subsonic responses are wrapped in {"subsonic-response": {...}} with status fields.
// This function validates the status and returns only the actual data.
func parseSubsonicPayload(body []byte) (json.RawMessage, error) {
	var outer struct {
		SR json.RawMessage `json:"subsonic-response"`
	}
	if err := json.Unmarshal(body, &outer); err != nil {
		return nil, err
	}
	var srMap map[string]json.RawMessage
	if err := json.Unmarshal(outer.SR, &srMap); err != nil {
		return nil, err
	}
	// Validate response status
	rawStatus, ok := srMap["status"]
	if !ok {
		return nil, fmt.Errorf("navidrome: missing status in response")
	}
	var status string
	if err := json.Unmarshal(rawStatus, &status); err != nil {
		return nil, err
	}
	if status != "ok" {
		if e, ok := srMap["error"]; ok {
			var errObj struct {
				Message string `json:"message"`
			}
			_ = json.Unmarshal(e, &errObj)
			if errObj.Message != "" {
				return nil, fmt.Errorf("navidrome: %s", errObj.Message)
			}
		}
		return nil, fmt.Errorf("navidrome: status=%s", status)
	}
	// Remove metadata fields, return data payload
	for _, k := range []string{"status", "version", "type", "serverVersion", "openSubsonic", "error"} {
		delete(srMap, k)
	}
	// Return known payload types
	if v, ok := srMap["playlists"]; ok {
		return v, nil
	}
	if v, ok := srMap["playlist"]; ok {
		return v, nil
	}
	// Return single remaining key if present
	if len(srMap) == 1 {
		for _, v := range srMap {
			return v, nil
		}
	}
	return nil, fmt.Errorf("navidrome: unexpected API response shape (%d keys)", len(srMap))
}

// truncate limits a string to n characters, appending "..." if truncated.
func truncate(b []byte, n int) string {
	s := string(b)
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// httpClient returns an HTTP client with a 60-second timeout.
func httpClient() *http.Client {
	return &http.Client{Timeout: 60 * time.Second}
}

// playlistSummary represents basic playlist info from getPlaylists API.
type playlistSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// playlistIDByName finds a playlist by name and returns its ID.
// Uses the getPlaylists API endpoint.
func (c *Client) playlistIDByName(want string) (string, error) {
	raw, err := c.getJSON("getPlaylists", nil)
	if err != nil {
		return "", err
	}
	var outer struct {
		Playlist json.RawMessage `json:"playlist"`
	}
	if err := json.Unmarshal(raw, &outer); err != nil {
		return "", err
	}
	lists, err := parseJSONArrayOrOne[playlistSummary](outer.Playlist)
	if err != nil {
		return "", err
	}
	for _, p := range lists {
		if p.Name == want {
			return p.ID, nil
		}
	}
	return "", fmt.Errorf("navidrome: no playlist named %q (create it in Navidrome or run 'ticktask music config')", want)
}

// songIDsInPlaylist returns all song IDs in a playlist.
// Filters out directory entries (isDir=true).
func (c *Client) songIDsInPlaylist(playlistID string) ([]string, error) {
	raw, err := c.getJSON("getPlaylist", url.Values{"id": {playlistID}})
	if err != nil {
		return nil, err
	}
	var pl struct {
		Entry json.RawMessage `json:"entry"`
	}
	if err := json.Unmarshal(raw, &pl); err != nil {
		return nil, err
	}
	entries, err := parseJSONArrayOrOne[songEntry](pl.Entry)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir {
			continue
		}
		if e.ID != "" {
			ids = append(ids, e.ID)
		}
	}
	return ids, nil
}

// songEntry represents a track in a playlist.
type songEntry struct {
	ID    string `json:"id"`
	IsDir bool   `json:"isDir"`
}

// parseJSONArrayOrOne handles Subsonic's inconsistent JSON responses.
// Some endpoints return an array, others return a single object when there's only one item.
// This function normalizes both cases to a slice.
func parseJSONArrayOrOne[T any](raw json.RawMessage) ([]T, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	if raw[0] == '[' {
		var out []T
		return out, json.Unmarshal(raw, &out)
	}
	var one T
	if err := json.Unmarshal(raw, &one); err != nil {
		return nil, err
	}
	return []T{one}, nil
}

// streamTrack downloads a track in its original format (e.g., FLAC).
// Uses the stream API endpoint with format=raw to avoid server-side transcoding.
// Has a longer timeout (15 minutes) to handle large files.
func (c *Client) streamTrack(id string) ([]byte, error) {
	base, err := c.base()
	if err != nil {
		return nil, err
	}
	auth, err := c.authValues()
	if err != nil {
		return nil, err
	}
	auth.Set("id", id)
	auth.Set("format", "raw") // Request original format (FLAC, etc.)
	delete(auth, "f")         // Remove JSON format for binary response
	reqURL := base + "/rest/stream?" + auth.Encode()
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	// Longer timeout for streaming large files
	client := &http.Client{Timeout: 15 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("navidrome: stream HTTP %d: %s", resp.StatusCode, truncate(body, 200))
	}
	// Check for error response disguised as success
	ct := resp.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") || strings.HasPrefix(ct, "text/xml") {
		return nil, fmt.Errorf("navidrome: stream failed: %s", truncate(body, 400))
	}
	return body, nil
}
