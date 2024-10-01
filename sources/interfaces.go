package sources

import ()

var (
	LoadableSources = map[string]LoadaleSource{}
	SoundSources    = map[string]SoundSource{}
)

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

type LoadaleSource func(map[string]string) (SoundSource, error)
