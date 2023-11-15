package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"git.nemunai.re/nemunaire/hathoris/api"
	"git.nemunai.re/nemunaire/hathoris/config"
	"git.nemunai.re/nemunaire/hathoris/sources"
	"git.nemunai.re/nemunaire/hathoris/ui"
)

type App struct {
	cfg    *config.Config
	router *gin.Engine
	srv    *http.Server
}

func NewApp(cfg *config.Config) *App {
	if cfg.DevProxy == "" {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.ForceConsoleColor()
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Next()
	})

	// Prepare struct
	app := &App{
		cfg:    cfg,
		router: router,
	}

	// Register routes
	ui.DeclareRoutes(router, cfg)
	api.DeclareRoutes(router, cfg)

	router.GET("/api/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"version": Version})
	})

	return app
}

func (app *App) Start() {
	app.srv = &http.Server{
		Addr:    app.cfg.Bind,
		Handler: app.router,
	}

	log.Printf("Ready, listening on %s\n", app.cfg.Bind)
	if err := app.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func (app *App) Stop() {
	// Disable all sources
	someEnabled := false
	for k, src := range sources.SoundSources {
		if src.IsEnabled() {
			someEnabled = true
			go func(k string, src sources.SoundSource) {
				log.Printf("Stopping %s...", k)
				err := src.Disable()
				if err != nil {
					log.Printf("Unable to disable %s source", k)
				}
				log.Printf("%s stopped", k)
			}(k, src)
		}
	}

	// Wait for fadeout
	if someEnabled {
		time.Sleep(2000 * time.Millisecond)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
