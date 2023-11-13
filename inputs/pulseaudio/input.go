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

func (s *PulseaudioInput) CurrentlyPlaying() map[string]string {
	cmd := exec.Command("pactl", "-f", "json", "list", "sink-inputs")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Unable to list sink-inputs:", err.Error())
		return nil
	}

	var sinkinputs []PASinkInput
	err = json.Unmarshal(stdoutStderr, &sinkinputs)
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
