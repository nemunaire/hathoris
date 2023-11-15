package inputs

import ()

var SoundInputs = map[string]SoundInput{}

type InputMixer struct {
	Volume        uint    `json:"volume"`
	VolumePercent string  `json:"volume_percent"`
	VolumeDB      string  `json:"volume_db"`
	Mute          bool    `json:"mute"`
	Balance       float64 `json:"balance"`
}

type SoundInput interface {
	GetName() string
	IsActive() bool
	CurrentlyPlaying() map[string]string
}

type ControlableInput interface {
	TogglePause(string) error
}

type MixableInput interface {
	GetMixers() (map[string]*InputMixer, error)
	SetMixer(string, *InputMixer) error
}
