package container

import (
	"github.com/alin-io/pkgstore/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"strings"
)

func (s *Service) DownloadHandler(c *gin.Context) {
	name := s.ConstructFullPkgName(c)
	inputDigest := c.Param("sha256")
	digest := strings.Replace(inputDigest, "sha256:", "", 1)
	pkg := models.Package[PackageMetadata]{}
	err := pkg.FillByName(name, s.Prefix)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error while trying to get package info"})
		return
	}
	if pkg.ID == uuid.Nil {
		c.JSON(404, gin.H{"error": "Package not found"})
		return
	}

	asset := models.Asset{}
	err = asset.FillByDigest(digest)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to check the DB for package version"})
		return
	}
	if asset.Digest != digest {
		c.JSON(404, gin.H{"error": "Uploaded asset not found"})
		return
	}

	fileData, err := s.Storage.GetFile(s.PackageFilename(asset.Digest))
	if err != nil || fileData == nil {
		c.JSON(404, gin.H{"error": "Not Found"})
		return
	}

	defer func(fileData io.ReadCloser) {
		err := fileData.Close()
		if err != nil {
			log.Println(err)
		}
	}(fileData)

	c.DataFromReader(200, int64(asset.Size), "application/octet-stream", fileData, map[string]string{
		"Content-Disposition": "attachment; filename=" + digest,
	})
}
