package cmd

import (
	"encoding/base64"
	"encoding/json"
	"github.com/alin-io/pkgstore/models"
	"github.com/alin-io/pkgstore/services/npm"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiAuthentication(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/packages", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("some-username:")))
	serverApp.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestApiPackagesList(t *testing.T) {
	pkgName := uuid.NewString()
	version := "0.0.1"
	w, req := UploadTestNpmPackage(pkgName, version)
	serverApp.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	versionInfo := npm.MetadataResponse{}
	err := json.Unmarshal(w.Body.Bytes(), &versionInfo)
	assert.Nil(t, err)

	w, req = UploadTestPypiPackage(pkgName, version)
	serverApp.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	_, _ = UploadTestContainerPackage(t, pkgName, version)

	t.Run("should respond with 200 and empty list if no packages are found", func(t *testing.T) {
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/packages", nil)
		serverApp.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		packages := make([]models.Package[any], 0)
		err = json.Unmarshal(w.Body.Bytes(), &packages)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(packages))
	})

	t.Run("should respond with package with a filtered name", func(t *testing.T) {
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/packages?q="+pkgName, nil)
		serverApp.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		packages := make([]models.Package[any], 0)
		err = json.Unmarshal(w.Body.Bytes(), &packages)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(packages))
	})
}

func TestApiPackageVersions(t *testing.T) {
	pkgName := uuid.NewString()
	version := "0.0.1"
	w, req := UploadTestNpmPackage(pkgName, version)
	serverApp.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	versionInfo := npm.MetadataResponse{}
	err := json.Unmarshal(w.Body.Bytes(), &versionInfo)
	assert.Nil(t, err)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/packages", nil)
	serverApp.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	packages := make([]models.Package[any], 0)
	err = json.Unmarshal(w.Body.Bytes(), &packages)
	assert.Nil(t, err)
	pkg := models.Package[any]{}
	for _, p := range packages {
		if p.Name == pkgName {
			pkg = p
			break
		}
	}
	assert.NotEqual(t, uuid.Nil, pkg.ID)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/packages/"+pkg.ID.String()+"/versions", nil)
	serverApp.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	versions := make([]models.PackageVersion[any], 0)
	err = json.Unmarshal(w.Body.Bytes(), &versions)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(versions))
	assert.Equal(t, version, versions[0].Version)
}

func TestApiPackageDelete(t *testing.T) {
	pkgName := uuid.NewString()
	version := "0.0.1"
	_, _ = UploadTestContainerPackage(t, pkgName, version)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/packages?q="+pkgName, nil)
	serverApp.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	packages := make([]models.Package[any], 0)
	err := json.Unmarshal(w.Body.Bytes(), &packages)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(packages))
	pkg := packages[0]
	assert.NotEqual(t, uuid.Nil, pkg.ID)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/packages/"+pkg.ID.String(), nil)
	serverApp.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/packages/"+pkg.ID.String()+"/versions", nil)
	serverApp.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}
