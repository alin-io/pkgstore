package router

import (
	"github.com/alin-io/pkgstore/services/api"
	"github.com/alin-io/pkgstore/storage"
	"github.com/gin-gonic/gin"
)

func initApiRoutes(r *gin.Engine, storageBackend storage.BaseStorageBackend) {
	apiService := api.NewApiService(storageBackend)
	apiRoutes := r.Group("/api")
	{
		apiRoutes.GET("/packages", apiService.ListPackagesHandler)
		apiRoutes.GET("/packages/:id", apiService.GetPackage)
		apiRoutes.GET("/packages/:id/versions", apiService.ListVersionsHandler)

		apiRoutes.DELETE("/packages/:id", apiService.DeletePackage)
		apiRoutes.DELETE("/packages/:id/versions/:versionId", apiService.DeleteVersion)
	}
}
