package pypi

import (
	"github.com/alin-io/pkgproxy/models"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"io"
	"log"
	"mime/multipart"
	"slices"
)

func (s *Service) UploadHandler(c *gin.Context) {
	pkgName := c.PostForm("name")
	pkgVersionName := c.PostForm("version")
	file, err := c.FormFile("content")
	if err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}
	fileHandle, err := file.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	defer func(fileHandle multipart.File) {
		err := fileHandle.Close()
		if err != nil {
			log.Println(err)
		}
	}(fileHandle)

	checksum, size, err := s.ChecksumReader(fileHandle)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to Upload Package"})
		return
	}
	_, err = fileHandle.Seek(0, io.SeekStart)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to Upload Package"})
		return
	}
	filenamePostfix := s.FilenamePostfix(file.Filename, pkgName, pkgVersionName)
	storageFilename := s.PackageFilename(checksum, filenamePostfix)

	packageModel := models.Package[PackageMetadata]{}
	pkgVersion := models.PackageVersion[PackageMetadata]{}
	_ = packageModel.FillByName(pkgName, s.Prefix)
	if packageModel.Id > 0 {
		pkgVersion, err = packageModel.Version(pkgVersionName)
		if err != nil {
			log.Println("Unable to fill package versions: ", err)
			c.JSON(500, gin.H{"error": "Unable to Upload Package"})
			return
		}
	}

	err = s.Storage.WriteFile(storageFilename, nil, fileHandle)

	if packageModel.Id > 0 && len(pkgVersion.Digest) > 0 {
		if slices.Contains(pkgVersion.Metadata.Data().OriginalFiles, file.Filename) {
			c.JSON(200, pkgVersion)
			return
		} else {
			versionMeta := pkgVersion.Metadata.Data()
			versionMeta.OriginalFiles = append(versionMeta.OriginalFiles, file.Filename)
			pkgVersion.Metadata = datatypes.NewJSONType(versionMeta)
			err = pkgVersion.SaveMeta()
			if err != nil {
				log.Println("Unable to update package version metadata: ", err)
			}
		}
	} else {
		pkgVersion = models.PackageVersion[PackageMetadata]{
			Digest:  checksum,
			Version: pkgVersionName,
			Size:    uint64(size),
			Metadata: datatypes.NewJSONType(PackageMetadata{
				RequiresPython: c.PostForm("requires_python"),
				OriginalFiles:  []string{file.Filename},
			}),
		}

		packageModel = models.Package[PackageMetadata]{
			Name:      pkgName,
			Service:   s.Prefix,
			Namespace: "",
			AuthId:    c.GetString("token"),
			Versions: []models.PackageVersion[PackageMetadata]{
				pkgVersion,
			},
		}

		err = packageModel.Insert()
		if err != nil {
			log.Println("Unable to create package in DB: ", err)
		}
	}

	if err != nil {
		err = s.Storage.DeleteFile(storageFilename)
		if err != nil {
			log.Println("Unable to Delete/Rollback package upload: ", err)
		}
		c.JSON(500, gin.H{"error": "Unable to Upload Package"})
		return
	}
	c.JSON(200, pkgVersion.Metadata.Data())
}
