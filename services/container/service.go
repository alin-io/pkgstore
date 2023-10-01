package container

import (
	"github.com/alin-io/pkgstore/services"
	"github.com/alin-io/pkgstore/storage"
	"github.com/gin-gonic/gin"
	"regexp"
)

type PackageMetadata struct {
	ContentType string `json:"contentType"`

	// Content Type - application/vnd.docker.distribution.manifest.v1+json
	ManifestV1 ManifestV1 `json:"manifestV1,omitempty"`

	// Content Type - application/vnd.docker.distribution.manifest.v2+json
	ManifestV2 ManifestV2 `json:"manifestV2,omitempty"`

	// Content Type - application/vnd.docker.distribution.manifest.list.v2+json
	ManifestListV2 ManifestListV2 `json:"manifestListV2,omitempty"`
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

const (
	ManifestV1ContentType     = "application/vnd.docker.distribution.manifest.v1+json"
	ManifestV2ContentType     = "application/vnd.docker.distribution.manifest.v2+json"
	ManifestListV2ContentType = "application/vnd.docker.distribution.manifest.list.v2+json"
)

type ManifestV1 struct {
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

type ManifestListV2 struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"` // application/vnd.docker.distribution.manifest.list.v2+json
	Manifests     []struct {
		MediaType string `json:"mediaType"` // application/vnd.docker.distribution.manifest.v2+json
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
		Platform  struct {
			Architecture string   `json:"architecture"`
			Os           string   `json:"os"`
			OsVersion    string   `json:"os.version"`
			OsFeatures   []string `json:"os.features"`
			Variant      string   `json:"variant"`
			Features     []string `json:"features"`
		} `json:"platform"`
	} `json:"manifests"`
}

type ManifestV2 struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"` // application/vnd.docker.distribution.manifest.v2+json
	Config        struct {
		MediaType string `json:"mediaType"` // application/vnd.docker.distribution.manifest.v2+json
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"config"`
	Layers []struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"layers"`
}
