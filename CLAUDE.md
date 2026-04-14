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
│   └── player.go           # Audio playback (oto + go-mp3)
├── config/
│   └── music.go            # Music backend config (local/navidrome)
├── navidrome/
│   └── client.go           # Subsonic API client for Navidrome
└── utils/
    ├── directories.go      # GetInstallationPath() → ~/.ticktask/
    ├── others.go           # GetRandom(), StringifyTasks()
    └── env.go
```

## Data Location

All data lives in `~/.ticktask/`:
- `tasks.db` - BoltDB database (tasks, workspaces, config)
- `config.yaml` - Music backend configuration (optional)
- `music/` - Local music files (focus/, idle/, generic/)

## Key Flows

### Focus Timer (cmd/focus.go → views/countdown.go)
1. User selects task from list
2. `initCountdown()` creates 3 players (focus, rest, generic)
3. Pomodoro cycle: 25min focus → 5min rest → repeat
4. Keys: Space=toggle focus/rest, Backspace=chore mode, q=quit

### Music Loading (player/player.go)
`loadMusicForPlaylist()` checks `config.LoadMusic().Backend`:
- "local" → reads from ~/.ticktask/music/{focus,idle,generic}/
- "navidrome" → streams from Navidrome server via Subsonic API

## Config File (~/.ticktask/config.yaml)

```yaml
music:
  backend: navidrome  # or "local" (default)
  navidrome:
    base_url: http://localhost:4533
    username: user
    password: pass  # or use TICKTASK_NAVIDROME_PASSWORD env
    playlists:
      focus: "Deep Focus"
      rest: "Chill"
      generic: "Background"
```

## Dependencies

- **cobra/viper** - CLI framework
- **bubbletea** - TUI framework
- **boltdb** - Embedded key-value store
- **oto/go-mp3** - Audio playback
- **go-keyring** - System keyring access
- **aws-sdk-go-v2** - S3 sync

## Testing

No test files currently exist.
