package amp1gpio

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"

	"git.nemunai.re/nemunaire/hathoris/sources"
)

type AMP1GPIOSource struct {
	process *exec.Cmd
	Path    string
}

const GPIODirectory = "/sys/class/gpio/gpio46/"

func init() {
	if _, err := os.Stat(GPIODirectory); err == nil {
		sources.SoundSources["amp1"] = &AMP1GPIOSource{
			Path: GPIODirectory,
		}
	}
}

func (s *AMP1GPIOSource) GetName() string {
	return "analog."
}

func (s *AMP1GPIOSource) read() ([]byte, error) {
	fd, err := os.Open(path.Join(s.Path, "value"))
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	return io.ReadAll(fd)
}

func (s *AMP1GPIOSource) IsActive() bool {
	return s.process != nil
}

func (s *AMP1GPIOSource) IsEnabled() bool {
	b, err := s.read()
	if err != nil {
		log.Println("Unable to get amp1 GPIO state:", err.Error())
		return false
	}

	return bytes.Compare(b, []byte{'1', '\n'}) == 0
}

func (s *AMP1GPIOSource) write(value string) error {
	fd, err := os.Create(path.Join(s.Path, "value"))
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.Write([]byte(value))

	return err
}

func (s *AMP1GPIOSource) Enable() error {
	if s.process != nil {
		return fmt.Errorf("Already running")
	}

	s.process = exec.Command("aplay", "-f", "cd", "/dev/zero")
	if err := s.process.Start(); err != nil {
		return err
	}

	go func() {
		err := s.process.Wait()
		if err != nil {
			s.process.Process.Kill()
		}

		s.process = nil
	}()

	return s.write("1")
}

func (s *AMP1GPIOSource) Disable() error {
	if s.process != nil {
		if s.process.Process != nil {
			s.process.Process.Kill()
		}
	}

	return s.write("0")
}
