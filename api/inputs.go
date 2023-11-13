package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.nemunai.re/nemunaire/hathoris/config"
	"git.nemunai.re/nemunaire/hathoris/inputs"
)

type InputState struct {
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	Controlable bool   `json:"controlable"`
}

func declareInputsRoutes(cfg *config.Config, router *gin.RouterGroup) {
	router.GET("/inputs", func(c *gin.Context) {
		ret := map[string]*InputState{}

		for k, inp := range inputs.SoundInputs {
			_, controlable := inp.(inputs.ControlableInput)

			ret[k] = &InputState{
				Name:        inp.GetName(),
				Active:      inp.IsActive(),
				Controlable: controlable,
			}
		}

		c.JSON(http.StatusOK, ret)
	})

	inputsRoutes := router.Group("/inputs/:input")
	inputsRoutes.Use(InputHandler)

	inputsRoutes.GET("", func(c *gin.Context) {
		src := c.MustGet("input").(inputs.SoundInput)

		c.JSON(http.StatusOK, &InputState{
			Name:   src.GetName(),
			Active: src.IsActive(),
		})
	})
	inputsRoutes.GET("/settings", func(c *gin.Context) {
		c.JSON(http.StatusOK, c.MustGet("input"))
	})
	inputsRoutes.GET("/currently", func(c *gin.Context) {
		src := c.MustGet("input").(inputs.SoundInput)

		if !src.IsActive() {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"errmsg": "Input not active"})
			return
		}

		c.JSON(http.StatusOK, src.CurrentlyPlaying())
	})

	streamRoutes := inputsRoutes.Group("/stream/:stream")
	streamRoutes.Use(StreamHandler)

	streamRoutes.POST("/pause", func(c *gin.Context) {
		input, ok := c.MustGet("input").(inputs.ControlableInput)
		if !ok {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"errmsg": "The source doesn't support that"})
			return
		}

		err := input.TogglePause(c.MustGet("streamid").(string))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errmsg": fmt.Sprintf("Unable to pause the input: %s", err.Error())})
			return
		}

		c.JSON(http.StatusOK, true)
	})
}

func InputHandler(c *gin.Context) {
	src, ok := inputs.SoundInputs[c.Param("input")]
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"errmsg": fmt.Sprintf("Input not found: %s", c.Param("input"))})
		return
	}

	c.Set("input", src)

	c.Next()
}

func StreamHandler(c *gin.Context) {
	input := c.MustGet("input").(inputs.SoundInput)

	streams := input.CurrentlyPlaying()

	stream, ok := streams[c.Param("stream")]
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"errmsg": fmt.Sprintf("Stream not found: %s", c.Param("stream"))})
		return
	}

	c.Set("streamid", c.Param("stream"))
	c.Set("stream", stream)

	c.Next()
}
