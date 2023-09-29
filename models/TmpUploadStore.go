package models

import (
	"github.com/alin-io/pkgproxy/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TmpUploadStore struct {
	gorm.Model

	Id          string `gorm:"column:id;primaryKey" json:"id" binding:"required"`
	Name        string `gorm:"column:name;not null" json:"name" binding:"required"`
	UploadRange string `gorm:"column:upload_range;not null" json:"upload_range" binding:"required"`
	Digest      string `gorm:"column:digest;index" json:"digest" binding:"required"`
	Size        uint64 `gorm:"column:size;not null" json:"size" binding:"required"`
}

func (*TmpUploadStore) TableName() string {
	return "tmp_upload_store"
}

func (t *TmpUploadStore) Insert() error {
	if t.Id == "" {
		t.Id = uuid.NewString()
	}
	return db.DB().Create(t).Error
}

func (t *TmpUploadStore) FillByDigest(digest, name string) error {
	return db.DB().Find(t, "digest = ? AND name = ?", digest, name).Error
}

func (t *TmpUploadStore) FillById(id string) error {
	return db.DB().Find(t, "id = ?", id).Error
}

func (t *TmpUploadStore) Update() error {
	return db.DB().Save(t).Error
}

func (t *TmpUploadStore) Delete() error {
	return db.DB().Delete(t).Error
}
