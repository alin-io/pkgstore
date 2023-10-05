package api

import "github.com/alin-io/pkgstore/storage"

type Service struct {
	Storage storage.BaseStorageBackend
}

func NewApiService(storageBackend storage.BaseStorageBackend) *Service {
	return &Service{
		Storage: storageBackend,
	}
}
