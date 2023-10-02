package pypi

import (
	"github.com/alin-io/pkgstore/models"
	"github.com/gin-gonic/gin"
	"io"
	"log"
)

func (s *Service) DownloadHandler(c *gin.Context) {
	filename := c.Param("filename")
	pkgName, version := s.PkgVersionFromFilename(filename)
	pkg := models.Package[PackageMetadata]{}
	versionInfo := models.PackageVersion[PackageMetadata]{}
	err := pkg.FillByName(pkgName, s.Prefix)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error while trying to get package info"})
		return
	}

	if pkg.Id == 0 {
		c.JSON(404, gin.H{"error": "Not Found"})
		return
	}

	versionInfo, err = pkg.Version(version)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error while trying to get package info"})
		return
	}

	if len(versionInfo.Digest) == 0 {
		c.JSON(404, gin.H{"error": "Not Found"})
		return
	}

	fileAsset, err := versionInfo.Asset()
	if err != nil {
		c.JSON(500, gin.H{"error": "Error while trying to get package info"})
		return
	}

	if fileAsset == nil || fileAsset.Id < 1 {
		c.JSON(404, gin.H{"error": "Not Found"})
		return
	}

	fileData, err := s.Storage.GetFile(s.PackageFilename(fileAsset.Digest))
	if err != nil {
		c.JSON(404, gin.H{"error": "Not Found"})
		return
	}

	defer func(fileData io.ReadCloser) {
		err := fileData.Close()
		if err != nil {
			log.Println(err)
		}
	}(fileData)

	c.DataFromReader(200, int64(fileAsset.Size), "application/octet-stream", fileData, map[string]string{
		"Content-Disposition": "attachment; filename=" + filename,
	})
}
