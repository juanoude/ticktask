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

type PlayerStatus string

const (
	PausedStatus  PlayerStatus = "paused"
	PlayingStatus PlayerStatus = "playing"
)

type TTPlayer struct {
	player      *oto.Player
	controlChan chan PlayerCommand
	status      PlayerStatus
}

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

func GetFocusPlayer() *TTPlayer {
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

	return &TTPlayer{
		player:      initializePlayer(fileBytes),
		controlChan: make(chan PlayerCommand),
		status:      PausedStatus,
	}
}

func GetRestPlayer() *TTPlayer {
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

	return &TTPlayer{
		player:      initializePlayer(fileBytes),
		controlChan: make(chan PlayerCommand),
		status:      PausedStatus,
	}
}

func GetGenericPlayer() *TTPlayer {
	path := utils.GetInstallationPath("/music/generic")
	songs, err := utils.ListFilesOnDir(path)
	if err != nil {
		log.Fatal("error picking generic songs")
	}
	chosenFile := utils.GetRandom(songs)
	fileBytes, err := os.ReadFile(path + "/" + chosenFile)
	if err != nil {
		log.Fatal("error picking generic music")
	}

	return &TTPlayer{
		player:      initializePlayer(fileBytes),
		controlChan: make(chan PlayerCommand),
		status:      PausedStatus,
	}
}

func (ttp *TTPlayer) Play() {
	ttp.controlChan <- PlayCommand
}

func (ttp *TTPlayer) Pause() {
	ttp.controlChan <- PauseCommand
}

func (ttp *TTPlayer) Close() {
	ttp.controlChan <- CloseCommand
	close(ttp.controlChan)
}

func (ttp *TTPlayer) InitPlayer() {
	go ttp.initListener()
}

func (ttp *TTPlayer) initListener() {
	for command := range ttp.controlChan {
		switch command {
		case PauseCommand:
			ttp.status = PausedStatus
			ttp.player.Pause()
		case PlayCommand:
			ttp.player.Play()
			ttp.status = PlayingStatus
			go func() {
				time.Sleep(time.Second)
				for ttp.player.IsPlaying() {
					time.Sleep(time.Second)
				}

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
