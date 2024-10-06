package mpv

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/DexterLB/mpvipc"

	"git.nemunai.re/nemunaire/hathoris/sources"
)

type MPVSource struct {
	process      *exec.Cmd
	ipcSocketDir string
	Name         string
	Options      []string
	File         string
}

func init() {
	sources.LoadableSources["mpv"] = NewMPVSource
}

func NewMPVSource(kv map[string]string) (sources.SoundSource, error) {
	var s MPVSource

	if name, ok := kv["name"]; ok {
		s.Name = name
	}

	if opts, ok := kv["opts"]; ok {
		s.Options = strings.Split(opts, " ")
	}

	if file, ok := kv["file"]; ok {
		s.File = file
	}

	return &s, nil
}

func (s *MPVSource) GetName() string {
	return s.Name
}

func (s *MPVSource) IsActive() bool {
	return s.process != nil
}

func (s *MPVSource) IsEnabled() bool {
	return s.process != nil
}

func (s *MPVSource) ipcSocket() string {
	return path.Join(s.ipcSocketDir, "mpv.socket")
}

func (s *MPVSource) Enable() (err error) {
	if s.process != nil {
		return fmt.Errorf("Already running")
	}

	s.ipcSocketDir, err = os.MkdirTemp("", "hathoris")

	opts := append([]string{"--no-video", "--no-terminal", "--prefetch-playlist=yes"}, s.Options...)
	if s.ipcSocketDir != "" {
		opts = append(opts, "--input-ipc-server="+s.ipcSocket(), "--pause")
	}
	opts = append(opts, s.File)

	s.process = exec.Command("mpv", opts...)
	if err = s.process.Start(); err != nil {
		log.Println("Unable to launch mpv:", err.Error())
		return
	}

	go func() {
		err := s.process.Wait()
		if err != nil {
			var exiterr *exec.ExitError
			if errors.As(err, &exiterr) {
				if exiterr.ExitCode() > 0 {
					log.Printf("mpv exited with error code = %d", exiterr.ExitCode())
				} else {
					log.Print("mpv exited successfully")
				}
			} else {
				s.process.Process.Kill()
			}
		}

		if s.ipcSocketDir != "" {
			os.RemoveAll(s.ipcSocketDir)
		}

		s.process = nil
	}()

	if s.ipcSocketDir != "" {
		_, err = os.Stat(s.ipcSocket())
		for i := 20; i >= 0 && err != nil; i-- {
			time.Sleep(100 * time.Millisecond)
			_, err = os.Stat(s.ipcSocket())
		}
		time.Sleep(200 * time.Millisecond)

		conn := mpvipc.NewConnection(s.ipcSocket())
		err = conn.Open()
		for i := 20; i >= 0 && err != nil; i-- {
			time.Sleep(100 * time.Millisecond)
			err = conn.Open()
		}
		if err != nil {
			log.Println("Unable to connect to mpv socket:", err.Error())
			return err
		}
		defer conn.Close()

		_, err = conn.Get("media-title")
		for err != nil {
			time.Sleep(100 * time.Millisecond)
			_, err = conn.Get("media-title")
		}

		conn.Set("ao-volume", 1)

		err = conn.Set("pause", false)
		if err != nil {
			log.Println("Unable to unpause:", err.Error())
			return err
		}

		var pfc interface{}
		pfc, err = conn.Get("core-idle")

		for err == nil && pfc.(bool) {
			time.Sleep(250 * time.Millisecond)
			pfc, err = conn.Get("core-idle")
		}

		if err != nil {
			log.Println("Unable to retrieve core-idle status:", err.Error())
		}

		s.FadeIn(conn, 3, 50)
	}

	return
}

func (s *MPVSource) FadeIn(conn *mpvipc.Connection, speed int, level int) {
	volume, err := conn.Get("ao-volume")
	if err != nil {
		volume = 1.0
	}

	for i := int(volume.(float64)) + 1; i <= level; i += speed {
		conn.Set("ao-volume", i)
		time.Sleep(time.Duration(300/speed) * time.Millisecond)
	}
}

func (s *MPVSource) FadeOut(conn *mpvipc.Connection, speed int) {
	volume, err := conn.Get("ao-volume")
	if err == nil {
		for i := int(volume.(float64)) - 1; i > 0; i -= speed {
			if conn.Set("ao-volume", i) == nil {
				time.Sleep(time.Duration(300/speed) * time.Millisecond)
			}
		}
	}
}

func (s *MPVSource) Disable() error {
	if s.process != nil {
		if s.process.Process != nil {
			if s.ipcSocketDir != "" {
				conn := mpvipc.NewConnection(s.ipcSocket())
				err := conn.Open()
				if err == nil {
					s.FadeOut(conn, 3)
					conn.Close()
				}
			}

			s.process.Process.Kill()
		}
	}

	return nil
}

func (s *MPVSource) CurrentlyPlaying() string {
	if s.ipcSocketDir != "" {
		conn := mpvipc.NewConnection(s.ipcSocket())
		err := conn.Open()
		if err != nil {
			log.Println("Unable to open mpv socket:", err.Error())
			return "!"
		}
		defer conn.Close()

		title, err := conn.Get("media-title")
		if err != nil {
			log.Println("Unable to retrieve title:", err.Error())
			return "!"
		}
		return title.(string)
	}

	return "-"
}

func (s *MPVSource) TogglePause(id string) error {
	if s.ipcSocketDir == "" {
		return fmt.Errorf("Not supported")
	}

	conn := mpvipc.NewConnection(s.ipcSocket())
	err := conn.Open()
	if err != nil {
		return err
	}
	defer conn.Close()

	paused, err := conn.Get("pause")
	if err != nil {
		return err
	}

	if !paused.(bool) {
		s.FadeOut(conn, 5)
	}

	err = conn.Set("pause", !paused.(bool))
	if err != nil {
		return err
	}

	if paused.(bool) {
		s.FadeIn(conn, 5, 50)
	}

	return nil
}

func (s *MPVSource) HasPlaylist() bool {
	if s.ipcSocketDir == "" {
		return false
	}

	conn := mpvipc.NewConnection(s.ipcSocket())
	err := conn.Open()
	if err != nil {
		return false
	}
	defer conn.Close()

	plistCount, err := conn.Get("playlist-count")
	if err != nil {
		return false
	}

	return plistCount.(float64) > 1
}

func (s *MPVSource) NextTrack() error {
	if s.ipcSocketDir == "" {
		return fmt.Errorf("Not supported")
	}

	conn := mpvipc.NewConnection(s.ipcSocket())
	err := conn.Open()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Call("playlist-next", "weak")
	if err != nil {
		return err
	}

	return nil
}

func (s *MPVSource) NextRandomTrack() error {
	if s.ipcSocketDir == "" {
		return fmt.Errorf("Not supported")
	}

	conn := mpvipc.NewConnection(s.ipcSocket())
	err := conn.Open()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Call("playlist-shuffle")
	if err != nil {
		return err
	}

	_, err = conn.Call("playlist-next", "weak")
	if err != nil {
		return err
	}

	_, err = conn.Call("playlist-unshuffle")
	if err != nil {
		return err
	}

	return nil
}

func (s *MPVSource) PreviousTrack() error {
	if s.ipcSocketDir == "" {
		return fmt.Errorf("Not supported")
	}

	conn := mpvipc.NewConnection(s.ipcSocket())
	err := conn.Open()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Call("playlist-prev", "weak")
	if err != nil {
		return err
	}

	return nil
}
