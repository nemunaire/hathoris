package sources

import ()

var SoundSources = map[string]SoundSource{}

type SoundSource interface {
	GetName() string
	IsActive() bool
	IsEnabled() bool
	Enable() error
	Disable() error
}

type PlayingSource interface {
	CurrentlyPlaying() string
}
