# Tick Task

**Tick Task** is a command-line productivity tool: prioritize tasks, work in workspaces, and run a focus timer (Pomodoro-style) from the terminal.

## Requirements

- [Go](https://go.dev/) **1.24.4** or compatible (see `go.mod`)

## Build

From the repository root:

```bash
go build -o ticktask
```

The binary is named `ticktask` by default.

## Install (optional)

The included `install.sh` builds the binary, installs it to `/usr/local/bin/ticktask`, and sets up music directories under `~/.ticktask/music/` if a `music/` tree exists beside the script (with `focus`, `idle`, and `generic` subfolders). Review the script before running it; it uses `sudo` for the binary install and copies.

## Data

Application data lives under **`~/.ticktask/`** (created as needed). Tasks are stored locally (BoltDB).

## Usage overview

Run `ticktask --help` or `ticktask <command> --help` for full flags.

| Command | Description |
|--------|-------------|
| `ticktask add <priority> <name>` | Add a task (integer priority, then name) |
| `ticktask list` | List tasks (`-t` / `--todo`: incomplete only) |
| `ticktask done` | Mark a task complete (interactive picker) |
| `ticktask cancel` | Cancel a task (interactive picker) |
| `ticktask focus` | Pick a task and start the focus timer (`-o` / `--open`: extend past the default 25 minutes) |
| `ticktask version` | Print version (**v0.3.0** as of this tree) |

### Workspaces

Group tasks into named workspaces (default name `default` is used when appropriate):

- `ticktask workspaces new <name>` — create a workspace
- `ticktask workspaces list` — list workspaces (current is marked with `->`)
- `ticktask workspaces select` — choose the active workspace
- `ticktask workspaces move` — move incomplete tasks between workspaces
- `ticktask workspaces remove` — delete a workspace (requires more than one)

### Sync (AWS S3)

Optional backup/sync of the local database to S3:

- `ticktask sync config` — set region, bucket, and credentials (secrets stored via the system keyring where supported)
- `ticktask sync up` — push local DB to the remote backup
- `ticktask sync down` — pull remote backup and overwrite local

## License

This project is licensed under the BSD 3-Clause License; see [LICENSE](LICENSE).
