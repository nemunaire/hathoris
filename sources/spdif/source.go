package spdif

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"git.nemunai.re/nemunaire/hathoris/alsacontrol"
	"git.nemunai.re/nemunaire/hathoris/sources"
)

type SPDIFSource struct {
	processRec  *exec.Cmd
	processPlay *exec.Cmd
	endChan     chan bool
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
						sr, err := getCardSampleRate("imxspdif")
						if err == nil {
							sources.SoundSources["imxspdif"] = &SPDIFSource{
								endChan:   make(chan bool, 1),
								DeviceIn:  "imxspdif",
								DeviceOut: "is31ap2121",
								Bitrate:   sr,
								Channels:  2,
								Format:    "S24_LE",
							}
						}
					}
					fd.Close()
				}
			}
		}
	}
}

func (s *SPDIFSource) GetName() string {
	if s.IsActive() {
		return fmt.Sprintf("S/PDIF %.1f kHz", float32(s.Bitrate)/1000)
	}
	return "S/PDIF"
}

func (s *SPDIFSource) IsActive() bool {
	return s.processRec != nil
}

func (s *SPDIFSource) IsEnabled() bool {
	return s.processRec != nil || s.processPlay != nil
}

func (s *SPDIFSource) Enable() error {
	if s.processRec != nil {
		return fmt.Errorf("Already running")
	}
	if s.processPlay != nil {
		s.processPlay.Process.Kill()
	}

	// Update bitrate
	sr, err := getCardSampleRate(s.DeviceIn)
	if err != nil {
		return err
	}
	s.Bitrate = sr

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
		}
		pipeR.Close()

		s.processPlay = nil
	}()

	s.processRec = exec.Command("arecord", "-t", "wav", "-f", s.Format, fmt.Sprintf("-r%d", s.Bitrate), fmt.Sprintf("-c%d", s.Channels), "-D", "hw:"+s.DeviceIn, "-F0", "--period-size=512", "-B0", "--buffer-size=512")
	s.processRec.Stdout = pipeW
	if err := s.processRec.Start(); err != nil {
		s.processPlay.Process.Kill()
		return err
	}

	go s.watchBitrate()

	go func() {
		err := s.processRec.Wait()
		if err != nil {
			s.processRec.Process.Kill()
		}
		pipeW.Close()

		s.endChan <- true
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

func getCardSampleRate(cardId string) (sr int64, err error) {
	cc, err := alsa.ParseAmixerContent("hw:" + cardId)
	if err != nil {
		return 0, fmt.Errorf("unable to parse amixer content: %w", err)
	}

	for _, c := range cc {
		if len(c.Values) == 1 {
			val, err := strconv.Atoi(c.Values[0])
			if c.Name == "RX Sample Rate" && err == nil {
				return int64(val), nil
			}
		}
	}

	return 0, fmt.Errorf("unable to find 'RX Sample Rate' control value")
}

func (s *SPDIFSource) watchBitrate() {
	ticker := time.NewTicker(time.Second)

	nbAt0 := 0
loop:
	for {
		select {
		case <-ticker.C:
			sr, err := getCardSampleRate(s.DeviceIn)
			if err == nil {
				if sr == 0 {
					nbAt0 += 1
					if nbAt0 >= 30 {
						log.Printf("[SPDIF] Sample rate is at %d Hz for %d seconds, disabling", sr, nbAt0)
						s.Disable()

						// Wait process exited
						for {
							if s.processPlay == nil && s.processRec == nil {
								break
							}
							time.Sleep(100 * time.Millisecond)
						}
					} else if nbAt0 == 1 || nbAt0%5 == 0 {
						log.Printf("[SPDIF] Sample rate is at %d Hz for %d seconds", sr, nbAt0)
					}
				} else {
					nbAt0 = 0
					if s.Bitrate/10 != sr/10 {
						log.Printf("[SPDIF] Sample rate changes from %d to %d Hz", s.Bitrate, sr)
						s.Bitrate = sr

						s.Disable()

						// Wait process exited
						for {
							if s.processPlay == nil && s.processRec == nil {
								break
							}
							time.Sleep(100 * time.Millisecond)
						}

						s.Enable()
					}
				}
			}
		case <-s.endChan:
			break loop
		}
	}

	ticker.Stop()
}
