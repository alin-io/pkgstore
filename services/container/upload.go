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
	"log"
	"strings"
)

// StartLayerUploadHandler POST /v2/<name>/blobs/uploads/
func (s *Service) StartLayerUploadHandler(c *gin.Context) {
	pkgName := s.ConstructFullPkgName(c)
	asset := models.Asset{}
	err := asset.StartUpload()
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to Start the upload process"})
		return
	}

	c.Header("Location", "/v2/"+pkgName+"/blobs/uploads/"+asset.UploadUUID)
	c.Header("Docker-Upload-UUID", asset.UploadUUID)
	c.Header("Range", "bytes="+asset.UploadRange)
	c.Header("Content-Length", "0")
	c.Status(202)
	c.Done()
}

func (s *Service) CheckBlobExistenceHandler(c *gin.Context) {
	digest := strings.Replace(c.Param("sha256"), "sha256:", "", 1)
	asset := models.Asset{}
	err := asset.FillByDigest(digest)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to check the DB for package version"})
		return
	}

	if asset.Id == 0 {
		c.JSON(404, gin.H{"error": "Blob not found"})
		return
	}

	c.Header("Docker-Content-Digest", "sha256:"+digest)
	c.Header("Content-Length", fmt.Sprintf("%d", asset.Size))
	c.Status(200)
	c.Done()
}

func (s *Service) GetUploadProgressHandler(c *gin.Context) {
	pkgName := s.ConstructFullPkgName(c)
	uploadUUID := c.Param("uuid")
	asset := models.Asset{}
	err := asset.FillByUploadUUID(uploadUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to get upload progress"})
		return
	}
	if asset.UploadUUID != uploadUUID {
		c.JSON(404, gin.H{"error": "Upload not found"})
		return
	}

	c.Header("Range", asset.UploadRange)
	c.Header("Location", "/v2/"+pkgName+"/blobs/uploads/"+uploadUUID)
	c.Header("Docker-Upload-UUID", uploadUUID)
	c.Status(204)
	c.Done()
}

func (s *Service) ChunkUploadHandler(c *gin.Context) {
	pkgName := s.ConstructFullPkgName(c)
	uploadUUID := c.Param("uuid")
	asset := models.Asset{}
	err := asset.FillByUploadUUID(uploadUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to get upload progress"})
		return
	}
	if asset.UploadUUID != uploadUUID {
		c.JSON(404, gin.H{"error": "Upload not found"})
		return
	}

	_, chunkSize, err := s.appendStorageData(uploadUUID, c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to save chunk"})
		return
	}

	if chunkSize > 0 {
		asset.UploadRange = fmt.Sprintf("%d-%d", asset.Size, chunkSize)
		asset.Size = chunkSize
	}
	err = asset.Update()
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to save chunk metadata"})
		return
	}

	c.Header("Location", "/v2/"+pkgName+"/blobs/uploads/"+uploadUUID)
	c.Header("Docker-Upload-UUID", uploadUUID)
	c.Header("Range", asset.UploadRange)
	c.Header("Content-Length", "0")
	c.Status(204)
	c.Done()
}

func (s *Service) UploadHandler(c *gin.Context) {
	pkgName := s.ConstructFullPkgName(c)
	inputDigest := strings.Replace(c.Query("digest"), "sha256:", "", 1)
	uploadUUID := c.Param("uuid")
	asset := models.Asset{}
	err := asset.FillByUploadUUID(uploadUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to get upload progress"})
		return
	}
	if asset.UploadUUID != uploadUUID {
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

	err = s.Storage.CopyFile(fmt.Sprintf("%s/%s", s.Prefix, uploadUUID), fmt.Sprintf("%s/%s", s.Prefix, digest))
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to store the file"})
		return
	}

	asset.Digest = digest
	asset.Size = totalSize
	err = asset.Update()
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to save chunk metadata"})
		return
	}

	err = s.Storage.DeleteFile(fmt.Sprintf("%s/%s", s.Prefix, uploadUUID))
	if err != nil {
		log.Println(err)
	}

	c.Header("Location", "/v2/"+pkgName+"/blobs/"+digest)
	c.Header("Content-Range", "0-"+fmt.Sprintf("%d", totalSize))
	c.Header("Content-Length", "0")
	c.Header("Docker-Content-Digest", "sha256:"+digest)
	c.Status(204)
	c.Done()
}

func (s *Service) ManifestUploadHandler(c *gin.Context) {
	metadata := PackageMetadata{
		ContentType: c.Request.Header.Get("Content-Type"),
	}
	var (
		tagName = c.Param("reference")
		pkgName = s.ConstructFullPkgName(c)
		digest  string
		err     error
	)
	switch metadata.ContentType {
	case ManifestV1ContentType:
		err = c.ShouldBindJSON(&metadata.ManifestV1)
		tagName = metadata.ManifestV1.Tag
		pkgName = metadata.ManifestV1.Name
		if len(metadata.ManifestV1.FsLayers) == 0 {
			c.JSON(400, gin.H{"error": "Bad Request: No layers found"})
			return
		}
		digest = metadata.ManifestV1.FsLayers[len(metadata.ManifestV1.FsLayers)-1].BlobSum
	case ManifestV2ContentType:
		err = c.ShouldBindJSON(&metadata.ManifestV2)
		digest = metadata.ManifestV2.Config.Digest
	case ManifestListV2ContentType:
		err = c.ShouldBindJSON(&metadata.ManifestListV2)
	default:
		c.JSON(400, gin.H{"error": "Bad Request"})
	}

	digest = strings.Replace(digest, "sha256:", "", 1)

	asset := models.Asset{}
	err = asset.FillByDigest(digest)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to get upload progress"})
		return
	}
	if asset.Digest != digest {
		c.JSON(404, gin.H{"error": "Uploaded asset not found"})
		return
	}

	pkg := models.Package[PackageMetadata]{}
	err = pkg.FillByName(pkgName, s.Prefix)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to check the DB for package"})
		return
	}

	if pkg.Id == 0 {
		pkg = models.Package[PackageMetadata]{
			Name:    pkgName,
			Service: s.Prefix,
			AuthId:  c.GetString("token"),
		}
		err = pkg.Insert()
		if err != nil {
			log.Println(err)
		}
	}

	pkgVersion := models.PackageVersion[PackageMetadata]{}
	err = pkgVersion.FillByDigest(digest)
	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to check the DB for package version"})
		return
	}
	if pkgVersion.Id == 0 {
		pkgVersion = models.PackageVersion[PackageMetadata]{
			PackageId: pkg.Id,
			Service:   s.Prefix,
			Digest:    digest,
			Version:   tagName,
			Tag:       tagName,
			Metadata:  datatypes.NewJSONType[PackageMetadata](metadata),
		}
		err = pkgVersion.Save()
	} else {
		if pkgVersion.PackageId != pkg.Id || pkgVersion.Service != s.Prefix {
			c.JSON(404, gin.H{"error": "Package version not found"})
			return
		}
		pkgVersion.Version = tagName
		pkgVersion.Tag = tagName
		pkgVersion.Metadata = datatypes.NewJSONType[PackageMetadata](metadata)
		err = pkgVersion.Save()
	}

	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to insert package version"})
		return
	}

	c.Header("Docker-Content-Digest", "sha256:"+digest)
	c.Status(201)
	c.Done()
}

type sizeHandler struct {
	size uint64
}

type partialReadWriter struct {
	io.Reader
	input     io.Reader
	output    io.Writer
	totalSize *sizeHandler
}

func (p partialReadWriter) Read(b []byte) (n int, err error) {
	n, err = p.input.Read(b)
	p.totalSize.size += uint64(n)
	if err != nil {
		return n, err
	}
	_, _ = p.output.Write(b[:n])
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

	defer func() {
		_ = fileReader.Close()
	}()

	hasher := sha256.New()

	sh := &sizeHandler{}
	rw := partialReadWriter{
		input:     io.MultiReader(fileReader, input),
		output:    hasher,
		totalSize: sh,
	}

	err = s.Storage.WriteFile(fmt.Sprintf("%s/%s", s.Prefix, uploadUUID), nil, rw)
	if err != nil {
		return "", 0, err
	}

	return hex.EncodeToString(hasher.Sum(nil)), sh.size, nil
}
