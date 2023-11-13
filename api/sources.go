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
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Active      *bool  `json:"active,omitempty"`
	Controlable bool   `json:"controlable,omitempty"`
}

func declareSourcesRoutes(cfg *config.Config, router *gin.RouterGroup) {
	router.GET("/sources", func(c *gin.Context) {
		ret := map[string]*SourceState{}

		for k, src := range sources.SoundSources {
			active := src.IsActive()
			_, controlable := src.(inputs.ControlableInput)

			ret[k] = &SourceState{
				Name:        src.GetName(),
				Enabled:     src.IsEnabled(),
				Active:      &active,
				Controlable: controlable,
			}
		}

		c.JSON(http.StatusOK, ret)
	})

	sourcesRoutes := router.Group("/sources/:source")
	sourcesRoutes.Use(SourceHandler)

	sourcesRoutes.GET("", func(c *gin.Context) {
		src := c.MustGet("source").(sources.SoundSource)

		active := src.IsActive()

		c.JSON(http.StatusOK, &SourceState{
			Name:    src.GetName(),
			Enabled: src.IsEnabled(),
			Active:  &active,
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
