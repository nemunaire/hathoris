package mpris

import (
	"fmt"
	"log"
	"strings"

	"git.nemunai.re/nemunaire/hathoris/inputs"

	"github.com/godbus/dbus/v5"
	"github.com/leberKleber/go-mpris"
)

type MPRISClient struct {
	Id   string
	Name string
	Path string
}

var KNOWN_CLIENTS = []MPRISClient{
	MPRISClient{"shairport", "ShairportSync", "org.mpris.MediaPlayer2.ShairportSync."},
	MPRISClient{"firefox", "Firefox", "org.mpris.MediaPlayer2.firefox."},
}

type MPRISInput struct {
	player *mpris.Player
	Name   string
	Path   string
}

var dbusConn *dbus.Conn

func init() {
	var err error
	dbusConn, err = dbus.ConnectSessionBus()
	if err != nil {
		dbusConn, err = dbus.ConnectSystemBus()
		if err != nil {
			log.Println("Unable to connect to DBus. MPRIS will be unavailable:", err.Error())
			return
		}
	}

	var s []string
	err = dbusConn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&s)
	if err != nil {
		log.Println("DBus unavailable:", err.Error())
		return
	}
	log.Println("Available DBus entries:", strings.Join(s, ","))

	for _, ss := range s {
		for _, c := range KNOWN_CLIENTS {
			if strings.HasPrefix(ss, c.Path) {
				inputs.SoundInputs[c.Id] = &MPRISInput{
					Name: c.Name,
					Path: c.Path,
				}
			}
		}
	}
}

func (i *MPRISInput) getPlayer() (*mpris.Player, error) {
	if i.player == nil {
		var s []string
		err := dbusConn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&s)
		if err != nil {
			return nil, err
		}

		for _, ss := range s {
			if strings.HasPrefix(ss, i.Path) {
				player := mpris.NewPlayerWithConnection(ss, dbusConn)
				if err != nil {
					return nil, err
				}
				i.player = &player
				break
			}
		}

		if i.player == nil {
			return nil, fmt.Errorf("Unable to find such dBus entry")
		}
	}

	return i.player, nil
}

func (i *MPRISInput) GetName() string {
	return i.Name
}

func (i *MPRISInput) IsActive() bool {
	p, err := i.getPlayer()
	if err != nil || p == nil {
		log.Println(err)
		return false
	}

	_, err = p.Metadata()
	return err == nil
}

func (i *MPRISInput) CurrentlyPlaying() map[string]string {
	p, err := i.getPlayer()
	if err != nil || p == nil {
		log.Println(err)
		return nil
	}

	meta, err := p.Metadata()
	if err != nil {
		log.Println(err)
		return nil
	}

	var infos []string
	if artists, err := meta.XESAMArtist(); err == nil {
		for _, artist := range artists {
			if artist != "" {
				infos = append(infos, artist)
			}
		}
	}
	if title, err := meta.XESAMTitle(); err == nil && title != "" {
		infos = append(infos, title)
	}

	ret := strings.Join(infos, " - ")
	return map[string]string{
		"default": ret,
	}
}

func (i *MPRISInput) TogglePause(id string) error {
	p, err := i.getPlayer()
	if err != nil {
		return err
	}

	if ok, err := p.CanPause(); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("The player doesn't support pause")
	}

	p.Pause()
	return nil
}
