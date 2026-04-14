// Package config handles loading and caching of application configuration.
// Configuration is stored in BoltDB (non-sensitive) and system keyring (password).
package config

import (
	"os"
	"sync"

	"ticktask/persistence"
)

// NavidromePasswordEnv is the environment variable that overrides the stored password.
// Use this to avoid storing the password in the keyring.
const NavidromePasswordEnv = "TICKTASK_NAVIDROME_PASSWORD"

// Music holds playback source settings.
// Supports two backends:
//   - "local" (default): Reads MP3/FLAC files from ~/.ticktask/music/{focus,idle,generic}/
//   - "navidrome": Streams from a Navidrome server using the Subsonic API
type Music struct {
	Backend   string         // "local" or "navidrome"
	Navidrome NavidromeMusic // Navidrome connection settings
}

// NavidromeMusic configures connection to a Subsonic-compatible server (e.g., Navidrome).
// The Subsonic API uses salt+token authentication with MD5 hashing.
type NavidromeMusic struct {
	BaseURL   string // Server URL (e.g., "http://localhost:4533")
	Username  string // Subsonic username
	Password  string // Subsonic password
	Playlists struct {
		Focus   string // Playlist name for focus mode
		Rest    string // Playlist name for rest mode
		Generic string // Playlist name for generic/chore mode
	}
}

// Mutex and cache for thread-safe config loading.
var (
	musicMu     sync.Mutex
	musicCached *Music
)

// LoadMusic loads and returns music configuration settings.
// The configuration is loaded once and cached for subsequent calls.
//
// Loading priority:
//  1. Reads from BoltDB config bucket
//  2. Falls back to defaults (local backend) if not configured
//  3. Password is loaded from system keyring (or env var override)
//  4. Sets default playlist names if not specified
//
// Returns a copy of the cached config to prevent external modification.
func LoadMusic() (*Music, error) {
	musicMu.Lock()
	defer musicMu.Unlock()
	if musicCached != nil {
		m := *musicCached
		return &m, nil
	}

	db := persistence.GetDB()
	wallet := persistence.GetWallet()

	m := &Music{}

	// Load backend (default: local)
	if backend, err := db.GetConfig(persistence.MusicBackendConfig); err == nil {
		m.Backend = backend
	} else {
		m.Backend = "local"
	}

	// Load Navidrome settings
	if baseURL, err := db.GetConfig(persistence.NavidromeBaseURLConfig); err == nil {
		m.Navidrome.BaseURL = baseURL
	}
	if username, err := db.GetConfig(persistence.NavidromeUsernameConfig); err == nil {
		m.Navidrome.Username = username
	}

	// Load password: env var takes priority, then keyring
	if p := os.Getenv(NavidromePasswordEnv); p != "" {
		m.Navidrome.Password = p
	} else if password, err := wallet.GetKey(persistence.NavidromePasswordKey); err == nil {
		m.Navidrome.Password = password
	}

	// Load playlist names
	if focus, err := db.GetConfig(persistence.NavidromePlaylistFocusConfig); err == nil {
		m.Navidrome.Playlists.Focus = focus
	}
	if rest, err := db.GetConfig(persistence.NavidromePlaylistRestConfig); err == nil {
		m.Navidrome.Playlists.Rest = rest
	}
	if generic, err := db.GetConfig(persistence.NavidromePlaylistGenericConfig); err == nil {
		m.Navidrome.Playlists.Generic = generic
	}

	applyPlaylistDefaults(&m.Navidrome)
	musicCached = m
	out := *musicCached
	return &out, nil
}

// applyPlaylistDefaults sets default playlist names if not specified.
// Defaults: "focus", "rest", "generic"
func applyPlaylistDefaults(n *NavidromeMusic) {
	if n.Playlists.Focus == "" {
		n.Playlists.Focus = "focus"
	}
	if n.Playlists.Rest == "" {
		n.Playlists.Rest = "rest"
	}
	if n.Playlists.Generic == "" {
		n.Playlists.Generic = "generic"
	}
}
