package models

import (
	"errors"
	"github.com/alin-io/pkgstore/db"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"regexp"
	"time"
)

var digestRegex = regexp.MustCompile(`^[a-f0-9]{64}$`)

type PackageVersion[MetaType any] struct {
	gorm.Model `json:"-"`

	ID      uuid.UUID `gorm:"column:id;primaryKey;" json:"id" binding:"required"`
	Service string    `gorm:"column:service;not null" json:"service" binding:"required"`

	Digest string `gorm:"column:digest;index" json:"digest"`

	PackageId uuid.UUID `gorm:"column:package_id;uniqueIndex:pkg_id_version;uniqueIndex:pkg_id_tag;" json:"package_id" binding:"required"`

	Version string `gorm:"column:version;not null;uniqueIndex:pkg_id_version" json:"version" binding:"required"`
	Tag     string `gorm:"column:tag;uniqueIndex:pkg_id_tag" json:"tag"`

	Metadata datatypes.JSONType[MetaType] `gorm:"column:metadata" json:"metadata"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (p *PackageVersion[T]) BeforeCreate(_ *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}

func (*PackageVersion[T]) TableName() string {
	return "package_versions"
}

func (p *PackageVersion[T]) FillByName(version string) error {
	return db.DB().Find(p, "version = ?", version).Error
}

func (p *PackageVersion[T]) FillById(id uint64) error {
	return db.DB().Find(p, "id = ?", id).Error
}

func (p *PackageVersion[T]) FillByDigest(digest string) error {
	match := digestRegex.MatchString(digest)
	if !match {
		return errors.New("invalid digest")
	}
	return db.DB().Find(p, "digest = ?", digest).Error
}

func (p *PackageVersion[T]) Insert() error {
	return db.DB().Create(p).Error
}

func (p *PackageVersion[T]) SaveMeta() error {
	return db.DB().Model(p).Update("metadata", p.Metadata).Error
}

func (p *PackageVersion[T]) Save() error {
	return db.DB().Save(p).Error
}

func (p *PackageVersion[T]) Delete() error {
	return db.DB().Delete(&PackageVersion[T]{}, "id = ?", p.ID).Error
}

func (p *PackageVersion[T]) Asset() (*Asset, error) {
	asset := &Asset{}
	if len(p.Digest) == 0 {
		return nil, nil
	}
	err := db.DB().Where("digest = ?", p.Digest).Find(&asset).Error
	return asset, err
}
