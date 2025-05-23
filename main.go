package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"git.nemunai.re/nemunaire/hathoris/config"
	_ "git.nemunai.re/nemunaire/hathoris/inputs/mpris"
	_ "git.nemunai.re/nemunaire/hathoris/inputs/pulseaudio"
	_ "git.nemunai.re/nemunaire/hathoris/sources/amp1_gpio"
	_ "git.nemunai.re/nemunaire/hathoris/sources/mpv"
	_ "git.nemunai.re/nemunaire/hathoris/sources/spdif"
)

var (
	Version = "custom-build"
)

func main() {
	cfg, err := config.Consolidated()
	if err != nil {
		log.Fatal("Unable to read configuration:", err)
	}

	// Start app
	a := NewApp(cfg)
	go a.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Stopping the service...")
	a.Stop()
	log.Println("Stopped")
}
