package music

import (
	"fmt"
	"ticktask/persistence"
	"ticktask/views"

	"github.com/spf13/cobra"
)

func init() {
	MusicCmd.AddCommand(configCmd)
}

// configCmd configures the music backend and Navidrome connection.
// Prompts for backend choice, and if Navidrome is selected, prompts for
// server URL, username, password (stored in keyring), and playlist names.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure music backend settings",
	Long: `Configure the music playback backend.

For local backend: Music files are read from ~/.ticktask/music/{focus,idle,generic}/

For Navidrome backend: Streams music from a Navidrome server.
  - Server URL, username, and playlist names are stored in the local database
  - Password is stored securely in the system keyring`,
	Run: func(cmd *cobra.Command, args []string) {
		db := persistence.GetDB()
		wallet := persistence.GetWallet()

		// Select backend
		backends := []string{"local", "navidrome"}
		backendIndex := views.RunSelector(backends, "Select music backend:")
		if backendIndex < 0 {
			return
		}
		backend := backends[backendIndex]

		if err := db.StoreConfig(persistence.MusicBackendConfig, backend); err != nil {
			fmt.Printf("Error saving backend: %v\n", err)
			return
		}

		if backend == "local" {
			fmt.Println("Music backend set to local.")
			fmt.Println("Place your music files in ~/.ticktask/music/{focus,idle,generic}/")
			return
		}

		// Navidrome configuration
		baseURL, cancelled := views.RunInput("Navidrome server URL (e.g., http://localhost:4533):", false)
		if cancelled {
			return
		}
		if err := db.StoreConfig(persistence.NavidromeBaseURLConfig, baseURL); err != nil {
			fmt.Printf("Error saving base URL: %v\n", err)
			return
		}

		username, cancelled := views.RunInput("Navidrome username:", false)
		if cancelled {
			return
		}
		if err := db.StoreConfig(persistence.NavidromeUsernameConfig, username); err != nil {
			fmt.Printf("Error saving username: %v\n", err)
			return
		}

		password, cancelled := views.RunInput("Navidrome password:", true)
		if cancelled {
			return
		}
		if err := wallet.StoreKey(persistence.NavidromePasswordKey, password); err != nil {
			fmt.Printf("Error saving password to keyring: %v\n", err)
			return
		}

		// Playlist names (with defaults)
		fmt.Println("\nPlaylist names (press Enter to use defaults):")

		focusPlaylist, cancelled := views.RunInput("Focus playlist name [default: focus]:", false)
		if cancelled {
			return
		}
		if focusPlaylist == "" {
			focusPlaylist = "focus"
		}
		db.StoreConfig(persistence.NavidromePlaylistFocusConfig, focusPlaylist)

		restPlaylist, cancelled := views.RunInput("Rest playlist name [default: rest]:", false)
		if cancelled {
			return
		}
		if restPlaylist == "" {
			restPlaylist = "rest"
		}
		db.StoreConfig(persistence.NavidromePlaylistRestConfig, restPlaylist)

		genericPlaylist, cancelled := views.RunInput("Generic playlist name [default: generic]:", false)
		if cancelled {
			return
		}
		if genericPlaylist == "" {
			genericPlaylist = "generic"
		}
		db.StoreConfig(persistence.NavidromePlaylistGenericConfig, genericPlaylist)

		fmt.Println("\nNavidrome configuration saved successfully!")
	},
}
