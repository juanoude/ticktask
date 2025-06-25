package player

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"ticktask/utils"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

var otoCtx *oto.Context
var readyChan chan struct{}

type PlayerCommand string

const (
	PauseCommand PlayerCommand = "pause"
	PlayCommand  PlayerCommand = "play"
	CloseCommand PlayerCommand = "close"
)

func initializePlayer(fileBytes []byte) *oto.Player {
	fileBytesReader := bytes.NewReader(fileBytes)
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	if otoCtx == nil {
		op := &oto.NewContextOptions{}
		op.SampleRate = 44100
		op.ChannelCount = 2
		op.Format = oto.FormatSignedInt16LE
		otoCtx, readyChan, err = oto.NewContext(op)
		if err != nil {
			panic("oto.NewContext failed: " + err.Error())
		}
		<-readyChan
	}

	player := otoCtx.NewPlayer(decodedMp3)
	return player
}

func InitFocusPlayer() *oto.Player {
	path := utils.GetInstallationPath("/music/focus")
	songs, err := utils.ListFilesOnDir(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("error picking focus songs %v - %s", err, path))
	}
	chosenFile := utils.GetRandom(songs)
	fileBytes, err := os.ReadFile(path + "/" + chosenFile)
	if err != nil {
		log.Fatal("error picking focus music")
	}
	return initializePlayer(fileBytes)
}

func InitRestPlayer() *oto.Player {
	path := utils.GetInstallationPath("/music/idle")
	songs, err := utils.ListFilesOnDir(path)
	if err != nil {
		log.Fatal("error picking idle songs")
	}
	chosenFile := utils.GetRandom(songs)
	fileBytes, err := os.ReadFile(path + "/" + chosenFile)
	if err != nil {
		log.Fatal("error picking idle music")
	}
	return initializePlayer(fileBytes)
}

func InitPlayerListener(player *oto.Player, controlChan chan PlayerCommand) {
	status := "playing"

	for command := range controlChan {
		switch command {
		case PauseCommand:
			status = "paused"
			player.Pause()
		case PlayCommand:
			player.Play()
			status = "playing"
			go func() {
				for player.IsPlaying() {
					time.Sleep(time.Second)
				}

				if status == "playing" {
					_, err := player.Seek(0, io.SeekStart)
					if err != nil {
						panic("player.Seek failed: " + err.Error())
					}
					controlChan <- PlayCommand
				}
			}()

		case CloseCommand:
			player.Close()
		}
	}
}
