// Package music implements music configuration commands.
// Allows configuring the music backend (local files or Navidrome server).
package music

import (
	"fmt"

	"github.com/spf13/cobra"
)

// MusicCmd is the parent command for music configuration.
// Subcommands: config
var MusicCmd = &cobra.Command{
	Use:   "music",
	Short: "Configure music playback settings",
	Long: `Configure the music backend for the focus timer.

Subcommands:
  config - Set up music backend (local or Navidrome)

Backends:
  local     - Play MP3/FLAC files from ~/.ticktask/music/
  navidrome - Stream from a Navidrome server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'ticktask music config' to configure music settings")
	},
}
