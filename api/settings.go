package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.nemunai.re/nemunaire/hathoris/config"
	"git.nemunai.re/nemunaire/hathoris/settings"
	"git.nemunai.re/nemunaire/hathoris/sources"
)

type loadableSourceExposed struct {
	Description string                 `json:"description"`
	Fields      []*sources.SourceField `json:"fields"`
}

func declareSettingsRoutes(cfg *config.Config, getsettings SettingsGetter, reloadsettings SettingsReloader, router *gin.RouterGroup) {
	router.GET("/settings", func(c *gin.Context) {
		c.JSON(http.StatusOK, getsettings())
	})

	router.POST("/settings", func(c *gin.Context) {
		var params settings.Settings

		// Parse settings
		err := c.ShouldBindJSON(&params)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errmsg": fmt.Sprintf("Something is wrong in received settings: %s", err.Error())})
			return
		}

		// Reload current settings
		*getsettings() = params

		// Save settings
		getsettings().Save(cfg.SettingsPath)

		err = reloadsettings()
		if err != nil {
			log.Println("Unable to reload settings:", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errmsg": fmt.Sprintf("Unable to reload settings: %s", err.Error())})
			return
		}

		c.JSON(http.StatusOK, getsettings())
	})

	router.GET("/settings/loadable_sources", func(c *gin.Context) {
		ret := map[string]loadableSourceExposed{}

		for k, ls := range sources.LoadableSources {
			ret[k] = loadableSourceExposed{
				Description: ls.Description,
				Fields:      sources.GenFields(ls.SourceDefinition),
			}
		}

		c.JSON(http.StatusOK, ret)
	})

}
