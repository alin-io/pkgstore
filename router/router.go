package router

import (
	"github.com/alin-io/pkgstore/middlewares"
	"github.com/alin-io/pkgstore/services"
	"github.com/alin-io/pkgstore/storage"
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
	initApiRoutes(r, storageBackend)
}
