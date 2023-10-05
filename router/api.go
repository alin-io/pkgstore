package router

import (
	"github.com/alin-io/pkgstore/services"
	"github.com/alin-io/pkgstore/storage"
	"github.com/gin-gonic/gin"
)

func initApiRoutes(r *gin.Engine, storageBackend storage.BaseStorageBackend) {
	apiService := services.NewApiService(storageBackend)
	apiRoutes := r.Group("/api")
	{
		apiRoutes.GET("/packages", apiService.ListPackagesHandler)
	}
}
