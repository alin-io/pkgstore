package router

import (
	"github.com/alin-io/pkgstore/services/npm"
	"github.com/alin-io/pkgstore/storage"
	"github.com/gin-gonic/gin"
)

func initNpmRoutes(r *gin.Engine, storageBackend storage.BaseStorageBackend) {
	npmService := npm.NewService(storageBackend)
	npmRoutes := r.Group("/npm/:name")
	{
		npmRoutes.GET("", npmService.MetadataHandler)
		npmRoutes.GET("/-/:filename", npmService.DownloadHandler)
		npmRoutes.GET(":name2", npmService.MetadataHandler)
		npmRoutes.GET(":name2/-/:filename", npmService.DownloadHandler)

		npmRoutes.PUT("", npmService.UploadHandler)
		npmRoutes.PUT(":name2", npmService.UploadHandler)
	}
}
