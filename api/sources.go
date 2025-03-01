package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.nemunai.re/nemunaire/hathoris/config"
	"git.nemunai.re/nemunaire/hathoris/inputs"
	"git.nemunai.re/nemunaire/hathoris/sources"
)

type SourceState struct {
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
	Active       *bool  `json:"active,omitempty"`
	Controlable  bool   `json:"controlable,omitempty"`
	HasPlaylist  bool   `json:"hasplaylist,omitempty"`
	CurrentTitle string `json:"currentTitle,omitempty"`
}

func declareSourcesRoutes(cfg *config.Config, router *gin.RouterGroup) {
	router.GET("/sources", func(c *gin.Context) {
		ret := map[string]*SourceState{}

		for k, src := range sources.SoundSources {
			active := src.IsActive()
			_, controlable := src.(inputs.ControlableInput)

			var hasPlaylist bool
			if p, withPlaylist := src.(inputs.PlaylistInput); withPlaylist {
				hasPlaylist = p.HasPlaylist()
			}

			var title string
			if s, ok := src.(sources.PlayingSource); ok && active {
				title = s.CurrentlyPlaying()
			}

			ret[k] = &SourceState{
				Name:         src.GetName(),
				Enabled:      src.IsEnabled(),
				Active:       &active,
				Controlable:  controlable,
				HasPlaylist:  hasPlaylist,
				CurrentTitle: title,
			}
		}

		c.JSON(http.StatusOK, ret)
	})

	sourcesRoutes := router.Group("/sources/:source")
	sourcesRoutes.Use(SourceHandler)

	sourcesRoutes.GET("", func(c *gin.Context) {
		src := c.MustGet("source").(sources.SoundSource)

		active := src.IsActive()
		_, controlable := src.(inputs.ControlableInput)

		var hasPlaylist bool
		if p, withPlaylist := src.(inputs.PlaylistInput); withPlaylist {
			hasPlaylist = p.HasPlaylist()
		}

		var title string
		if s, ok := src.(sources.PlayingSource); ok && active {
			title = s.CurrentlyPlaying()
		}

		c.JSON(http.StatusOK, &SourceState{
			Name:         src.GetName(),
			Enabled:      src.IsEnabled(),
			Active:       &active,
			Controlable:  controlable,
			HasPlaylist:  hasPlaylist,
			CurrentTitle: title,
		})
	})
	sourcesRoutes.GET("/settings", func(c *gin.Context) {
		c.JSON(http.StatusOK, c.MustGet("source"))
	})
	sourcesRoutes.GET("/currently", func(c *gin.Context) {
		src := c.MustGet("source").(sources.SoundSource)

		if !src.IsActive() {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"errmsg": "Source not active"})
			return
		}

		s, ok := src.(sources.PlayingSource)
		if !ok {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"errmsg": "The source doesn't support"})
			return
		}

		c.JSON(http.StatusOK, s.CurrentlyPlaying())
	})
	sourcesRoutes.POST("/enable", func(c *gin.Context) {
		src := c.MustGet("source").(sources.SoundSource)

		if src.IsEnabled() {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errmsg": "The source is already enabled"})
			return
		}

		// Disable all sources
		for k, src := range sources.SoundSources {
			if src.IsEnabled() {
				err := src.Disable()
				if err != nil {
					log.Printf("Unable to disable %s: %s", k, err.Error())
				}
			}
		}

		err := src.Enable()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errmsg": fmt.Sprintf("Unable to enable the source: %s", err.Error())})
			return
		}

		c.JSON(http.StatusOK, true)
	})
	sourcesRoutes.POST("/disable", func(c *gin.Context) {
		src := c.MustGet("source").(sources.SoundSource)

		err := src.Disable()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errmsg": fmt.Sprintf("Unable to disable the source: %s", err.Error())})
			return
		}

		c.JSON(http.StatusOK, true)
	})

	// ControlableInput
	sourcesRoutes.POST("/pause", func(c *gin.Context) {
		src := c.MustGet("source").(sources.SoundSource)

		if !src.IsActive() {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"errmsg": "Source not active"})
			return
		}

		s, ok := src.(inputs.ControlableInput)
		if !ok {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"errmsg": "The source doesn't support"})
			return
		}

		c.JSON(http.StatusOK, s.TogglePause("default"))
	})

	// PlaylistInput
	sourcesRoutes.POST("/has_playlist", func(c *gin.Context) {
		src := c.MustGet("source").(sources.SoundSource)

		if !src.IsActive() {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"errmsg": "Source not active"})
			return
		}

		s, ok := src.(inputs.PlaylistInput)
		if !ok {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"errmsg": "The source doesn't support"})
			return
		}

		c.JSON(http.StatusOK, s.HasPlaylist())
	})
	sourcesRoutes.POST("/next_track", func(c *gin.Context) {
		src := c.MustGet("source").(sources.SoundSource)

		if !src.IsActive() {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"errmsg": "Source not active"})
			return
		}

		s, ok := src.(inputs.PlaylistInput)
		if !ok || !s.HasPlaylist() {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"errmsg": "The source doesn't support"})
			return
		}

		err := s.NextTrack()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"errmsg": err.Error()})
			return
		}

		c.JSON(http.StatusOK, true)
	})
	sourcesRoutes.POST("/next_random_track", func(c *gin.Context) {
		src := c.MustGet("source").(sources.SoundSource)

		if !src.IsActive() {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"errmsg": "Source not active"})
			return
		}

		s, ok := src.(inputs.PlaylistInput)
		if !ok || !s.HasPlaylist() {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"errmsg": "The source doesn't support"})
			return
		}

		err := s.NextRandomTrack()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"errmsg": err.Error()})
			return
		}

		c.JSON(http.StatusOK, true)
	})
	sourcesRoutes.POST("/prev_track", func(c *gin.Context) {
		src := c.MustGet("source").(sources.SoundSource)

		if !src.IsActive() {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"errmsg": "Source not active"})
			return
		}

		s, ok := src.(inputs.PlaylistInput)
		if !ok || !s.HasPlaylist() {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"errmsg": "The source doesn't support"})
			return
		}

		err := s.PreviousTrack()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"errmsg": err.Error()})
			return
		}

		c.JSON(http.StatusOK, true)
	})
}

func SourceHandler(c *gin.Context) {
	src, ok := sources.SoundSources[c.Param("source")]
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"errmsg": fmt.Sprintf("Source not found: %s", c.Param("source"))})
		return
	}

	c.Set("source", src)

	c.Next()
}
