package cmd

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/alin-io/pkgstore/config"
	"github.com/alin-io/pkgstore/services/npm"
	"github.com/alin-io/pkgstore/services/pypi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPypiAuthentication(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/pypi/some-package-name", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("some-username:")))
	serverApp.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestPypiPackageUpload(t *testing.T) {
	pkgName := uuid.NewString()
	version := "0.0.1"
	w, req := UploadTestPypiPackage(pkgName, version)
	serverApp.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	versionInfo := pypi.PackageMetadata{}
	err := json.Unmarshal(w.Body.Bytes(), &versionInfo)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("%s-%s.tar.gz", pkgName, version), versionInfo.OriginalFiles[0])
	assert.Equal(t, "", versionInfo.RequiresPython)

	err = DeleteTestPackage(pkgName, "pypi")
	assert.Nil(t, err)
}

func TestPackageMetadata(t *testing.T) {
	t.Run("should respond with 404 if requested package doesn't exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/pypi/simple/some-package-name", nil)
		serverApp.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})

	t.Run("should respond with metadata JSON if package exists", func(t *testing.T) {
		for _, pkgName := range []string{uuid.NewString(), uuid.NewString() + "/" + uuid.NewString()} {
			w, req := UploadTestPypiPackage(pkgName, "0.0.1")
			serverApp.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
			w.Flush()

			w = httptest.NewRecorder()
			req, _ = http.NewRequest("GET", "/pypi/simple/"+pkgName, nil)
			serverApp.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
			assert.Contains(t, w.Body.String(), fmt.Sprintf(`href="%s/files/`, config.Get().RegistryHost))
			err := DeleteTestPackage(pkgName, "pypi")
			assert.Nil(t, err)
		}
	})
}

func TestPypiPackageDownload(t *testing.T) {
	pkgName := uuid.NewString()
	version := "0.0.1"
	w, req := UploadTestPypiPackage(pkgName, version)
	serverApp.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	versionInfo := npm.MetadataResponse{}
	err := json.Unmarshal(w.Body.Bytes(), &versionInfo)
	assert.Nil(t, err)

	w.Flush()

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/pypi/files/%[1]s/%[2]s-%[3]s.tar.gz", fmt.Sprintf("%x", sha256.Sum256([]byte(time.Now().String()))), pkgName, version), nil)
	serverApp.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), fmt.Sprintf("%[1]s-%[2]s.tar.gz", pkgName, version))
	assert.Equal(t, "1024", w.Header().Get("Content-Length"))

	err = DeleteTestPackage(pkgName, "pypi")
	assert.Nil(t, err)
}

func UploadTestPypiPackage(name, version string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	filename := fmt.Sprintf("%s-%s.tar.gz", name, version)
	bodyBuffer := bytes.NewBuffer([]byte{})
	formWriter := multipart.NewWriter(bodyBuffer)
	randomBytes := make([]byte, 1024)
	_, _ = rand.Read(randomBytes)

	part, _ := formWriter.CreateFormFile("content", filename)
	_, _ = part.Write(randomBytes)

	part, _ = formWriter.CreateFormField("name")
	_, _ = part.Write([]byte(name))

	part, _ = formWriter.CreateFormField("version")
	_, _ = part.Write([]byte(version))

	_ = formWriter.Close()

	req, _ := http.NewRequest("POST", "/pypi", bodyBuffer)
	req.Header.Set("Content-Type", formWriter.FormDataContentType())
	return w, req
}
