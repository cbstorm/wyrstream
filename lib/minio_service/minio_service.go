package minio_service

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOOption struct {
	ctx context.Context
}

type MinIOOptionFunc func(*MinIOOption)

func WithContext(c context.Context) MinIOOptionFunc {
	return func(mi *MinIOOption) {
		mi.ctx = c
	}
}

type IMinioConfig interface {
	LoadMinioConfig() error
	MINIO_HOST() string
	MINIO_PORT() uint16
	MINIO_ACCESS_KEY() string
	MINIO_SECRET_KEY() string
	MINIO_BUCKET_NAME() string
	MINIO_PUBLIC_URL() string
}

var minio_instance *MinIOService
var minio_instance_sync sync.Once

func GetMinioService() *MinIOService {
	if minio_instance == nil {
		minio_instance_sync.Do(func() {
			minio_instance = &MinIOService{
				logger: logger.NewLogger("MINIO_SERVICE"),
			}
		})
	}
	return minio_instance
}

type MinIOService struct {
	client                 *minio.Client
	bucket_name            string
	config                 IMinioConfig
	logger                 *logger.Logger
	health_check_cancel_fn context.CancelFunc
}

func (m *MinIOService) LoadConfig(config IMinioConfig) error {
	if err := config.LoadMinioConfig(); err != nil {
		return err
	}
	m.config = config
	return nil
}

func (m *MinIOService) Init() error {
	client, err := minio.New(fmt.Sprintf("%s:%d", m.config.MINIO_HOST(), m.config.MINIO_PORT()), &minio.Options{
		Creds:  credentials.NewStaticV4(m.config.MINIO_ACCESS_KEY(), m.config.MINIO_SECRET_KEY(), ""),
		Secure: false,
	})
	if err != nil {
		return err
	}
	m.client = client
	m.bucket_name = m.config.MINIO_BUCKET_NAME()
	health_check_cancel_fn, err := m.client.HealthCheck(5 * time.Second)
	if err != nil {
		return err
	}
	m.health_check_cancel_fn = health_check_cancel_fn
	endpoint := m.client.EndpointURL()
	m.logger.Info("Initialized the Minio client successfully at %v", endpoint.String())
	return nil
}

func (m *MinIOService) AssertBucket() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := m.client.MakeBucket(ctx, m.bucket_name, minio.MakeBucketOptions{}); err != nil {
		_, existed_err := m.client.BucketExists(ctx, m.bucket_name)
		if existed_err != nil {
			return err
		}
	}
	m.logger.Info("Assert bucket %s successfully", m.bucket_name)
	return nil
}

func (m *MinIOService) FPutObject(object MinIOFObject, opts ...MinIOOptionFunc) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	mio := &MinIOOption{}
	for _, f := range opts {
		f(mio)
	}
	if mio.ctx != nil {
		cancel()
		ctx = mio.ctx
	}
	po := minio.PutObjectOptions{}
	if object.ContentType() != "" {
		po.ContentType = object.ContentType()
	}
	result, err := m.client.FPutObject(ctx, m.bucket_name, object.ObjectName(), object.FilePath(), po)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", m.config.MINIO_PUBLIC_URL(), result.Bucket, result.Key), nil
}

func (m *MinIOService) FPutObjects(objects *[]MinIOFObject) *[]*BulkPutObjectResult {
	no_obj := len(*objects)
	t := time.Duration(30 * no_obj)
	ctx, cancel := context.WithTimeout(context.Background(), t*time.Second)
	defer cancel()
	upload_ch := make(chan *BulkPutObjectResult, no_obj)
	defer close(upload_ch)
	for _, e := range *objects {
		go func(o MinIOFObject) {
			a, err := m.FPutObject(o, WithContext(ctx))
			upload_ch <- &BulkPutObjectResult{ObjectName: o.ObjectName(), PutResult: a, Error: err}
		}(e)
	}
	result := make([]*BulkPutObjectResult, 0)
	for i := 0; i < no_obj; i++ {
		r := <-upload_ch
		result = append(result, r)
	}
	return &result
}

func (m *MinIOService) PutObject(object MinIOObject, opts ...MinIOOptionFunc) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	mio := &MinIOOption{}
	for _, f := range opts {
		f(mio)
	}
	if mio.ctx != nil {
		cancel()
		ctx = mio.ctx
	}
	put_opts := minio.PutObjectOptions{}
	if object.ContentType() != "" {
		put_opts.ContentType = object.ContentType()
	}
	result, err := m.client.PutObject(ctx, m.bucket_name, object.ObjectName(), bytes.NewReader(object.Data()), object.Size(), put_opts)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", m.config.MINIO_PUBLIC_URL(), result.Bucket, result.Key), nil
}

func (m *MinIOService) PutObjects(objects *[]MinIOObject) (*[]*BulkPutObjectResult, error) {
	t := time.Duration(30 * len(*objects))
	ctx, cancel := context.WithTimeout(context.Background(), t*time.Second)
	defer cancel()
	upload_ch := make(chan *BulkPutObjectResult, len(*objects))
	defer close(upload_ch)
	for _, e := range *objects {
		go func(o MinIOObject) {
			a, err := m.PutObject(o, WithContext(ctx))
			upload_ch <- &BulkPutObjectResult{ObjectName: o.ObjectName(), PutResult: a, Error: err}
		}(e)
	}
	result := make([]*BulkPutObjectResult, 0)
	for r := range upload_ch {
		result = append(result, r)
	}
	return &result, nil
}

func (m *MinIOService) CancelHealthCheck() {
	m.health_check_cancel_fn()
}
