package inputs

import ()

var SoundInputs = map[string]SoundInput{}

type SoundInput interface {
	GetName() string
	IsActive() bool
	CurrentlyPlaying() map[string]string
}

type ControlableInput interface {
	TogglePause(string) error
}
