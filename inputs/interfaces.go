package inputs

import ()

var SoundInputs = map[string]SoundInput{}

type SoundInput interface {
	GetName() string
	IsActive() bool
	CurrentlyPlaying() *string
}

type ControlableInput interface {
	TogglePause() error
}
