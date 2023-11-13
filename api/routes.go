package api

import (
	"github.com/gin-gonic/gin"

	"git.nemunai.re/nemunaire/hathoris/config"
)

func DeclareRoutes(router *gin.Engine, cfg *config.Config) {
	apiRoutes := router.Group("/api")

	declareInputsRoutes(cfg, apiRoutes)
	declareSourcesRoutes(cfg, apiRoutes)
	declareVolumeRoutes(cfg, apiRoutes)
}
