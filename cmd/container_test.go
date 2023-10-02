package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContainerAuthentication(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v2/some-package-name", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("some-username:")))
	serverApp.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestContainerPackageUpload(t *testing.T) {
	name := uuid.NewString()
	tag := "latest"
	digest, _ := UploadTestContainerPackage(t, name, tag)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("HEAD", "/v2/"+name+"/manifests/"+tag, nil)
	serverApp.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("HEAD", "/v2/"+name+"/blobs/sha256:"+digest, nil)
	serverApp.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, fmt.Sprintf("sha256:%s", digest), w.Header().Get("Docker-Content-Digest"))
}

func TestContainerPackageDownload(t *testing.T) {
	name := uuid.NewString()
	tag := "latest"
	digest, blobBuffer := UploadTestContainerPackage(t, name, tag)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v2/%[1]s/blobs/sha256:%[2]s", name, digest), nil)
	serverApp.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "filename="+digest)
	assert.Equal(t, fmt.Sprintf("%d", len(blobBuffer)), w.Header().Get("Content-Length"))
}

func UploadTestContainerPackage(t *testing.T, name, tag string) (digest string, blob []byte) {
	blob = []byte("some layer blob")
	blobBuffer := bytes.NewBuffer([]byte("some layer blob"))
	digest = fmt.Sprintf("%x", sha256.Sum256(blobBuffer.Bytes()))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v2/"+name+"/blobs/uploads/", nil)
	serverApp.ServeHTTP(w, req)
	assert.Equal(t, 202, w.Code)
	uploadUUID := w.Header().Get("Docker-Upload-UUID")
	assert.Equal(t, len(uploadUUID), 36)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/v2/"+name+"/blobs/uploads/"+uploadUUID+"?digest=sha256:"+digest, blobBuffer)
	serverApp.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)
	w.Flush()

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/v2/"+name+"/manifests/"+tag, ContainerManifestReader(digest, len(blob)))
	req.Header.Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")
	serverApp.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	w.Flush()

	return
}

func ContainerManifestReader(layerDigest string, size int) *bytes.Buffer {
	return bytes.NewBuffer([]byte(fmt.Sprintf(`{
    "schemaVersion": 2,
    "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
    "config": {
        "mediaType": "application/vnd.docker.container.image.v1+json",
        "size": %[2]d,
        "digest": "sha256:%[1]s"
    },
    "layers": [
        {
            "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
            "size": %[2]d,
            "digest": "sha256:%[1]s"
        }
    ]
}
`, layerDigest, size)))
}
