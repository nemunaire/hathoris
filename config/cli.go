package config

import (
	"flag"
)

// declareFlags registers flags for the structure Options.
func (c *Config) declareFlags() {
	flag.StringVar(&c.BaseURL, "baseurl", c.BaseURL, "URL prepended to each URL")
	flag.StringVar(&c.Bind, "bind", c.Bind, "Bind port/socket")
	flag.StringVar(&c.DevProxy, "dev", c.DevProxy, "Use ui directory instead of embedded assets")

	// Others flags are declared in some other files when they need specials configurations
}

func Consolidated() (cfg *Config, err error) {
	// Define defaults options
	cfg = &Config{
		Bind: "127.0.0.1:8080",
	}

	cfg.declareFlags()

	// Then, overwrite that by what is present in the environment
	err = cfg.FromEnv()
	if err != nil {
		return
	}

	// Finaly, command line takes precedence
	err = cfg.parseCLI()
	if err != nil {
		return
	}

	return
}

// parseCLI parse the flags and treats extra args as configuration filename.
func (c *Config) parseCLI() error {
	flag.Parse()

	return nil
}
