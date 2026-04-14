// Package player provides audio playback functionality for the focus timer.
// It uses the oto library for cross-platform audio output with support for
// both MP3 (local files) and FLAC (Navidrome streaming) formats.
//
// The package supports three music modes corresponding to timer states:
//   - Focus: Background music for concentrated work
//   - Rest: Music for break periods
//   - Generic: Music for general/chore tasks
package player

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"ticktask/config"
	"ticktask/navidrome"
	"ticktask/utils"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
	"github.com/mewkiz/flac"
)

// otoCtx is the shared audio context for all players.
// Initialized once on first player creation.
var otoCtx *oto.Context

// readyChan signals when the audio context is ready for playback.
var readyChan chan struct{}

// PlayerCommand represents control commands for the audio player.
type PlayerCommand string

const (
	PauseCommand PlayerCommand = "pause" // Pause playback
	PlayCommand  PlayerCommand = "play"  // Start/resume playback
	CloseCommand PlayerCommand = "close" // Stop and release resources
)

// PlayerStatus represents the current playback state.
type PlayerStatus string

const (
	PausedStatus  PlayerStatus = "paused"  // Playback is paused
	PlayingStatus PlayerStatus = "playing" // Playback is active
)

// TTPlayer wraps an oto.Player with channel-based control.
// Supports play, pause, and close operations via the control channel.
// Automatically loops the audio when playback completes.
type TTPlayer struct {
	player      *oto.Player         // The underlying audio player
	controlChan chan PlayerCommand  // Channel for receiving control commands
	status      PlayerStatus        // Current playback state
}

// audioFormat represents the detected audio file format.
type audioFormat int

const (
	formatUnknown audioFormat = iota
	formatMP3
	formatFLAC
)

// detectFormat identifies the audio format from the file header bytes.
// FLAC files start with "fLaC" magic bytes.
// MP3 files start with ID3 tag or 0xFF sync byte.
func detectFormat(data []byte) audioFormat {
	if len(data) < 4 {
		return formatUnknown
	}
	// FLAC: starts with "fLaC"
	if data[0] == 'f' && data[1] == 'L' && data[2] == 'a' && data[3] == 'C' {
		return formatFLAC
	}
	// MP3: ID3 tag or sync word
	if (data[0] == 'I' && data[1] == 'D' && data[2] == '3') || (data[0] == 0xFF && (data[1]&0xE0) == 0xE0) {
		return formatMP3
	}
	return formatUnknown
}

// flacReader wraps decoded FLAC audio as an io.ReadSeeker for oto.
// Converts FLAC samples to 16-bit signed little-endian PCM.
type flacReader struct {
	data       []byte       // original FLAC data for seeking
	stream     *flac.Stream
	sampleRate uint32
	channels   uint8
	bps        uint8  // bits per sample
	buf        []byte // buffer of decoded PCM bytes
	pos        int    // current position in buf
}

// newFLACReader creates a reader from FLAC data bytes.
func newFLACReader(data []byte) (*flacReader, error) {
	stream, err := flac.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("flac.Parse: %w", err)
	}
	return &flacReader{
		data:       data,
		stream:     stream,
		sampleRate: stream.Info.SampleRate,
		channels:   stream.Info.NChannels,
		bps:        stream.Info.BitsPerSample,
	}, nil
}

// Read implements io.Reader by decoding FLAC frames to PCM.
func (r *flacReader) Read(p []byte) (n int, err error) {
	// Return buffered data first
	if r.pos < len(r.buf) {
		n = copy(p, r.buf[r.pos:])
		r.pos += n
		return n, nil
	}

	// Decode next frame
	frame, err := r.stream.ParseNext()
	if err != nil {
		return 0, err
	}

	// Convert samples to 16-bit PCM
	r.buf = r.buf[:0]
	r.pos = 0
	nSamples := len(frame.Subframes[0].Samples)
	for i := 0; i < nSamples; i++ {
		for ch := 0; ch < int(r.channels); ch++ {
			sample := frame.Subframes[ch].Samples[i]
			// Scale to 16-bit based on bits per sample
			if r.bps > 16 {
				sample >>= (r.bps - 16)
			} else if r.bps < 16 {
				sample <<= (16 - r.bps)
			}
			// Clamp to int16 range
			if sample > 32767 {
				sample = 32767
			} else if sample < -32768 {
				sample = -32768
			}
			var buf [2]byte
			binary.LittleEndian.PutUint16(buf[:], uint16(int16(sample)))
			r.buf = append(r.buf, buf[:]...)
		}
	}

	n = copy(p, r.buf[r.pos:])
	r.pos += n
	return n, nil
}

// Seek implements io.Seeker (limited: only supports seeking to start for looping).
func (r *flacReader) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekStart && offset == 0 {
		// Re-parse stream from beginning for looping
		stream, err := flac.Parse(bytes.NewReader(r.data))
		if err != nil {
			return 0, fmt.Errorf("flacReader: failed to re-parse for seek: %w", err)
		}
		r.stream = stream
		r.pos = 0
		r.buf = r.buf[:0]
		return 0, nil
	}
	return 0, fmt.Errorf("flacReader: only seek to start is supported")
}

// initializePlayer creates an oto.Player from audio bytes (MP3 or FLAC).
// Automatically detects the format and uses the appropriate decoder.
// Initializes the shared audio context on first call.
// Panics if decoding fails or audio context cannot be created.
func initializePlayer(fileBytes []byte) *oto.Player {
	format := detectFormat(fileBytes)

	var reader io.Reader
	var sampleRate int

	switch format {
	case formatFLAC:
		flacRdr, err := newFLACReader(fileBytes)
		if err != nil {
			panic("FLAC decode failed: " + err.Error())
		}
		reader = flacRdr
		sampleRate = int(flacRdr.sampleRate)

	case formatMP3:
		fileBytesReader := bytes.NewReader(fileBytes)
		decodedMp3, err := mp3.NewDecoder(fileBytesReader)
		if err != nil {
			panic("mp3.NewDecoder failed: " + err.Error())
		}
		reader = decodedMp3
		sampleRate = decodedMp3.SampleRate()

	default:
		panic("unsupported audio format")
	}

	// Initialize audio context once (stereo, 16-bit signed LE)
	// LIMITATION: Sample rate is locked to the first track's rate.
	// All subsequent tracks should have the same sample rate for correct playback.
	// Most music libraries use 44.1kHz or 48kHz consistently.
	if otoCtx == nil {
		var err error
		op := &oto.NewContextOptions{}
		op.SampleRate = sampleRate
		op.ChannelCount = 2
		op.Format = oto.FormatSignedInt16LE
		otoCtx, readyChan, err = oto.NewContext(op)
		if err != nil {
			panic("oto.NewContext failed: " + err.Error())
		}
		<-readyChan
	}

	player := otoCtx.NewPlayer(reader)
	return player
}

// GetFocusPlayer creates a player for focus/work mode music.
// Loads a random track from the "focus" playlist (Navidrome) or ~/.ticktask/music/focus/ (local).
func GetFocusPlayer() *TTPlayer {
	fileBytes, err := loadMusicForPlaylist("/music/focus", func(m *config.Music) string {
		return m.Navidrome.Playlists.Focus
	})
	if err != nil {
		log.Fatal(err)
	}
	return newTTPlayer(fileBytes)
}

// GetRestPlayer creates a player for rest/break mode music.
// Loads a random track from the "rest" playlist (Navidrome) or ~/.ticktask/music/idle/ (local).
func GetRestPlayer() *TTPlayer {
	fileBytes, err := loadMusicForPlaylist("/music/idle", func(m *config.Music) string {
		return m.Navidrome.Playlists.Rest
	})
	if err != nil {
		log.Fatal(err)
	}
	return newTTPlayer(fileBytes)
}

// GetGenericPlayer creates a player for generic/chore mode music.
// Loads a random track from the "generic" playlist (Navidrome) or ~/.ticktask/music/generic/ (local).
func GetGenericPlayer() *TTPlayer {
	fileBytes, err := loadMusicForPlaylist("/music/generic", func(m *config.Music) string {
		return m.Navidrome.Playlists.Generic
	})
	if err != nil {
		log.Fatal(err)
	}
	return newTTPlayer(fileBytes)
}

// newTTPlayer wraps raw audio bytes (MP3 or FLAC) in a TTPlayer with control channel.
// The player starts in paused state.
func newTTPlayer(fileBytes []byte) *TTPlayer {
	return &TTPlayer{
		player:      initializePlayer(fileBytes),
		controlChan: make(chan PlayerCommand),
		status:      PausedStatus,
	}
}

// loadMusicForPlaylist loads a random track based on the configured backend.
// For "navidrome" backend: streams from the Navidrome server using the Subsonic API.
// For "local" backend (default): reads from the local music directory.
func loadMusicForPlaylist(localRel string, playlistName func(*config.Music) string) ([]byte, error) {
	m, err := config.LoadMusic()
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(strings.TrimSpace(m.Backend), "navidrome") {
		return navidrome.RandomTrackFromPlaylist(&m.Navidrome, playlistName(m))
	}
	return loadLocalDir(localRel)
}

// loadLocalDir reads a random audio file (MP3 or FLAC) from the specified local directory.
// The path is relative to ~/.ticktask/ (e.g., "/music/focus" → ~/.ticktask/music/focus/).
func loadLocalDir(rel string) ([]byte, error) {
	path := utils.GetInstallationPath(rel)
	songs, err := utils.ListFilesOnDir(path)
	if err != nil {
		return nil, fmt.Errorf("local music %s: %w", path, err)
	}
	chosenFile := utils.GetRandom(songs)
	return os.ReadFile(filepath.Join(path, chosenFile))
}

// Play sends a play command to start or resume playback.
func (ttp *TTPlayer) Play() {
	ttp.controlChan <- PlayCommand
}

// Pause sends a pause command to pause playback.
func (ttp *TTPlayer) Pause() {
	ttp.controlChan <- PauseCommand
}

// Close sends a close command and closes the control channel.
// Should be called when the player is no longer needed.
func (ttp *TTPlayer) Close() {
	ttp.controlChan <- CloseCommand
	close(ttp.controlChan)
}

// InitPlayer starts the background goroutine that listens for control commands.
// Must be called before Play(), Pause(), or Close().
func (ttp *TTPlayer) InitPlayer() {
	go ttp.initListener()
}

// initListener is the main control loop that processes player commands.
// Runs in a separate goroutine and handles:
//   - PauseCommand: Pauses playback
//   - PlayCommand: Starts playback and sets up auto-loop when track ends
//   - CloseCommand: Releases player resources
func (ttp *TTPlayer) initListener() {
	for command := range ttp.controlChan {
		switch command {
		case PauseCommand:
			ttp.status = PausedStatus
			ttp.player.Pause()

		case PlayCommand:
			ttp.player.Play()
			ttp.status = PlayingStatus
			// Monitor playback and loop when track ends
			go func() {
				time.Sleep(time.Second)
				for ttp.player.IsPlaying() {
					time.Sleep(time.Second)
				}

				// If still in playing status when track ends, loop
				if ttp.status == "playing" {
					_, err := ttp.player.Seek(0, io.SeekStart)
					if err != nil {
						panic("player.Seek failed: " + err.Error())
					}
					ttp.controlChan <- PlayCommand
				}
			}()

		case CloseCommand:
			ttp.player.Close()
		}
	}
}
