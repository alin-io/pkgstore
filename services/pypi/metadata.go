package pypi

import (
	"fmt"
	"github.com/alin-io/pkgstore/config"
	"github.com/alin-io/pkgstore/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

	versionLinks := ""
	for _, versionData := range pkg.Versions {
		for _, originalFilename := range versionData.Metadata.Data().OriginalFiles {
			versionLinks = fmt.Sprintf(
				`%[1]s<a href="%[2]s/files/%[3]s/%[4]s#sha256=%[3]s" data-requires-python="%[5]s">%[4]s</a></br>`,
				versionLinks,
				config.Get().RegistryHost,
				versionData.Digest,
				originalFilename,
				versionData.Metadata.Data().RequiresPython,
			)
		}
	}

	if pkg.ID == uuid.Nil || len(pkg.Versions) == 0 {
		c.JSON(404, gin.H{"error": "Package not found"})
		return
	}

	c.Data(200, "text/html; charset=utf-8", []byte(fmt.Sprintf(`<!DOCTYPE html>
<html>
  <head>
    <title>Links for %[1]s</title>
  </head>
  <body>
    <h1>Links for %[1]s</h1>
    %[2]s
  </body>
</html>
`, pkgName, versionLinks)))
}
