package api

import (
	"github.com/gin-gonic/gin"

	"git.nemunai.re/nemunaire/hathoris/config"
	"git.nemunai.re/nemunaire/hathoris/settings"
)

type SettingsGetter func() *settings.Settings

type SettingsReloader func() error

func DeclareRoutes(router *gin.Engine, cfg *config.Config, getsettings SettingsGetter, reloadsettings SettingsReloader) {
	apiRoutes := router.Group("/api")

	declareInputsRoutes(cfg, apiRoutes)
	declareSettingsRoutes(cfg, getsettings, reloadsettings, apiRoutes)
	declareSourcesRoutes(cfg, apiRoutes)
	declareVolumeRoutes(cfg, apiRoutes)
}
