package container

import (
	"github.com/alin-io/pkgproxy/services"
	"github.com/alin-io/pkgproxy/storage"
	"github.com/gin-gonic/gin"
	"regexp"
)

type PackageMetadata struct {
	Name     string `json:"name"`
	Tag      string `json:"tag"`
	FsLayers []struct {
		BlobSum string `json:"blobSum"`
	} `json:"fsLayers"`
	History []struct {
		V1Compatibility string `json:"v1Compatibility"`
	} `json:"history"`
	SchemaVersion int `json:"schemaVersion"`
	Signatures    []struct {
		Header struct {
			Jwk struct {
				Crv string `json:"crv"`
				Kid string `json:"kid"`
				Kty string `json:"kty"`
				X   string `json:"x"`
				Y   string `json:"y"`
			} `json:"jwk"`
			Alg string `json:"alg"`
		} `json:"header"`
		Signature string `json:"signature"`
		Protected string `json:"protected"`
	} `json:"signatures"`
}

type Service struct {
	services.BasePackageService
}

func NewService(storage storage.BaseStorageBackend) *Service {
	return &Service{
		BasePackageService: services.BasePackageService{
			Prefix:                   "container",
			Storage:                  storage,
			PublicRegistryPathPrefix: "/v2/",
			PublicRegistryUrl:        "https://registry.hub.docker.com",
		},
	}
}

func (s *Service) PkgInfoFromRequest(c *gin.Context) (pkgName string, filename string) {
	pkgPath := c.Param("path")

	// /:pkgName/
	pattern := `^/v2/(?P<pkgName>([^/]+/)?[^/]+)/(blob/|manifest/)(?P<filename>[a-z0-9]+)(?:/)?$`
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
