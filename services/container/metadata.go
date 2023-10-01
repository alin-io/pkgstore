package container

import (
	"fmt"
	"github.com/alin-io/pkgproxy/models"
	"github.com/gin-gonic/gin"
)

func (s *Service) MetadataHandler(c *gin.Context) {
	pkgVersionMeta, _, _ := s.pkgVersionMetadata(c)
	if pkgVersionMeta == nil {
		return
	}
	var responseMeta interface{}
	c.Header("Content-Type", pkgVersionMeta.ContentType)
	switch pkgVersionMeta.ContentType {
	case ManifestV1ContentType:
		responseMeta = pkgVersionMeta.ManifestV1
	case ManifestV2ContentType:
		responseMeta = pkgVersionMeta.ManifestV2
	case ManifestListV2ContentType:
		responseMeta = pkgVersionMeta.ManifestListV2
	}
	c.JSON(200, responseMeta)
}

func (s *Service) CheckMetadataHandler(c *gin.Context) {
	pkgVersionMeta, _, pkgVersion := s.pkgVersionMetadata(c)
	if pkgVersionMeta == nil {
		return
	}

	assets, err := pkgVersion.Assets()
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to get package assets"})
		return
	}
	if len(assets) == 0 {
		c.JSON(404, gin.H{"error": "Package not found"})
		return
	}
	c.Header("Content-Length", fmt.Sprintf("%d", assets[0].Size))
	c.Header("Docker-Content-Digest", "sha256:"+assets[0].Digest)
	c.Status(200)
	c.Done()
}

func (s *Service) pkgVersionMetadata(c *gin.Context) (meta *PackageMetadata, pkg models.Package[PackageMetadata], pkgVersion models.PackageVersion[PackageMetadata]) {
	name := s.ConstructFullPkgName(c)
	tag := c.Param("reference")
	pkg = models.Package[PackageMetadata]{}
	err := pkg.FillByName(name, s.Prefix)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error while trying to get package info"})
		return
	}
	if pkg.Id < 1 {
		c.JSON(404, gin.H{"error": "Package not found"})
		return
	}
	pkgVersion, err = pkg.Version(tag)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error while trying to get package info"})
		return
	}
	if pkgVersion.Id < 1 {
		c.JSON(404, gin.H{"error": "Package version not found"})
		return
	}
	pkgMeta := pkgVersion.Metadata.Data()
	meta = &pkgMeta
	return
}
