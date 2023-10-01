package main

import (
	"github.com/alin-io/pkgstore/config"
	"github.com/alin-io/pkgstore/db"
	_ "github.com/alin-io/pkgstore/db"
	"github.com/alin-io/pkgstore/models"
	"github.com/alin-io/pkgstore/router"
	"github.com/alin-io/pkgstore/storage"
)

func main() {
	var storageBackend storage.BaseStorageBackend
	if config.Get().Storage.ActiveBackend == config.StorageS3 {
		storageBackend = storage.NewS3Backend()
	} else {
		panic("Unknown storage backend")
	}

	// Initialize the DB connection
	db.InitDatabase()

	// Sync Models with the DB
	models.SyncModels()

	r := router.SetupGinServer()
	router.PackageRouter(r, storageBackend)

	err := r.Run(config.Get().ListenAddress)
	if err != nil {
		panic(err)
	}
}
