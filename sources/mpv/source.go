package mpv

import (
	"fmt"
	"os/exec"

	"git.nemunai.re/nemunaire/hathoris/sources"
)

type MPVSource struct {
	process *exec.Cmd
	Options []string
	File    string
}

func init() {
	sources.SoundSources["mpv"] = &MPVSource{
		Options: []string{"--no-video"},
		File:    "https://mediaserv38.live-streams.nl:18030/stream",
	}
}

func (s *MPVSource) GetName() string {
	return "Radio 1"
}

func (s *MPVSource) IsActive() bool {
	return s.process != nil
}

func (s *MPVSource) IsEnabled() bool {
	return s.process != nil
}

func (s *MPVSource) Enable() (err error) {
	if s.process != nil {
		return fmt.Errorf("Already running")
	}

	var opts []string
	opts = append(opts, s.Options...)
	opts = append(opts, s.File)

	s.process = exec.Command("mpv", opts...)
	if err = s.process.Start(); err != nil {
		return
	}

	go func() {
		err := s.process.Wait()
		if err != nil {
			s.process.Process.Kill()
		}

		s.process = nil
	}()

	return
}

func (s *MPVSource) Disable() error {
	if s.process != nil {
		if s.process.Process != nil {
			s.process.Process.Kill()
		}
	}

	return nil
}
