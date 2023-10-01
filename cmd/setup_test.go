package cmd

import (
	"github.com/alin-io/pkgstore/db"
	"github.com/alin-io/pkgstore/models"
	"github.com/alin-io/pkgstore/router"
	"github.com/alin-io/pkgstore/services/npm"
	"github.com/alin-io/pkgstore/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	storageBackend storage.BaseStorageBackend
	serverApp      *gin.Engine
)

func init() {
	storageBackend = storage.NewInMemoryBackend()
	db.InitDatabaseForTest()
	models.SyncModels()

	serverApp = router.SetupGinServer()
	serverApp.Use(func(c *gin.Context) {
		c.Set("testing", true)
	})
	router.PackageRouter(serverApp, storageBackend)
}

func TestServerHealthCheck(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	serverApp.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func DeleteTestPackage(name, service string) error {
	return db.DB().Delete(&models.Package[npm.MetadataResponse]{}, "name = ? AND service = ?", name, service).Error
}
