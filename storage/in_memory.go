package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

type InMemoryBackend struct {
	BaseStorageBackend

	storage map[string]InMemoryFile
}

type InMemoryFile struct {
	name     string
	data     []byte
	fileMeta interface{}
}

func NewInMemoryBackend() *InMemoryBackend {
	return &InMemoryBackend{
		storage: make(map[string]InMemoryFile),
	}
}

func (s *InMemoryBackend) WriteFile(key string, fileMeta interface{}, r io.Reader) error {
	fileBuffer := bytes.NewBuffer([]byte{})
	_, err := io.Copy(fileBuffer, r)
	if err != nil {
		return err
	}
	s.storage[key] = InMemoryFile{
		name:     key,
		data:     fileBuffer.Bytes(),
		fileMeta: fileMeta,
	}
	return nil
}

func (s *InMemoryBackend) GetFile(key string) (io.ReadCloser, error) {
	file, ok := s.storage[key]
	if !ok {
		return nil, nil
	}
	return io.NopCloser(bytes.NewReader(file.data)), nil
}

func (s *InMemoryBackend) GetMetadata(key string, value interface{}) error {
	file, ok := s.storage[key]
	if !ok {
		return errors.New("file not found")
	}
	fileMetaBytes, err := json.Marshal(file.fileMeta)
	if err != nil {
		return err
	}
	return json.Unmarshal(fileMetaBytes, value)
}

func (s *InMemoryBackend) CopyFile(fromKey, toKey string) error {
	s.storage[toKey] = s.storage[fromKey]
	return nil
}

func (s *InMemoryBackend) DeleteFile(key string) error {
	delete(s.storage, key)
	return nil
}
