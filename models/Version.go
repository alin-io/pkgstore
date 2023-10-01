package models

import (
	"errors"
	"github.com/alin-io/pkgstore/db"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

var digestRegex = regexp.MustCompile(`^[a-f0-9]{64}$`)

type PackageVersion[MetaType any] struct {
	gorm.Model

	Id      uint64 `gorm:"column:id;primaryKey;autoincrement" json:"id" binding:"required"`
	Service string `gorm:"column:service;not null" json:"service" binding:"required"`

	// Single Digest of an Asset or multiple with comma separated values - <sha256>,<sha256>,...
	Digest string `gorm:"column:digest;index" json:"digest"`

	PackageId uint64 `gorm:"column:package_id;uniqueIndex:pkg_id_version;uniqueIndex:pkg_id_tag" json:"package_id" binding:"required"`

	Version string `gorm:"column:version;not null;uniqueIndex:pkg_id_version" json:"version" binding:"required"`
	Tag     string `gorm:"column:tag;uniqueIndex:pkg_id_tag" json:"tag"`

	Metadata datatypes.JSONType[MetaType] `gorm:"column:metadata" json:"metadata"`
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
	return db.DB().Find(p, "digest LIKE ?", "%"+digest+"%").Error
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
	return db.DB().Delete(&PackageVersion[T]{}, "id = ?", p.Id).Error
}

func (p *PackageVersion[T]) Assets() ([]Asset, error) {
	assets := make([]Asset, 0)
	if len(p.Digest) == 0 {
		return assets, nil
	}
	digests := strings.Split(p.Digest, ",")
	err := db.DB().Find(&assets, "digest IN ?", digests).Error
	return assets, err
}
