package router

import (
	"github.com/alin-io/pkgstore/services/pypi"
	"github.com/alin-io/pkgstore/storage"
	"github.com/gin-gonic/gin"
)

func initPypiRoutes(r *gin.Engine, storageBackend storage.BaseStorageBackend) {
	pypiService := pypi.NewService(storageBackend)
	pypiRoutes := r.Group("/pypi")
	{
		pypiRoutes.GET("/simple/:name", pypiService.MetadataHandler)
		pypiRoutes.GET("/simple/:name/:name2", pypiService.MetadataHandler)

		pypiRoutes.GET("/files/:sha256/:filename", pypiService.DownloadHandler)

		pypiRoutes.POST("", pypiService.UploadHandler)
	}
}
