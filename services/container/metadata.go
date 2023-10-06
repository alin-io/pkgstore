package container

import (
	"fmt"
	"github.com/alin-io/pkgstore/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
)

func (s *Service) MetadataHandler(c *gin.Context) {
	_, pkgVersion := s.pkgVersionMetadata(c)
	if pkgVersion.ID == uuid.Nil {
		return
	}

	metadata := pkgVersion.Metadata.Data()

	c.Header("Docker-Content-Digest", "sha256:"+pkgVersion.Digest)
	c.Data(200, metadata.ContentType, metadata.MetadataBuffer)
}

func (s *Service) CheckMetadataHandler(c *gin.Context) {
	_, pkgVersion := s.pkgVersionMetadata(c)
	if pkgVersion.ID == uuid.Nil {
		return
	}

	metadata := pkgVersion.Metadata.Data()

	c.Header("Content-Type", metadata.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", len(metadata.MetadataBuffer)))
	c.Header("Docker-Content-Digest", fmt.Sprintf("sha256:%s", pkgVersion.Digest))
	c.Status(200)
	c.Done()
}

func (s *Service) pkgVersionMetadata(c *gin.Context) (pkg models.Package[PackageMetadata], pkgVersion models.PackageVersion[PackageMetadata]) {
	name := s.ConstructFullPkgName(c)
	tagOrDigest := c.Param("reference")
	pkg = models.Package[PackageMetadata]{}
	err := pkg.FillByName(name, s.Prefix)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error while trying to get package info"})
		return
	}
	if pkg.ID == uuid.Nil {
		c.JSON(404, gin.H{"error": "Package not found"})
		return
	}
	if strings.Contains(tagOrDigest, "sha256:") {
		err = pkgVersion.FillByDigest(strings.Replace(tagOrDigest, "sha256:", "", 1))
	} else {
		pkgVersion, err = pkg.Version(tagOrDigest)
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "Error while trying to get package info"})
		return
	}
	if pkgVersion.ID == uuid.Nil || pkgVersion.PackageId != pkg.ID {
		c.JSON(404, gin.H{"error": "Package version not found"})
		return
	}
	return
}
