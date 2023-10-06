package npm

import (
	"github.com/alin-io/pkgstore/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MetadataResponse struct {
	Name     string                     `json:"name"`
	DistTags map[string]string          `json:"dist-tags"`
	Versions map[string]PackageMetadata `json:"versions"`
}

func (s *Service) MetadataHandler(c *gin.Context) {
	pkgName := s.ConstructFullPkgName(c)
	pkg := models.Package[PackageMetadata]{}
	err := pkg.FillByName(pkgName, s.Prefix)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error while trying to get package info"})
		return
	}

	err = pkg.FillVersions()
	if err != nil {
		c.JSON(500, gin.H{"error": "Error while trying to get package info"})
		return
	}

	if !c.GetBool("testing") && (pkg.ID == uuid.Nil || len(pkg.Versions) == 0) {
		s.ProxyToPublicRegistry(c)
		return
	}

	if pkg.ID == uuid.Nil || len(pkg.Versions) == 0 {
		c.JSON(404, gin.H{"error": "Package not found"})
		return
	}

	result := MetadataResponse{Name: pkgName, DistTags: make(map[string]string), Versions: make(map[string]PackageMetadata)}

	for _, version := range pkg.Versions {
		result.Versions[version.Version] = version.Metadata.Data()
		if len(version.Tag) > 0 {
			result.DistTags[version.Tag] = version.Version
		}
	}
	c.JSON(200, result)
}
