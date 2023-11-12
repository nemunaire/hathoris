package config

import (
	"flag"
	"strings"
)

type Config struct {
	DevProxy string
	Bind     string
	BaseURL  string
}

// parseLine treats a config line and place the read value in the variable
// declared to the corresponding flag.
func (c *Config) parseLine(line string) (err error) {
	fields := strings.SplitN(line, "=", 2)
	orig_key := strings.TrimSpace(fields[0])
	value := strings.TrimSpace(fields[1])

	key := strings.TrimPrefix(orig_key, "HATHORIS_")
	key = strings.Replace(key, "_", "-", -1)
	key = strings.ToLower(key)

	err = flag.Set(key, value)

	return
}
