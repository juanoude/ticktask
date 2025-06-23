package player

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

func InitFocusPlayer() *oto.Player {
	fileBytes, err := os.ReadFile("./music/focus/ch-belight.mp3")
	fileBytesReader := bytes.NewReader(fileBytes)
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	op := &oto.NewContextOptions{}
	op.SampleRate = 44100
	op.ChannelCount = 2
	op.Format = oto.FormatSignedInt16LE
	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	<-readyChan
	player := otoCtx.NewPlayer(decodedMp3)
	return player
}

func PlayLoop(player *oto.Player) {
	player.Play()

	for player.IsPlaying() {
		time.Sleep(time.Second)
	}

	newPos, err := player.Seek(0, io.SeekStart)
	if err != nil {
		panic("player.Seek failed: " + err.Error())
	}
	println("Player is now at position:", newPos)
	player.Play()
}
