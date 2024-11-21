package api

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.nemunai.re/nemunaire/hathoris/alsacontrol"
	"git.nemunai.re/nemunaire/hathoris/config"
)

var cardId string = "0"

func init() {
	flag.StringVar(&cardId, "card-id", cardId, "ALSA card identifier for volume handling")
}

func declareVolumeRoutes(cfg *config.Config, router *gin.RouterGroup) {
	router.GET("/mixer", func(c *gin.Context) {
		cnt, err := alsa.ParseAmixerContent(cardId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errmsg": err.Error()})
			return
		}

		var ret []*alsa.CardControlState

		for _, cc := range cnt {
			ret = append(ret, cc.ToCardControlState())
		}

		c.JSON(http.StatusOK, ret)
	})

	mixerRoutes := router.Group("/mixer/:mixer")
	mixerRoutes.Use(MixerHandler)

	mixerRoutes.GET("", func(c *gin.Context) {
		cc := c.MustGet("mixer").(*alsa.CardControl)

		c.JSON(http.StatusOK, cc.ToCardControlState())
	})
	mixerRoutes.POST("/values", func(c *gin.Context) {
		cc := c.MustGet("mixer").(*alsa.CardControl)

		var valuesINT []interface{}

		err := c.BindJSON(&valuesINT)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errmsg": fmt.Sprintf("Unable to parse values: %s", err.Error())})
			return
		}

		var values []string
		for _, v := range valuesINT {
			if t, ok := v.(float64); ok {
				if float64(int64(t)) == t {
					values = append(values, strconv.FormatInt(int64(t), 10))
				} else {
					values = append(values, fmt.Sprintf("%f", t))
				}
			} else if t, ok := v.(bool); ok {
				if t {
					values = append(values, "on")
				} else {
					values = append(values, "off")
				}
			}
		}

		err = cc.CsetAmixer(cardId, values...)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errmsg": fmt.Sprintf("Unable to set values: %s", err.Error())})
			return
		}

		c.JSON(http.StatusOK, values)
	})
}

func MixerHandler(c *gin.Context) {
	mixers, err := alsa.ParseAmixerContent(cardId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errmsg": err.Error()})
		return
	}

	mixer := c.Param("mixer")

	for _, m := range mixers {
		if strconv.FormatInt(m.NumID, 10) == mixer || m.Name == mixer {
			c.Set("mixer", m)
			c.Next()
			return
		}
	}

	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"errmsg": "Mixer not found"})
}
