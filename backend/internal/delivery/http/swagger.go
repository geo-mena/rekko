package http

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"

	swaggerFiles "github.com/swaggo/files/v2"
)

func swaggerHandler() gin.HandlerFunc {
	staticFS := http.FS(swaggerFiles.FS)
	fileServer := http.StripPrefix("/swagger", http.FileServer(staticFS))

	return func(c *gin.Context) {
		switch c.Param("any") {
		case "/", "":
			f, err := swaggerFiles.FS.(fs.ReadFileFS).ReadFile("index.html")
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.Data(http.StatusOK, "text/html; charset=utf-8", f)
		case "/doc.json":
			doc, err := swag.ReadDoc()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			c.Header("Content-Type", "application/json")
			c.String(http.StatusOK, doc)
		case "/swagger-initializer.js":
			c.Header("Content-Type", "application/javascript")
			c.String(http.StatusOK, swaggerInitJS)
		default:
			fileServer.ServeHTTP(c.Writer, c.Request)
		}
	}
}

const swaggerInitJS = `window.onload = function() {
  window.ui = SwaggerUIBundle({
    url: "/swagger/doc.json",
    dom_id: '#swagger-ui',
    presets: [
      SwaggerUIBundle.presets.apis
    ],
    layout: "BaseLayout"
  });
};`
