package api

import (
	"github.com/alin-io/pkgstore/db"
	"github.com/alin-io/pkgstore/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
)

func (s *Service) ListVersionsHandler(c *gin.Context) {
	packageIdString := c.Param("id")
	packageId, err := uuid.Parse(packageIdString)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid package id"})
		return
	}

	pkg := models.Package[any]{}
	err = db.DB().Model(&pkg).Where(`id = ?`, packageId).Preload("Versions").Find(&pkg).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if pkg.ID == uuid.Nil {
		c.JSON(404, gin.H{"error": "Package not found"})
		return
	}
	c.JSON(200, pkg.Versions)
}

func (s *Service) DeleteVersion(c *gin.Context) {
	packageIdString := c.Param("id")
	versionIdString := c.Param("versionId")
	packageId, err := uuid.Parse(packageIdString)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid package id"})
		return
	}

	versionId, err := strconv.ParseUint(versionIdString, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid version id"})
		return
	}

	version := models.PackageVersion[any]{}
	err = db.DB().Model(&version).Delete(`"package_id" = ? AND id = ?`, packageId, versionId).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, version)
}
