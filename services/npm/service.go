package npm

import (
	"github.com/alin-io/pkgproxy/services"
	"github.com/alin-io/pkgproxy/storage"
	"github.com/gin-gonic/gin"
	"regexp"
)

type Service struct {
	services.BasePackageService
}

type npmPackageMetadata struct {
	Id          string            `json:"_id"`
	Description string            `json:"description"`
	Readme      string            `json:"readme"`
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	NodeVersion string            `json:"_nodeVersion"`
	NpmVersion  string            `json:"_npmVersion"`
	Author      map[string]string `json:"author"`
	Dist        struct {
		Integrity string `json:"integrity"`
		Shasum    string `json:"shasum"`
		Tarball   string `json:"tarball"`
	} `json:"dist"`
	PublishConfig map[string]string `json:"publishConfig"`
	Scripts       map[string]string `json:"scripts"`
	Keywords      []string          `json:"keywords"`
	License       string            `json:"license"`
	Main          string            `json:"main"`
}

func NewService(storage storage.BaseStorageBackend) *Service {
	return &Service{
		BasePackageService: services.BasePackageService{Prefix: "npm", Storage: storage},
	}
}

func (s *Service) PkgInfoFromRequestPath(c *gin.Context) (pkgName string, filename string) {
	pkgPath := c.Param("path")

	// /:pkgName/-/:filename
	// /@orgname/pkgName/-/:filename
	pattern := `^/(?P<pkgName>(@[^/]+/)?[^/]+)(?:/-/)(?P<filename>[^/]+)$`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(pkgPath)
	if matches == nil {
		return "", ""
	}

	for i, name := range re.SubexpNames() {
		if name == "pkgName" {
			pkgName = matches[i]
		} else if name == "filename" {
			filename = matches[i]
		}
	}

	return pkgName, filename
}
