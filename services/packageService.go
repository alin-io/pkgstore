package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/alin-io/pkgproxy/storage"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"
)

type PackageService interface {
	PackageFilename(digest, postfix string) string
	PkgVersionFromFilename(filename string) (pkgName string, version string)
	PkgInfoFromRequestPath(c *gin.Context) (pkgName string, filename string)

	UploadHandler(c *gin.Context)
	DownloadHandler(c *gin.Context)
	MetadataHandler(c *gin.Context)
}

type BasePackageService struct {
	PackageService

	Prefix  string
	Storage storage.BaseStorageBackend

	PublicRegistryUrl        string
	PublicRegistryPathPrefix string
}

func (s *BasePackageService) PackageFilename(digest, postfix string) string {
	if len(postfix) == 0 {
		postfix = ".tar.gz"
	}
	return fmt.Sprintf("%s/%s%s", s.Prefix, digest, postfix)
}

func (s *BasePackageService) PkgVersionFromFilename(filename string) (pkgName string, version string) {
	base := filepath.Base(filename)
	for _, ext := range []string{".tar.gz", ".tgz", ".whl"} {
		if strings.HasSuffix(base, ext) {
			base = strings.Replace(base, ext, "", 1)
		}
	}

	filenameSplit := strings.Split(base, "-")
	pkgName = strings.Join(filenameSplit[:len(filenameSplit)-1], "-")
	version = filenameSplit[len(filenameSplit)-1]
	return pkgName, version
}

func (s *BasePackageService) ChecksumReader(r io.Reader) (checksum string, size int64, err error) {
	h := sha256.New()
	if size, err = io.Copy(h, r); err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(h.Sum(nil)), size, nil
}

func (s *BasePackageService) ProxyToPublicRegistry(c *gin.Context) {
	urlPath := s.PublicRegistryUrl + c.Param("path")
	remote, err := url.Parse(s.PublicRegistryPathPrefix + urlPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header

		// Remove Authorization header
		req.Header.Del("Authorization")

		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = urlPath
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
