package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.nemunai.re/nemunaire/hathoris/config"
	"git.nemunai.re/nemunaire/hathoris/inputs"
)

type InputState struct {
	Name        string                        `json:"name"`
	Active      bool                          `json:"active"`
	Controlable bool                          `json:"controlable"`
	Streams     map[string]string             `json:"streams,omitempty"`
	Mixable     bool                          `json:"mixable"`
	Mixer       map[string]*inputs.InputMixer `json:"mixer,omitempty"`
}

func declareInputsRoutes(cfg *config.Config, router *gin.RouterGroup) {
	router.GET("/inputs", func(c *gin.Context) {
		ret := map[string]*InputState{}

		for k, inp := range inputs.SoundInputs {
			var mixer map[string]*inputs.InputMixer

			_, controlable := inp.(inputs.ControlableInput)
			im, mixable := inp.(inputs.MixableInput)
			if mixable {
				mixer, _ = im.GetMixers()
			}

			ret[k] = &InputState{
				Name:        inp.GetName(),
				Active:      inp.IsActive(),
				Controlable: controlable,
				Streams:     inp.CurrentlyPlaying(),
				Mixable:     mixable,
				Mixer:       mixer,
			}
		}

		c.JSON(http.StatusOK, ret)
	})

	inputsRoutes := router.Group("/inputs/:input")
	inputsRoutes.Use(InputHandler)

	inputsRoutes.GET("", func(c *gin.Context) {
		inp := c.MustGet("input").(inputs.SoundInput)

		var mixer map[string]*inputs.InputMixer
		_, controlable := inp.(inputs.ControlableInput)
		im, mixable := inp.(inputs.MixableInput)
		if mixable {
			mixer, _ = im.GetMixers()
		}

		c.JSON(http.StatusOK, &InputState{
			Name:        inp.GetName(),
			Active:      inp.IsActive(),
			Controlable: controlable,
			Streams:     inp.CurrentlyPlaying(),
			Mixable:     mixable,
			Mixer:       mixer,
		})
	})
	inputsRoutes.GET("/settings", func(c *gin.Context) {
		c.JSON(http.StatusOK, c.MustGet("input"))
	})
	inputsRoutes.GET("/streams", func(c *gin.Context) {
		inp := c.MustGet("input").(inputs.SoundInput)

		if !inp.IsActive() {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"errmsg": "Input not active"})
			return
		}

		c.JSON(http.StatusOK, inp.CurrentlyPlaying())
	})

	streamRoutes := inputsRoutes.Group("/streams/:stream")
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
	streamRoutes.POST("/volume", func(c *gin.Context) {
		input, ok := c.MustGet("input").(inputs.MixableInput)
		if !ok {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"errmsg": "The source doesn't support that"})
			return
		}

		var mixer inputs.InputMixer
		err := c.BindJSON(&mixer)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errmsg": err.Error()})
			return
		}

		err = input.SetMixer(c.MustGet("streamid").(string), &mixer)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errmsg": fmt.Sprintf("Unable to pause the input: %s", err.Error())})
			return
		}

		c.JSON(http.StatusOK, true)
	})
}

func InputHandler(c *gin.Context) {
	inp, ok := inputs.SoundInputs[c.Param("input")]
	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"errmsg": fmt.Sprintf("Input not found: %s", c.Param("input"))})
		return
	}

	c.Set("input", inp)

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
