package main

import (
	"embed"
	"github.com/alin-io/pkgstore/config"
	"github.com/alin-io/pkgstore/db"
	_ "github.com/alin-io/pkgstore/db"
	"github.com/alin-io/pkgstore/models"
	"github.com/alin-io/pkgstore/router"
	"github.com/alin-io/pkgstore/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

//go:embed all:ui
var frontendFS embed.FS

func main() {
	var storageBackend storage.BaseStorageBackend
	if config.Get().Storage.ActiveBackend == config.StorageS3 {
		storageBackend = storage.NewS3Backend()
	} else {
		panic("Unknown storage backend")
	}

	// Initialize the DB connection
	db.InitDatabase()

	// Sync Models with the DB
	models.SyncModels()

	r := router.SetupGinServer()
	// Setup Cors if we are in Debug mode, otherwise UI would be under the same domain name
	if gin.Mode() == gin.DebugMode {
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = []string{"http://localhost:3000"}
		corsConfig.AllowCredentials = true
		corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
		r.Use(cors.New(corsConfig))
	}

	templates := template.Must(template.New("").ParseFS(frontendFS, "ui/index.html"))
	r.SetHTMLTemplate(templates)
	r.GET("/ui", serveIndexTemplate)
	r.GET("/ui/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		f, err := frontendFS.Open("ui" + filepath)
		if f != nil {
			f.Close()
		}

		if err != nil || filepath == "" || filepath == "/" || filepath == "index.html" {
			serveIndexTemplate(c)
		} else {
			http.FileServer(http.FS(frontendFS)).ServeHTTP(c.Writer, c.Request)
		}
	})

	router.PackageRouter(r, storageBackend)

	err := r.Run(config.Get().ListenAddress)
	if err != nil {
		panic(err)
	}
}

func serveIndexTemplate(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})

}
