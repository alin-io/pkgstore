package models

import (
	"github.com/alin-io/pkgstore/db"
	"github.com/google/uuid"
)

func SyncModels() {
	err := db.DB().AutoMigrate(&Package[any]{}, &PackageVersion[any]{}, Asset{})
	if err != nil {
		panic(err)
	}
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
