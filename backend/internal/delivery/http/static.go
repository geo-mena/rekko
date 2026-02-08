package http

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func registerStaticRoutes(router *gin.Engine, staticDir string) {
	fs := http.FileServer(http.Dir(staticDir))

	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/swagger/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		filePath := filepath.Join(staticDir, path)
		if _, err := os.Stat(filePath); err == nil {
			fs.ServeHTTP(c.Writer, c.Request)
			return
		}

		c.File(filepath.Join(staticDir, "index.html"))
	})
}
