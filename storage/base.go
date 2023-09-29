package storage

import (
	"io"
)

type BaseStorageBackend interface {
	// GetFile Get the package from the storage backend
	GetFile(key string) (io.ReadCloser, error)
	// GetMetadata Get the package JSON metadata from the storage backend
	GetMetadata(key string, value interface{}) error
	// WriteFile Write the package to the storage backend
	WriteFile(key string, metadata interface{}, r io.Reader) error
	// CopyFile Copy the package from the storage backend
	CopyFile(fromKey, toKey string) error
	// DeleteFile Delete the package from the storage backend
	DeleteFile(key string) error
}
