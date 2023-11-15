package pa

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"git.nemunai.re/nemunaire/hathoris/inputs"
)

type PulseaudioInput struct {
}

type PASink struct {
	Index               int    `json:"index"`
	Name                string `json:"name"`
	Driver              string `json:"driver"`
	SampleSpecification string `json:"sample_specifications"`
	State               string `json:"state"`
}

type PAVolume struct {
	Value        uint   `json:"value"`
	ValuePercent string `json:"value_percent"`
	DB           string `json:"db"`
}

type PASinkInput struct {
	Index               int64               `json:"index"`
	Driver              string              `json:"driver"`
	OwnerModule         string              `json:"owner_module"`
	Client              string              `json:"client"`
	Sink                int                 `json:"sink"`
	SampleSpecification string              `json:"sample_specifications"`
	ChannelMap          string              `json:"channel_map"`
	Format              string              `json:"format"`
	Corked              bool                `json:"corked"`
	Mute                bool                `json:"mute"`
	Volume              map[string]PAVolume `json:"volume"`
	Balance             float64             `json:"balance"`
	BufferLatencyUsec   float64             `json:"buffer_latency_usec"`
	SinkLatencyUsec     float64             `json:"sink_latency_usec"`
	ResampleMethod      string              `json:"resample_method"`
	Properties          map[string]string   `json:"properties"`
}

func init() {
	cmd := exec.Command("pactl", "-f", "json", "list", "sinks", "short")
	err := cmd.Run()
	if err == nil {
		inputs.SoundInputs["pulseaudio"] = &PulseaudioInput{}
	} else {
		log.Println("Unable to access pulseaudio:", err.Error())
	}
}

func (s *PulseaudioInput) GetName() string {
	return "pulseaudio"
}

func (s *PulseaudioInput) IsActive() bool {
	cmd := exec.Command("pactl", "-f", "json", "list", "sinks", "short")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	var sinks []PASink
	err = json.Unmarshal(stdoutStderr, &sinks)
	if err != nil {
		return false
	}

	for _, sink := range sinks {
		if sink.State != "SUSPENDED" {
			return true
		}
	}

	return false
}

func (s *PulseaudioInput) getPASinkInputs() ([]PASinkInput, error) {
	cmd := exec.Command("pactl", "-f", "json", "list", "sink-inputs")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("unable to list sink-inputs: %w", err)
	}

	var sinkinputs []PASinkInput
	err = json.Unmarshal(stdoutStderr, &sinkinputs)
	if err != nil {
		return nil, fmt.Errorf("unable to parse sink-inputs list: %w", err)
	}

	return sinkinputs, nil
}

func (s *PulseaudioInput) CurrentlyPlaying() map[string]string {
	sinkinputs, err := s.getPASinkInputs()
	if err != nil {
		log.Println("Unable to list sink-inputs:", err.Error())
		return nil
	}

	ret := map[string]string{}
	for _, input := range sinkinputs {
		if v, ok := input.Properties["media.name"]; ok {
			ret[strconv.FormatInt(input.Index, 10)] = v
		} else if v, ok := input.Properties["device.description"]; ok {
			ret[strconv.FormatInt(input.Index, 10)] = v
		} else {
			ret[strconv.FormatInt(input.Index, 10)] = fmt.Sprintf("#%d", input.Index)
		}
	}

	return ret
}

func (s *PulseaudioInput) GetMixers() (map[string]*inputs.InputMixer, error) {
	sinkinputs, err := s.getPASinkInputs()
	if err != nil {
		return nil, err
	}

	ret := map[string]*inputs.InputMixer{}
	for _, input := range sinkinputs {
		var maxvolume string
		for k, vol := range input.Volume {
			if maxvolume == "" || vol.Value > input.Volume[maxvolume].Value {
				maxvolume = k
			}
		}

		ret[strconv.FormatInt(input.Index, 10)] = &inputs.InputMixer{
			Volume:        input.Volume[maxvolume].Value,
			VolumePercent: input.Volume[maxvolume].ValuePercent,
			VolumeDB:      input.Volume[maxvolume].DB,
			Balance:       input.Balance,
			Mute:          input.Mute,
		}
	}

	return ret, nil
}

func (s *PulseaudioInput) SetMixer(stream string, volume *inputs.InputMixer) error {
	sinkinputs, err := s.getPASinkInputs()
	if err != nil {
		return err
	}

	for _, input := range sinkinputs {
		if strconv.FormatInt(input.Index, 10) == stream {
			cmd := exec.Command("pactl", "set-sink-input-volume", stream, strconv.FormatUint(uint64(volume.Volume), 10))
			err := cmd.Run()
			if err != nil {
				return fmt.Errorf("unable to set volume: %w", err)
			}

			if input.Mute != volume.Mute {
				cmd := exec.Command("pactl", "set-sink-input-mute", stream, "toggle")
				err := cmd.Run()
				if err != nil {
					return fmt.Errorf("unable to change mute state: %w", err)
				}
			}

			return nil
		}
	}

	return fmt.Errorf("unable to find stream %q", stream)
}
