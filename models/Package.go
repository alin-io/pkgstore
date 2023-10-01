package models

import (
	"github.com/alin-io/pkgstore/db"
	"gorm.io/gorm"
)

type Package[MetaType any] struct {
	gorm.Model

	Id            uint64                     `gorm:"column:id;primaryKey;autoincrement" json:"id" binding:"required"`
	Name          string                     `gorm:"column:name;uniqueIndex:name_service;not null" json:"name" binding:"required"`
	Service       string                     `gorm:"column:service;uniqueIndex:name_service;not null" json:"service" binding:"required"`
	AuthId        string                     `gorm:"column:auth_id;index;not null" json:"auth_id" binding:"required"`
	LatestVersion string                     `gorm:"column:latest_version" json:"latest_version"`
	Versions      []PackageVersion[MetaType] `gorm:"foreignKey:PackageId;references:Id;constraint:OnDelete:CASCADE;" json:"versions"`
}

func (*Package[T]) TableName() string {
	return "packages"
}

func (p *Package[T]) FillByName(name, service string) error {
	return db.DB().Find(&p, "name = ? AND service = ?", name, service).Error
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
	if p.Id == 0 {
		return version, nil
	}

	err := db.DB().Find(&version, "package_id = ? AND version = ?", p.Id, name).Error
	if err != nil {
		return version, err
	}
	return version, nil
}

func (p *Package[T]) Insert() error {
	return db.DB().Create(p).Error
}

func (p *Package[T]) InsertVersion(version PackageVersion[T]) error {
	if p.Id == 0 {
		return nil
	}
	version.PackageId = p.Id
	if p.Versions == nil {
		p.Versions = make([]PackageVersion[T], 0)
	}
	p.Versions = append(p.Versions, version)
	return db.DB().Create(&version).Error
}

func (p *Package[T]) Delete() error {
	return db.DB().Delete(&Package[T]{}, "id = ?", p.Id).Error
}
