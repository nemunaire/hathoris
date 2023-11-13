package mpv

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/DexterLB/mpvipc"

	"git.nemunai.re/nemunaire/hathoris/sources"
)

type MPVSource struct {
	process   *exec.Cmd
	ipcSocket string
	Name      string
	Options   []string
	File      string
}

func init() {
	sources.SoundSources["mpv-1"] = &MPVSource{
		Name:      "Radio 1",
		ipcSocket: "/tmp/tmpmpv.radio-1",
		Options:   []string{"--no-video", "--no-terminal"},
		File:      "https://mediaserv38.live-streams.nl:18030/stream",
	}
	sources.SoundSources["mpv-2"] = &MPVSource{
		Name:      "Radio 2",
		ipcSocket: "/tmp/tmpmpv.radio-2",
		Options:   []string{"--no-video", "--no-terminal"},
		File:      "https://mediaserv38.live-streams.nl:18040/live",
	}
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

func (s *MPVSource) Enable() (err error) {
	if s.process != nil {
		return fmt.Errorf("Already running")
	}

	var opts []string
	opts = append(opts, s.Options...)
	if s.ipcSocket != "" {
		opts = append(opts, "--input-ipc-server="+s.ipcSocket, "--pause")
	}
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

	if s.ipcSocket != "" {
		_, err = os.Stat(s.ipcSocket)
		for i := 20; i >= 0 && err != nil; i-- {
			time.Sleep(100 * time.Millisecond)
			_, err = os.Stat(s.ipcSocket)
		}
		time.Sleep(200 * time.Millisecond)

		conn := mpvipc.NewConnection(s.ipcSocket)
		err = conn.Open()
		for i := 20; i >= 0 && err != nil; i-- {
			time.Sleep(100 * time.Millisecond)
			err = conn.Open()
		}
		if err != nil {
			return err
		}
		defer conn.Close()

		_, err = conn.Get("media-title")
		for err != nil {
			time.Sleep(100 * time.Millisecond)
			_, err = conn.Get("media-title")
		}

		conn.Set("ao-volume", 50)

		err = conn.Set("pause", false)
		if err != nil {
			return err
		}
	}

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

func (s *MPVSource) CurrentlyPlaying() string {
	if s.ipcSocket != "" {
		conn := mpvipc.NewConnection(s.ipcSocket)
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
