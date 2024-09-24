package minio_service

type MinIOFObject interface {
	ObjectName() string
	FilePath() string
	ContentType() string
}

type MinIOObject interface {
	ObjectName() string
	Data() []byte
	Size() int64
	ContentType() string
}

type BulkPutObjectResult struct {
	ObjectName string
	PutResult  string
	Error      error
}

type BulkListDirResult struct {
	Dir    string
	Result *[]string
	Error  error
}
