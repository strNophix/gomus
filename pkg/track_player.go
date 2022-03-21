package gomus

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

const (
	base        = 2
	minVolume   = -5
	maxVolume   = 1
	startVolume = -2
)

type TrackPlayerEffects struct {
	volume float64
}

func newTrackPlayerEffects() TrackPlayerEffects {
	return TrackPlayerEffects{volume: startVolume}
}

type TrackPlayer struct {
	streamer   *beep.StreamSeekCloser
	playerCtrl *beep.Ctrl
	playerVol  *effects.Volume
}

func (t *TrackPlayer) Play(streamer *beep.StreamSeekCloser, trackEffects *TrackPlayerEffects) {
	if t.streamer != nil {
		t.Close()
	}

	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, *streamer), Paused: false}
	volume := &effects.Volume{Streamer: ctrl, Base: base, Volume: trackEffects.volume, Silent: false}
	speaker.Play(volume)

	t.streamer = streamer
	t.playerCtrl = ctrl
	t.playerVol = volume
}

func (t *TrackPlayer) TogglePause() bool {
	speaker.Lock()
	newState := !t.playerCtrl.Paused
	t.playerCtrl.Paused = newState
	speaker.Unlock()
	return newState
}

func (t *TrackPlayer) SetVolume(volume float64) {
	speaker.Lock()
	(*t.playerVol).Volume = volume
	speaker.Unlock()
}

func (t *TrackPlayer) Close() {
	(*t.streamer).Close()
}
