package gomus

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

type TrackPlayer struct {
	streamer   *beep.StreamSeekCloser
	playerCtrl *beep.Ctrl
}

func (t *TrackPlayer) Play(streamer *beep.StreamSeekCloser) {
	if t.streamer != nil {
		t.Close()
	}

	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, *streamer), Paused: false}
	speaker.Play(ctrl)

	t.playerCtrl = ctrl
	t.streamer = streamer
}

func (t *TrackPlayer) TogglePause() bool {
	speaker.Lock()
	newState := !t.playerCtrl.Paused
	t.playerCtrl.Paused = newState
	speaker.Unlock()
	return newState
}

func (t *TrackPlayer) Close() {
	(*t.streamer).Close()
}
