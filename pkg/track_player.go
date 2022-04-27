package gomus

import (
	"time"

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

func NewTrackPlayerEffects() TrackPlayerEffects {
	return TrackPlayerEffects{volume: startVolume}
}

type TrackPlayer struct {
	*beep.Ctrl
	*effects.Volume

	beep.SampleRate
}

func NewTrackPlayer() TrackPlayer {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))

	ctrl := &beep.Ctrl{Streamer: beep.Silence(1), Paused: false}
	volume := &effects.Volume{Streamer: ctrl, Base: base, Volume: startVolume, Silent: false}

	return TrackPlayer{
		SampleRate: sr,
		Ctrl:       ctrl,
		Volume:     volume,
	}
}

func (t *TrackPlayer) Play(streamer beep.Streamer) {
	speaker.Clear()

	speaker.Lock()
	t.Ctrl.Streamer = streamer
	speaker.Unlock()

	speaker.Play(t.Volume)
}

func (t *TrackPlayer) TogglePause() bool {
	speaker.Lock()
	newState := !t.Ctrl.Paused
	t.Ctrl.Paused = newState
	speaker.Unlock()
	return newState
}

func (t *TrackPlayer) SetVolume(volume float64) {
	speaker.Lock()
	t.Volume.Volume = volume
	speaker.Unlock()
}
