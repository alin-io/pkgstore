package router

import (
	"github.com/alin-io/pkgstore/services/container"
	"github.com/alin-io/pkgstore/storage"
	"github.com/gin-gonic/gin"
)

func initContainerRoutes(r *gin.Engine, storageBackend storage.BaseStorageBackend) {
	containerService := container.NewService(storageBackend)
	containerRoutes := r.Group("/v2")
	{
		containerRoutes.GET("/", func(context *gin.Context) {
			context.JSON(200, gin.H{"status": "ok"})
		})
		// Upload Process
		containerRoutes.GET(":name/blobs/uploads/:uuid", containerService.GetUploadProgressHandler)
		containerRoutes.HEAD(":name/blobs/:sha256", containerService.CheckBlobExistenceHandler)
		containerRoutes.POST(":name/blobs/uploads/", containerService.StartLayerUploadHandler)
		containerRoutes.PATCH(":name/blobs/uploads/:uuid", containerService.ChunkUploadHandler)
		containerRoutes.PUT(":name/blobs/uploads/:uuid", containerService.UploadHandler)
		containerRoutes.PUT(":name/manifests/:reference", containerService.ManifestUploadHandler)

		// Download Process
		containerRoutes.GET(":name/manifests/:reference", containerService.MetadataHandler)
		containerRoutes.HEAD(":name/manifests/:reference", containerService.CheckMetadataHandler)
		containerRoutes.GET(":name/blobs/:sha256", containerService.DownloadHandler)
	}
}
