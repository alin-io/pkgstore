package pypi

import (
	"fmt"
	"github.com/alin-io/pkgproxy/services"
	"github.com/alin-io/pkgproxy/storage"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
)

type PackageMetadata struct {
	RequiresPython string   `json:"requires_python"`
	OriginalFiles  []string `json:"original_files"`
}

type Service struct {
	services.BasePackageService
}

func NewService(storage storage.BaseStorageBackend) *Service {
	return &Service{
		BasePackageService: services.BasePackageService{
			Prefix:                   "pypi",
			Storage:                  storage,
			PublicRegistryPathPrefix: "/simple/",
			PublicRegistryUrl:        "https://pypi.org",
		},
	}
}

func (s *Service) constructPackageOriginalFilename(name, version, postfix string) string {
	if len(postfix) > 0 {
		postfix = "-" + postfix
	}
	return fmt.Sprintf("%s-%s%s", name, version, postfix)
}

func (s *Service) FilenamePostfix(filename, pkgName, pkgVersionName string) (postfix string) {
	return strings.Replace(filename, s.constructPackageOriginalFilename(pkgName, pkgVersionName, ""), "", 1)
}

func (s *Service) PkgInfoFromRequestPath(c *gin.Context) (pkgName string, filename string) {
	pkgPath := c.Param("path")

	// /:pkgName/
	pattern := `^/(files/|simple/)?([a-z0-9]{64}/)?(?P<pkgName>[^/]+)(?:/)?$`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(pkgPath)
	if matches == nil {
		return "", ""
	}

	for i, name := range re.SubexpNames() {
		if name == "pkgName" {
			if strings.Index(pkgPath, "/files/") == 0 {
				filename = matches[i]
				pkgName, _ = s.PkgVersionFromFilename(filename)
			} else {
				pkgName = matches[i]
			}
		}
	}

	return pkgName, filename
}
