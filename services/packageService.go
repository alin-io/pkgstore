package services

import (
	"fmt"
	"github.com/alin-io/pkgproxy/storage"
	"github.com/gin-gonic/gin"
	"strings"
)

type PackageService interface {
	PackageFilename(digest string) string
	PkgVersionFromFilename(filename string) (pkgName string, version string)
	ShouldHandleRequest(c *gin.Context) bool
	PkgInfoFromRequestPath(c *gin.Context) (pkgName string, filename string)

	UploadHandler(c *gin.Context)
	DownloadHandler(c *gin.Context)
	MetadataHandler(c *gin.Context)
}

type BasePackageService struct {
	PackageService

	Prefix  string
	Storage storage.BaseStorageBackend
}

func (s *BasePackageService) PackageFilename(digest string) string {
	return fmt.Sprintf("%s/%s.tgs", s.Prefix, digest)
}

func (s *BasePackageService) PkgVersionFromFilename(filename string) (pkgName string, version string) {
	filenameSplit := strings.Split(filename, "-")
	pkgName = filenameSplit[0]
	version = strings.Replace(filenameSplit[1], ".tgz", "", 1)
	return pkgName, version
}