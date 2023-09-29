package container

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/alin-io/pkgproxy/models"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"io"
)

// StartLayerUploadHandler POST /v2/<name>/blobs/uploads/
func (s *Service) StartLayerUploadHandler(c *gin.Context) {
	pkgName := s.ConstructFullPkgName(c)
	uploadItem := models.TmpUploadStore{Name: pkgName, UploadRange: "0-0"}
	err := uploadItem.Insert()
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to Start the upload process"})
		return
	}

	c.Header("Location", "/v2/"+pkgName+"/blobs/uploads/"+uploadItem.Id)
	c.Header("Docker-Upload-UUID", uploadItem.Id)
	c.Header("Range", "bytes="+uploadItem.UploadRange)
	c.Header("Content-Length", "0")
	c.Status(202)
	c.Done()
}

func (s *Service) GetUploadProgressHandler(c *gin.Context) {
	pkgName := s.ConstructFullPkgName(c)
	uploadUUID := c.Param("uuid")
	uploadItem := models.TmpUploadStore{}
	err := uploadItem.FillById(uploadUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to get upload progress"})
		return
	}
	if uploadItem.Name != pkgName || uploadItem.Id != uploadUUID {
		c.JSON(404, gin.H{"error": "Upload not found"})
		return
	}

	c.Header("Range", "bytes="+uploadItem.UploadRange)
	c.Header("Location", "/v2/"+pkgName+"/blobs/uploads/"+uploadUUID)
	c.Header("Docker-Upload-UUID", uploadUUID)
	c.Status(204)
	c.Done()
}

func (s *Service) ChunkUploadHandler(c *gin.Context) {
	pkgName := s.ConstructFullPkgName(c)
	uploadUUID := c.Param("uuid")
	uploadItem := models.TmpUploadStore{}
	err := uploadItem.FillById(uploadUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to get upload progress"})
		return
	}
	if uploadItem.Name != pkgName || uploadItem.Id != uploadUUID {
		c.JSON(404, gin.H{"error": "Upload not found"})
		return
	}

	_, _, err = s.appendStorageData(uploadUUID, c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to save chunk"})
		return
	}

	uploadItem.UploadRange = c.GetHeader("Content-Range")
	err = uploadItem.Update()
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to save chunk metadata"})
		return
	}

	c.Header("Location", "/v2/"+pkgName+"/blobs/uploads/"+uploadItem.Id)
	c.Header("Docker-Upload-UUID", uploadItem.Id)
	c.Header("Range", "bytes="+uploadItem.UploadRange)
	c.Header("Content-Length", "0")
	c.Status(202)
}

func (s *Service) UploadHandler(c *gin.Context) {
	pkgName := s.ConstructFullPkgName(c)
	inputDigest := c.Query("digest")
	uploadUUID := c.Param("uuid")
	uploadItem := models.TmpUploadStore{}
	err := uploadItem.FillById(uploadUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to get upload progress"})
		return
	}
	if uploadItem.Name != pkgName || uploadItem.Id != uploadUUID {
		c.JSON(404, gin.H{"error": "Upload not found"})
		return
	}

	digest, totalSize, err := s.appendStorageData(uploadUUID, c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to save chunk"})
		return
	}

	if inputDigest != "" && inputDigest != digest {
		c.JSON(400, gin.H{"error": "Digest mismatch"})
		return
	}

	uploadItem.Digest = digest
	uploadItem.Size = totalSize
	err = uploadItem.Update()
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to save chunk metadata"})
		return
	}

	c.Header("Location", "/v2/"+pkgName+"/blobs/"+digest)
	c.Header("Content-Length", "0")
	c.Header("Docker-Content-Digest", digest)
	c.Status(201)
	c.Done()
}

func (s *Service) ManifestUploadHandler(c *gin.Context) {
	metadata := PackageMetadata{}
	err := c.ShouldBind(&metadata)
	if err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	if len(metadata.FsLayers) == 0 {
		c.JSON(400, gin.H{"error": "Bad Request: No layers found"})
		return
	}
	tagLayer := metadata.FsLayers[len(metadata.FsLayers)-1]
	if tagLayer.BlobSum == "" {
		c.JSON(400, gin.H{"error": "Bad Request: No tag layer found"})
		return
	}

	uploadItem := models.TmpUploadStore{}
	err = uploadItem.FillByDigest(tagLayer.BlobSum, metadata.Name)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to get upload progress"})
		return
	}

	pkg := models.Package[PackageMetadata]{}
	err = pkg.FillByName(metadata.Name, s.Prefix)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to check the DB for package"})
		return
	}
	pkgVersion := models.PackageVersion[PackageMetadata]{
		Version:  metadata.Tag,
		Digest:   tagLayer.BlobSum,
		Metadata: datatypes.NewJSONType[PackageMetadata](metadata),
		Size:     uploadItem.Size,
		Tag:      metadata.Tag,
	}
	if pkg.Id == 0 {
		pkg = models.Package[PackageMetadata]{
			Name:      metadata.Name,
			Service:   s.Prefix,
			Namespace: "",
			AuthId:    c.GetString("token"),
			Versions: []models.PackageVersion[PackageMetadata]{
				pkgVersion,
			},
		}
		err = pkg.Insert()
		if err != nil {
			c.JSON(500, gin.H{"error": "Unable to insert package"})
			return
		}
	} else {
		pkgVersion.PackageId = pkg.Id
		err = pkgVersion.SaveMeta()
		if err != nil {
			c.JSON(500, gin.H{"error": "Unable to insert package version"})
			return
		}
	}
	c.Status(202)
	c.Done()
}

type partialReadWriter struct {
	io.Reader
	input     io.Reader
	output    io.Writer
	totalSize uint64
}

func (p partialReadWriter) Read(b []byte) (n int, err error) {
	n, err = p.input.Read(b)
	if err != nil {
		return n, err
	}
	_, _ = p.output.Write(b[:n])
	p.totalSize += uint64(n)
	return n, nil
}

func (s *Service) appendStorageData(uploadUUID string, input io.Reader) (digest string, size uint64, err error) {
	fileReader, err := s.Storage.GetFile(fmt.Sprintf("%s/%s", s.Prefix, uploadUUID))
	if err != nil {
		return "", 0, err
	}

	if fileReader == nil {
		fileReader = io.NopCloser(bytes.NewReader([]byte{}))
	}

	hasher := sha256.New()

	rw := partialReadWriter{
		input:  io.MultiReader(fileReader, input),
		output: hasher,
	}

	err = s.Storage.WriteFile(fmt.Sprintf("%s/%s", s.Prefix, uploadUUID), nil, rw)
	if err != nil {
		return "", 0, err
	}

	return hex.EncodeToString(hasher.Sum(nil)), rw.totalSize, nil
}
