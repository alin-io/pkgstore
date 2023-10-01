package npm

import (
	"github.com/alin-io/pkgstore/services"
	"github.com/alin-io/pkgstore/storage"
)

type Service struct {
	services.BasePackageService
}

type PackageMetadata struct {
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
		BasePackageService: services.BasePackageService{
			Prefix:                   "npm",
			Storage:                  storage,
			PublicRegistryPathPrefix: "",
			PublicRegistryUrl:        "https://registry.npmjs.org",
		},
	}
}
