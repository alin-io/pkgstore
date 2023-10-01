package pypi

import (
	"fmt"
	"github.com/alin-io/pkgstore/services"
	"github.com/alin-io/pkgstore/storage"
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
