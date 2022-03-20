package gomus

import (
	"io"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/speaker"
)

type trackPlayer struct {
	currentStream beep.StreamSeekCloser
}

func (t trackPlayer) play(reader io.Reader) {
	if t.currentStream != nil {
		t.currentStream.Close()
	}
	streamer, format, err := flac.Decode(reader)
	check(err)

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	check(err)

	speaker.Play(streamer)
	t.currentStream = streamer
}

func (t trackPlayer) close() {
	t.currentStream.Close()
}
