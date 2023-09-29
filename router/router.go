package router

import (
	"github.com/alin-io/pkgproxy/middlewares"
	"github.com/alin-io/pkgproxy/services"
	"github.com/alin-io/pkgproxy/storage"
	"github.com/gin-gonic/gin"
)

func SetupGinServer() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	return r
}

func PackageRouter(r *gin.Engine, storageBackend storage.BaseStorageBackend) {
	r.Use(middlewares.AuthMiddleware)

	r.GET("/", services.HealthCheckHandler)

	r.RedirectTrailingSlash = false

	initNpmRoutes(r, storageBackend)
	initPypiRoutes(r, storageBackend)
	initContainerRoutes(r, storageBackend)
}
