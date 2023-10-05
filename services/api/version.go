package api

import (
	"github.com/alin-io/pkgstore/db"
	"github.com/alin-io/pkgstore/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (s *Service) ListVersionsHandler(c *gin.Context) {
	packageIdString := c.Param("id")
	packageId, err := strconv.ParseUint(packageIdString, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid package id"})
		return
	}

	versions := make([]models.PackageVersion[any], 0)
	err = db.DB().Model(&versions).Where("package_id = ?", packageId).Find(&versions).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, versions)
}

func (s *Service) DeleteVersion(c *gin.Context) {
	packageIdString := c.Param("id")
	versionIdString := c.Param("versionId")
	packageId, err := strconv.ParseUint(packageIdString, 10, 64)
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
