package router

import (
	"github.com/alin-io/pkgproxy/services/container"
	"github.com/alin-io/pkgproxy/storage"
	"github.com/gin-gonic/gin"
)

func initContainerRoutes(r *gin.Engine, storageBackend storage.BaseStorageBackend) {
	containerService := container.NewService(storageBackend)
	containerRoutes := r.Group("/container/v2")
	{
		// Upload Process
		containerRoutes.GET(":name/blobs/uploads/:uuid", containerService.GetUploadProgressHandler)
		containerRoutes.POST(":name/blobs/uploads", containerService.StartLayerUploadHandler)
		containerRoutes.PATCH(":name/blobs/uploads/:uuid", containerService.ChunkUploadHandler)
		containerRoutes.PUT(":name/blobs/uploads/:uuid", containerService.UploadHandler)
		containerRoutes.PUT(":name/manifests/:reference", containerService.ManifestUploadHandler)
	}
}
