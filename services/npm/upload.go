package npm

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"github.com/alin-io/pkgproxy/db"
	"github.com/alin-io/pkgproxy/models"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type npmUploadRequestBody struct {
	Attachments map[string]struct {
		ContentType string `json:"content_type"`
		Data        string `json:"data"`
		Length      int    `json:"length"`
	} `json:"_attachments"`
	Id          string                                   `json:"_id"`
	Description string                                   `json:"description"`
	Name        string                                   `json:"name"`
	Readme      string                                   `json:"readme"`
	DistTags    map[string]string                        `json:"dist-tags"`
	Versions    map[string]models.PackageVersionMetadata `json:"versions"`
}

func (s *Service) UploadHandler(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	requestBody := npmUploadRequestBody{}
	err := c.ShouldBind(&requestBody)
	if err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	decodedBytes := make([]byte, 0)
	for _, attachment := range requestBody.Attachments {
		decodedBytes, err = base64.StdEncoding.DecodeString(attachment.Data)
		if err != nil {
			c.JSON(500, gin.H{"error": "Unable to Upload Package"})
			return
		}
		break
	}

	hasher := sha1.New()
	hasher.Write(decodedBytes)
	checksum := hex.EncodeToString(hasher.Sum(nil))

	currentVersion := ""
	var pkgVersion models.PackageVersion
	for _, versionInfo := range requestBody.Versions {
		currentVersion = versionInfo.Version

		pkgVersion = models.PackageVersion{
			Version:  currentVersion,
			Digest:   checksum,
			Metadata: datatypes.NewJSONType(versionInfo),
			Size:     uint64(len(decodedBytes)),
		}

		for tagName, tagVersion := range requestBody.DistTags {
			if tagVersion == versionInfo.Version {
				pkgVersion.Tag = tagName
				break
			}
		}
		break
	}

	err = s.storage.WriteFile(s.PackageFilename(checksum), nil, bytes.NewReader(decodedBytes))
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to Upload Package"})
		return
	}

	db.DB().Create(&models.Package{
		Name:      requestBody.Name,
		Namespace: "",
		AuthId:    "",
		Versions:  []models.PackageVersion{pkgVersion},
	})
	c.JSON(200, requestBody)
}
