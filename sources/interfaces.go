package sources

import (
	"encoding/json"
)

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

type LoadaleSource struct {
	LoadSource       func(map[string]interface{}) (SoundSource, error)
	Description      string
	SourceDefinition interface{}
}

func Unmarshal(in map[string]interface{}, out interface{}) error {
	jin, err := json.Marshal(in)
	if err != nil {
		return err
	}

	return json.Unmarshal(jin, out)
}
