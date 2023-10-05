package services

import (
	"github.com/alin-io/pkgstore/db"
	"github.com/alin-io/pkgstore/models"
	"github.com/alin-io/pkgstore/storage"
	"github.com/gin-gonic/gin"
)

type ApiService struct {
	Storage storage.BaseStorageBackend
}

func NewApiService(storageBackend storage.BaseStorageBackend) *ApiService {
	return &ApiService{
		Storage: storageBackend,
	}
}

func (s *ApiService) ListPackagesHandler(c *gin.Context) {
	pkgs := make([]models.Package[map[string]interface{}], 0)
	db.DB().Model(&pkgs).Preload("Versions").Find(&pkgs)
	c.JSON(200, pkgs)
}
