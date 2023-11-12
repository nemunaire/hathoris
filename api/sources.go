package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.nemunai.re/nemunaire/hathoris/config"
	"git.nemunai.re/nemunaire/hathoris/sources"
)

type SourceState struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Active  *bool  `json:"active,omitempty"`
}

func declareSourcesRoutes(cfg *config.Config, router *gin.RouterGroup) {
	router.GET("/sources", func(c *gin.Context) {
		ret := map[string]*SourceState{}

		for k, src := range sources.SoundSources {
			ret[k] = &SourceState{
				Name:    src.GetName(),
				Enabled: src.IsEnabled(),
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
	sourcesRoutes.POST("/enable", func(c *gin.Context) {
		src := c.MustGet("source").(sources.SoundSource)

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
