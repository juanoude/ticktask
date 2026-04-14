# TickTask

CLI productivity tool: task management with Pomodoro-style focus timer.

## Build & Run

```bash
go build -o ticktask
./ticktask --help
```

## Project Structure

```
ticktask/
├── main.go                 # Entry point → cmd.Execute()
├── cmd/                    # Cobra CLI commands
│   ├── root.go             # Root command, registers subcommands
│   ├── add.go              # ticktask add <priority> <name>
│   ├── list.go             # ticktask list [-t for todo only]
│   ├── done.go             # ticktask done (interactive)
│   ├── cancel.go           # ticktask cancel (interactive)
│   ├── focus.go            # ticktask focus [-o for open-ended]
│   ├── version.go
│   ├── workspace/          # ticktask workspaces <subcommand>
│   ├── music/              # ticktask music config
│   └── sync/               # ticktask sync <up|down|config>
├── models/
│   └── task.go             # Task struct
├── persistence/
│   ├── interface.go        # PersistenceLayer, WalletLayer, SyncLayer interfaces
│   ├── bolt/               # BoltDB implementation (tasks, workspaces, config)
│   ├── gkeyring/           # System keyring for secrets
│   └── sync/amazon/        # S3 backup implementation
├── views/                  # Bubble Tea TUI components
│   ├── selector.go         # Interactive list picker
│   ├── input.go            # Text input
│   └── countdown.go        # Focus timer (uses player)
├── player/
│   └── player.go           # Audio playback (oto + go-mp3/flac)
├── config/
│   └── music.go            # Music config loader (from DB/keyring)
├── navidrome/
│   └── client.go           # Subsonic API client for Navidrome
└── utils/
    ├── directories.go      # GetInstallationPath() → ~/.ticktask/
    ├── others.go           # GetRandom(), StringifyTasks()
    └── env.go
```

## Data Location

All data lives in `~/.ticktask/`:
- `data/ticktask.db` - BoltDB database (tasks, workspaces, config)
- `music/` - Local music files (focus/, idle/, generic/)

Sensitive credentials (AWS keys, Navidrome password) are stored in the system keyring.

## Key Flows

### Focus Timer (cmd/focus.go → views/countdown.go)
1. User selects task from list
2. `initCountdown()` creates 3 players (focus, rest, generic)
3. Pomodoro cycle: 25min focus → 5min rest → repeat
4. Keys: Space=toggle focus/rest, Backspace=chore mode, q=quit

### Music Loading (player/player.go)
`loadMusicForPlaylist()` checks `config.LoadMusic().Backend`:
- "local" → reads MP3/FLAC files from ~/.ticktask/music/{focus,idle,generic}/
- "navidrome" → streams FLAC from Navidrome server via Subsonic API (raw format)

Audio format is auto-detected from file headers (FLAC: "fLaC", MP3: ID3/sync).

## Configuration

Configuration is stored in BoltDB (non-sensitive) and the system keyring (secrets).

### Music Configuration
Run `ticktask music config` to configure:
- Backend: "local" or "navidrome"
- Navidrome settings: URL, username, password, playlist names

Password can also be set via `TICKTASK_NAVIDROME_PASSWORD` environment variable.

### Sync Configuration
Run `ticktask sync config` to configure:
- AWS region, bucket name (stored in DB)
- AWS credentials (stored in system keyring)

## Dependencies

- **cobra/viper** - CLI framework
- **bubbletea** - TUI framework
- **boltdb** - Embedded key-value store
- **oto** - Cross-platform audio output
- **go-mp3** - MP3 decoding (local files)
- **mewkiz/flac** - FLAC decoding (Navidrome streaming)
- **go-keyring** - System keyring access
- **aws-sdk-go-v2** - S3 sync

## Testing

No test files currently exist.
