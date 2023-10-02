package container

import (
	"github.com/alin-io/pkgstore/services"
	"github.com/alin-io/pkgstore/storage"
	"github.com/gin-gonic/gin"
	"regexp"
)

const (
	ManifestV1ContentType         = "application/vnd.docker.distribution.manifest.v1+json"
	ManifestV2ContentType         = "application/vnd.docker.distribution.manifest.v2+json"
	ManifestListV2ContentType     = "application/vnd.docker.distribution.manifest.list.v2+json"
	ManifestOCIV1ContentType      = "application/vnd.oci.image.manifest.v1+json"
	ManifestOCIIndexV1ContentType = "application/vnd.oci.image.index.v1+json"
)

type PackageMetadata struct {
	ContentType    string `json:"contentType"`
	Digest         string `json:"digest"`
	MetadataBuffer []byte `json:"metadataBuffer"`
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

type ManifestV1 struct {
	Name     string `json:"name,omitempty"`
	Tag      string `json:"tag,omitempty"`
	FsLayers []struct {
		BlobSum string `json:"blobSum,omitempty"`
	} `json:"fsLayers,omitempty"`
	History []struct {
		V1Compatibility string `json:"v1Compatibility,omitempty"`
	} `json:"history,omitempty"`
	SchemaVersion int `json:"schemaVersion,omitempty"`
	Signatures    []struct {
		Header struct {
			Jwk struct {
				Crv string `json:"crv,omitempty"`
				Kid string `json:"kid,omitempty"`
				Kty string `json:"kty,omitempty"`
				X   string `json:"x,omitempty"`
				Y   string `json:"y,omitempty"`
			} `json:"jwk"`
			Alg string `json:"alg,omitempty"`
		} `json:"header,omitempty"`
		Signature string `json:"signature,omitempty"`
		Protected string `json:"protected,omitempty"`
	} `json:"signatures,omitempty"`
}

type ManifestListV2 struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"` // application/vnd.docker.distribution.manifest.list.v2+json
	ArtifactType  string `json:"artifactType,omitempty"`
	Manifests     []struct {
		MediaType string `json:"mediaType"` // application/vnd.docker.distribution.manifest.v2+json
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
		Platform  struct {
			Architecture string   `json:"architecture,omitempty"`
			Os           string   `json:"os,omitempty"`
			OsVersion    string   `json:"os.version,omitempty"`
			OsFeatures   []string `json:"os.features,omitempty"`
			Variant      string   `json:"variant,omitempty"`
			Features     []string `json:"features,omitempty"`
		} `json:"platform,omitempty"`
	} `json:"manifests,omitempty"`
	Subject struct {
		MediaType    string            `json:"mediaType,omitempty"`
		Size         int               `json:"size,omitempty"`
		Digest       string            `json:"digest,omitempty"`
		Urls         []string          `json:"urls,omitempty"`
		Annotations  map[string]string `json:"annotations,omitempty"`
		Data         string            `json:"data,omitempty"`
		ArtifactType string            `json:"artifactType,omitempty"`
	} `json:"subject,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type ManifestV2 struct {
	SchemaVersion int    `json:"schemaVersion,omitempty"`
	MediaType     string `json:"mediaType,omitempty"`
	Config        struct {
		MediaType string `json:"mediaType,omitempty"`
		Size      int    `json:"size,omitempty"`
		Digest    string `json:"digest,omitempty"`
	} `json:"config"`
	Layers []struct {
		MediaType   string            `json:"mediaType,omitempty"`
		Size        int               `json:"size,omitempty"`
		Digest      string            `json:"digest,omitempty"`
		Annotations map[string]string `json:"annotations,omitempty"`
	} `json:"layers"`
}
