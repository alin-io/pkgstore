package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/alin-io/pkgstore/services/npm"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNpmAuthentication(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/npm/some-package-name", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("some-username:")))
	serverApp.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestNpmPackageUpload(t *testing.T) {
	pkgName := uuid.NewString()
	w, req := UploadTestNpmPackage(pkgName, "0.0.1")
	serverApp.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	versionInfo := npm.MetadataResponse{}
	err := json.Unmarshal(w.Body.Bytes(), &versionInfo)
	assert.Nil(t, err)
	assert.Equal(t, "0.0.1", versionInfo.DistTags["latest"])
	assert.Equal(t, pkgName, versionInfo.Versions["0.0.1"].Name)

	err = DeleteTestPackage(pkgName, "npm")
	assert.Nil(t, err)
}

func TestNpmPackageMetadata(t *testing.T) {
	t.Run("should respond with 404 if requested package doesn't exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/npm/some-package-name", nil)
		serverApp.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
	})

	t.Run("should respond with metadata JSON if package exists", func(t *testing.T) {
		for _, pkgName := range []string{uuid.NewString(), uuid.NewString() + "/" + uuid.NewString()} {
			w, req := UploadTestNpmPackage(pkgName, "0.0.1")
			serverApp.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)

			w.Flush()

			w = httptest.NewRecorder()
			req, _ = http.NewRequest("GET", "/npm/"+pkgName, nil)
			serverApp.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
			err := DeleteTestPackage(pkgName, "npm")
			assert.Nil(t, err)
		}
	})
}

func TestNpmPackageDownload(t *testing.T) {
	pkgName := uuid.NewString()
	version := "0.0.1"
	w, req := UploadTestNpmPackage(pkgName, version)
	serverApp.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	versionInfo := npm.MetadataResponse{}
	err := json.Unmarshal(w.Body.Bytes(), &versionInfo)
	assert.Nil(t, err)

	w.Flush()

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/npm/%[1]s/-/%[1]s-%[2]s.tar.gz", pkgName, version), nil)
	serverApp.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/octet-stream", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), fmt.Sprintf("%[1]s-%[2]s.tar.gz", pkgName, version))
	assert.Equal(t, "354", w.Header().Get("Content-Length"))

	err = DeleteTestPackage(pkgName, "npm")
	assert.Nil(t, err)
}

func UploadTestNpmPackage(name, version string) (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/npm/"+name, NpmPackageDataReader(name, version))
	req.Header.Set("Content-Type", "application/json")
	return w, req
}

func NpmPackageDataReader(name, version string) *bytes.Buffer {
	return bytes.NewBuffer([]byte(fmt.Sprintf(`
{
    "_attachments": {
        "%[1]s-%[2]s.tgz": {
            "content_type": "application/octet-stream",
            "data": "H4sIAAAAAAAAE+1TQUvDMBjdeb/iI4edZEldV2dPwhARPIjiyXlI26zN1iYhSeeK7L+bNJtednMg4l4OKe+9PF7DF0XzNS0ZVmEfr4wUgxODEJLEMRzjPRJyCYPJNCFRlCTE+dzH1PvJqYscQ2ss1a7KT3PCv8DX/kfwMQRAgjYMpYBuIoIzKtwy6MILG6YNl8Jr0XgyvgpswUyuubJ75TGMDuSaUcsKyDooa1C6De6G8t7GRcG2br4CGxKME3wDR1hmrLexvJKwQLdaS52CkOAFMIrlfMlZsUAwGgHbcgsRcid3fdqade9SFz7u9a1naGsrqX3gHbcPNINDyydWcmN1By+W19x2oU7NcyZMfwn3z/PAqTaruanmUix5+V3UXVKq9yEoRZW1yqQYl9zWNBvnssFUcbyJsdJyxXJrcHQdz8gsTg6PzGChGty3H+6Gvz0BZ5xxxn/FJ1EDRNIACAAA",
            "length": 354
        }
    },
    "_id": "{}",
    "description": "Package created by me",
    "dist-tags": {
        "latest": "%[2]s"
    },
    "name": "%[1]s",
    "readme": "ERROR: No README data found!",
    "versions": {
        "%[2]s": {
            "_id": "%[1]s@%[2]s",
            "_nodeVersion": "12.18.4",
            "_npmVersion": "6.14.6",
            "author": {
                "name": "Alin.io Package Registry Utility"
            },
            "description": "Package created by me",
            "dist": {
                "integrity": "sha512-loy16p+Dtw2S43lBmD3Nye+t+Vwv7Tbhv143UN2mwcjaHJyBfGZdNCTXnma3gJCUSE/AR4FPGWEyCOOTJ+ev9g==",
                "shasum": "4a9dbd94ca6093feda03d909f3d7e6bd89d9d4bf",
                "tarball": "http://localhost:8080/npm/%[1]s/-/%[1]s-%[2]s.tgz"
            },
            "keywords": [],
            "license": "ISC",
            "main": "index.js",
            "name": "%[1]s",
            "publishConfig": {
                "registry": "http://localhost:8080/npm"
            },
            "readme": "ERROR: No README data found!",
            "scripts": {
                "test": "echo \"Error: no test specified\" && exit 1"
            },
            "version": "%[2]s"
        }
    }
}
`, name, version)))
}
