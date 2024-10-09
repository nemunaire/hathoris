package ui

import (
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/gin-gonic/gin"

	"git.nemunai.re/nemunaire/hathoris/config"
)

func serveOrReverse(forced_url string, cfg *config.Config) gin.HandlerFunc {
	if cfg.DevProxy != "" {
		// Forward to the Vue dev proxy
		return func(c *gin.Context) {
			if u, err := url.Parse(cfg.DevProxy); err != nil {
				http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			} else {
				if forced_url != "" {
					u.Path = path.Join(u.Path, forced_url)
				} else {
					u.Path = path.Join(u.Path, c.Request.URL.Path)
				}

				if r, err := http.NewRequest(c.Request.Method, u.String(), c.Request.Body); err != nil {
					http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
				} else if resp, err := http.DefaultClient.Do(r); err != nil {
					http.Error(c.Writer, err.Error(), http.StatusBadGateway)
				} else {
					defer resp.Body.Close()

					for key := range resp.Header {
						c.Writer.Header().Add(key, resp.Header.Get(key))
					}
					c.Writer.WriteHeader(resp.StatusCode)

					io.Copy(c.Writer, resp.Body)
				}
			}
		}
	} else if forced_url != "" {
		// Serve forced_url
		return func(c *gin.Context) {
			c.FileFromFS(forced_url, Assets)
		}
	} else {
		// Serve requested file
		return func(c *gin.Context) {
			c.FileFromFS(c.Request.URL.Path, Assets)
		}
	}
}

func DeclareRoutes(router *gin.Engine, cfg *config.Config) {
	if cfg.DevProxy != "" {
		router.GET("/.svelte-kit/*_", serveOrReverse("", cfg))
		router.GET("/node_modules/*_", serveOrReverse("", cfg))
		router.GET("/@vite/*_", serveOrReverse("", cfg))
		router.GET("/@id/*_", serveOrReverse("", cfg))
		router.GET("/@fs/*_", serveOrReverse("", cfg))
		router.GET("/src/*_", serveOrReverse("", cfg))
	}

	router.GET("/", serveOrReverse("", cfg))
	router.GET("/settings", serveOrReverse("", cfg))

	router.GET("/_app/*_", serveOrReverse("", cfg))
	router.GET("/img/*_", serveOrReverse("", cfg))
	router.GET("/favicon.ico", serveOrReverse("", cfg))
}
