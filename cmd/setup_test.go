package cmd

import (
	"github.com/alin-io/pkgproxy/db"
	"github.com/alin-io/pkgproxy/models"
	"github.com/alin-io/pkgproxy/router"
	"github.com/alin-io/pkgproxy/storage"
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
