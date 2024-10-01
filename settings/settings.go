package settings

import (
	"encoding/json"
	"fmt"
	"os"

	"git.nemunai.re/nemunaire/hathoris/config"
)

type CustomSource struct {
	Source string            `json:"src"`
	KV     map[string]string `json:"kv"`
}

type Settings struct {
	CustomSources []CustomSource `json:"custom_sources"`
}

func Load(cfg *config.Config) (*Settings, error) {
	if cfg.SettingsPath == "" {
		return &Settings{}, nil
	}

	if st, err := os.Stat(cfg.SettingsPath); os.IsNotExist(err) || (err == nil && st.Size() == 0) {
		fd, err := os.Create(cfg.SettingsPath)
		if err != nil {
			return nil, fmt.Errorf("unable to create settings file: %w", err)
		}

		_, err = fd.Write([]byte("{}"))
		fd.Close()
		if err != nil {
			return nil, fmt.Errorf("unable to write to settings file: %w", err)
		}
	}

	fd, err := os.Open(cfg.SettingsPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open settings: %w", err)
	}
	defer fd.Close()

	var settings Settings
	err = json.NewDecoder(fd).Decode(&settings)
	if err != nil {
		return nil, fmt.Errorf("unable to read settings: %w", err)
	}

	return &settings, nil
}

func (settings *Settings) Save(path string) error {
	fd, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create settings file: %w", err)
	}
	defer fd.Close()

	err = json.NewEncoder(fd).Encode(settings)
	if err != nil {
		return fmt.Errorf("unable to read settings: %w", err)
	}

	return nil
}
