package config

import (
	"fmt"
	"os"
	"strings"
)

// FromEnv analyzes all the environment variables to find each one
// starting by HATHORIS_
func (c *Config) FromEnv() error {
	for _, line := range os.Environ() {
		if strings.HasPrefix(line, "HATHORIS_") {
			err := c.parseLine(line)
			if err != nil {
				return fmt.Errorf("error in environment (%q): %w", line, err)
			}
		}
	}
	return nil
}
