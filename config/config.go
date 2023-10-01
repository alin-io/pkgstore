package config

import (
	// Autoload .env file
	_ "github.com/joho/godotenv/autoload"

	"os"
)

var (
	projectConfig = ProjectConfigType{}
)

const (
	StorageS3 = "s3"
)

func init() {
	projectConfig.ListenAddress = getEnv("LISTEN_ADDRESS", ":8080")
	projectConfig.RegistryHost = getEnv("REGISTRY_HOST", "http://localhost:8080")
	projectConfig.DatabaseUrl = getEnv("DATABASE_URL", "file::memory:?cache=shared")

	// Storage Backend
	projectConfig.Storage.ActiveBackend = getEnv("STORAGE_BACKEND", StorageS3)

	// S3 Storage Config
	projectConfig.Storage.S3.Region = getEnv("S3_REGION", "us-east-1")
	projectConfig.Storage.S3.Bucket = getEnv("S3_BUCKET", "pkgstore")
	projectConfig.Storage.S3.ApiKey = getEnv("S3_API_KEY", "minioadmin")
	projectConfig.Storage.S3.ApiSecret = getEnv("S3_API_SECRET", "minioadmin")
	projectConfig.Storage.S3.ApiHost = getEnv("S3_API_HOST", "")
}

func Get() *ProjectConfigType {
	return &projectConfig
}

type ProjectConfigType struct {
	ListenAddress string
	DatabaseUrl   string
	RegistryHost  string
	Storage       struct {
		ActiveBackend string
		S3            struct {
			Region    string
			Bucket    string
			ApiKey    string
			ApiSecret string
			ApiHost   string
		}
	}
}

func getEnv(keyAndFallback ...string) string {
	key := keyAndFallback[0]
	value := os.Getenv(key)
	if len(value) == 0 {
		if len(keyAndFallback) > 1 {
			return keyAndFallback[1]
		}
		panic("Missing environment variable - " + key)
	}
	return value
}
