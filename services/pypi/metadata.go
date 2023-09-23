package pypi

import (
	"fmt"
	"github.com/alin-io/pkgproxy/config"
	"github.com/alin-io/pkgproxy/db"
	"github.com/alin-io/pkgproxy/models"
	"github.com/gin-gonic/gin"
	"strings"
)

func (s *Service) MetadataHandler(c *gin.Context) {
	pkgName := strings.Replace(c.Param("path"), "simple/", "", 1)
	pkg := models.Package[pypiPackageMetadata]{}
	versions := make([]models.PackageVersion[pypiPackageMetadata], 0)
	db.DB().Find(&pkg, "name = ?", pkgName)
	db.DB().Find(&versions, "package_id = ?", pkg.Id)
	versionLinks := ""
	for _, versionData := range versions {
		versionLinks = fmt.Sprintf(
			`%[1]s\n<a href="%[2]s/files/%[3]s/%[4]s-%[5]s-%[6]s#sha256=%[3]s" data-requires-python="%[7]s">%[4]s-%[5]s-%[6]s</a><br>`,
			versionLinks,
			config.Get().RegistryHost,
			versionData.Digest,
			pkgName,
			versionData.Version,
			versionData.Metadata.Data().FilenamePostfix,
			versionData.Metadata.Data().RequiresPython,
		)
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
