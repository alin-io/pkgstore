package pypi

import (
	"fmt"
	"github.com/alin-io/pkgproxy/config"
	"github.com/alin-io/pkgproxy/db"
	"github.com/alin-io/pkgproxy/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func (s *Service) MetadataHandler(c *gin.Context) {
	pkgName := c.GetString("pkgName")
	pkg := models.Package[pypiPackageMetadata]{}
	versions := make([]models.PackageVersion[pypiPackageMetadata], 0)
	db.DB().Find(&pkg, "name = ?", pkgName)
	db.DB().Find(&versions, "package_id = ?", pkg.Id)
	if pkg.Id < 1 || len(versions) == 0 {
		s.ProxyToPypi(c)
		return
	}

	versionLinks := ""
	for _, versionData := range versions {
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

func (s *Service) ProxyToPypi(c *gin.Context) {
	urlPath := "/simple" + c.Param("path")
	remote, err := url.Parse("https://pypi.org" + urlPath)
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
