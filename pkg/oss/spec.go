package oss

// Provider todo
type Provider interface {
	UploadFile(bucketName, objectKey, localFilePath string) error
}
