package http

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func registerStaticRoutes(router *gin.Engine, staticDir string) {
	router.Static("/assets", filepath.Join(staticDir, "assets"))
	router.StaticFile("/favicon.ico", filepath.Join(staticDir, "favicon.ico"))

	router.NoRoute(func(c *gin.Context) {
		indexPath := filepath.Join(staticDir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			c.File(indexPath)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})
}
