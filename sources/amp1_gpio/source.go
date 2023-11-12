package amp1gpio

import (
	"bytes"
	"io"
	"log"
	"os"
	"path"

	"git.nemunai.re/nemunaire/hathoris/sources"
)

type AMP1GPIOSource struct {
	Path string
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
	return "entrée analogique"
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
	return s.IsEnabled()
}

func (s *AMP1GPIOSource) IsEnabled() bool {
	b, err := s.read()
	if err != nil {
		log.Println("Unable to get amp1 GPIO state:", err.Error())
		return false
	}

	return bytes.Compare(b, []byte{'1'}) == 0
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
	return s.write("1")
}

func (s *AMP1GPIOSource) Disable() error {
	return s.write("0")
}
