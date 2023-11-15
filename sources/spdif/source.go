package spdif

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"git.nemunai.re/nemunaire/hathoris/sources"
)

type SPDIFSource struct {
	processRec  *exec.Cmd
	processPlay *exec.Cmd
	DeviceIn    string
	DeviceOut   string
	Bitrate     int64
	Channels    int64
	Format      string
}

func init() {
	if dirs, err := os.ReadDir("/sys/class/sound"); err == nil {
		for _, dir := range dirs {
			thisdir := path.Join("/sys/class/sound", dir.Name())
			if s, err := os.Stat(thisdir); err == nil && s.IsDir() {
				idfile := path.Join(thisdir, "id")
				if fd, err := os.Open(idfile); err == nil {
					if cnt, err := io.ReadAll(fd); err == nil && string(cnt) == "imxspdif\n" {
						sources.SoundSources["imxspdif"] = &SPDIFSource{
							DeviceIn:  "imxspdif",
							DeviceOut: "is31ap2121",
							Bitrate:   48000,
							Channels:  2,
							Format:    "S24_LE",
						}
					}
					fd.Close()
				}
			}
		}
	}
}

func (s *SPDIFSource) GetName() string {
	return "S/PDIF"
}

func (s *SPDIFSource) IsActive() bool {
	return s.processRec != nil
}

func (s *SPDIFSource) IsEnabled() bool {
	return s.processRec != nil
}

func (s *SPDIFSource) Enable() error {
	if s.processRec != nil {
		return fmt.Errorf("Already running")
	}
	if s.processPlay != nil {
		s.processPlay.Process.Kill()
	}

	pipeR, pipeW, err := os.Pipe()
	if err != nil {
		return err
	}

	s.processPlay = exec.Command("aplay", "-c", fmt.Sprintf("%d", s.Channels), "-D", "hw:"+s.DeviceOut, "--period-size=512", "-B0", "--buffer-size=512")
	s.processPlay.Stdin = pipeR
	if err := s.processPlay.Start(); err != nil {
		return err
	}

	go func() {
		err := s.processPlay.Wait()
		if err != nil {
			if s.processPlay != nil && s.processPlay.Process != nil {
				s.processPlay.Process.Kill()
			}
			pipeR.Close()
			pipeW.Close()
		}

		s.processPlay = nil
	}()

	s.processRec = exec.Command("arecord", "-t", "wav", "-f", s.Format, fmt.Sprintf("-r%d", s.Bitrate), fmt.Sprintf("-c%d", s.Channels), "-D", "hw:"+s.DeviceIn, "-B0", "--buffer-size=512")
	s.processRec.Stdout = pipeW
	if err := s.processRec.Start(); err != nil {
		s.processPlay.Process.Kill()
		return err
	}

	go func() {
		err := s.processRec.Wait()
		if err != nil {
			s.processRec.Process.Kill()
		}

		s.processRec = nil
	}()

	return nil
}

func (s *SPDIFSource) Disable() error {
	if s.processRec != nil && s.processRec.Process != nil {
		s.processRec.Process.Kill()
	}
	if s.processPlay != nil && s.processPlay.Process != nil {
		s.processPlay.Process.Kill()
	}

	return nil
}
