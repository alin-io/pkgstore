package models

import (
	"github.com/alin-io/pkgproxy/db"
	"gorm.io/gorm"
)

type Package[MetaType any] struct {
	gorm.Model

	Id            uint64                     `gorm:"column:id;primaryKey;autoincrement" json:"id" binding:"required"`
	Name          string                     `gorm:"column:name;unique;not null" json:"name" binding:"required"`
	Namespace     string                     `gorm:"column:namespace;index;not null" json:"namespace" binding:"required"`
	Service       string                     `gorm:"column:service;index;not null" json:"service" binding:"required"`
	AuthId        string                     `gorm:"column:auth_id;index;not null" json:"auth_id" binding:"required"`
	LatestVersion string                     `gorm:"column:latest_version" json:"latest_version"`
	Versions      []PackageVersion[MetaType] `gorm:"foreignKey:PackageId;references:Id;constraint:OnDelete:CASCADE;" json:"versions"`
}

func (*Package[T]) TableName() string {
	return "packages"
}

func (p *Package[T]) FillByName(name string) error {
	return db.DB().Find(&p, "name = ?", name).Error
}

func (p *Package[T]) FillVersions() error {
	if p.Versions == nil {
		p.Versions = make([]PackageVersion[T], 0)
	}
	if p.Id == 0 {
		return nil
	}
	return db.DB().Find(&p.Versions, "package_id = ?", p.Id).Error
}

func (p *Package[T]) Version(name string) (PackageVersion[T], error) {
	version := PackageVersion[T]{}
	err := db.DB().Find(&version, "package_id = ? AND version = ?", p.Id, name).Error
	if err != nil {
		return version, err
	}
	return version, nil
}

func (p *Package[T]) Insert() error {
	return db.DB().Create(p).Error
}

func (p *Package[T]) Delete() error {
	return db.DB().Delete(&Package[T]{}, "id = ?", p.Id).Error
}
