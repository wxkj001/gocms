package oss

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3 struct {
	client *minio.Client
	err    error
}
type S3Config struct {
	Endpoint string
	// 访问密钥
	AccessKeyID string
	// 访问密钥
	SecretAccessKey string
	Token           string
	// 桶名
	BucketName string
	// 地区
	Location string
	UseSSL   bool
}

func NewS3() *S3 {
	return &S3{}
}

// ProgressReader 是一个用于跟踪读取进度的io.Reader包装器
type ProgressReader struct {
	reader        io.Reader
	totalBytes    int64
	bytesUploaded int64
	progressFunc  func(bytesUploaded int64, totalBytes int64)
}

// Read 实现io.Reader接口
func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	pr.bytesUploaded += int64(n)
	if pr.progressFunc != nil {
		pr.progressFunc(pr.bytesUploaded, pr.totalBytes)
	}
	return
}
func (s *S3) GetConfig() *S3Config {
	return &S3Config{}
}
func (s *S3) Err() error {
	return s.err
}

// GetObjectURL 获取对象的URL
// bucketName: 桶名称
// objectName: 对象名称
// expires: URL过期时间（秒），如果为0则使用默认值7天
// 返回对象的URL和错误
func (s *S3) GetObjectURL(bucketName, objectName string, expires int) (string, error) {
	if s.err != nil {
		return "", s.err
	}

	// 设置过期时间，默认7天
	expiryTime := time.Duration(expires) * time.Second
	if expires == 0 {
		expiryTime = time.Duration(7*24*60*60) * time.Second
	}

	// 获取预签名URL
	presignedURL, err := s.client.PresignedGetObject(context.Background(), bucketName, objectName, expiryTime, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

// 实例化
func (s *S3) Client(config *S3Config) *S3 {
	s.client, s.err = minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, config.Token),
		Secure: config.UseSSL,
	})
	if s.err != nil {
		return s
	}
	found, err := s.client.BucketExists(context.Background(), config.BucketName)
	if err != nil {
		s.err = err
		return s
	}
	// 检查桶是否存在
	if !found {
		// 创建桶
		err = s.client.MakeBucket(context.Background(), config.BucketName, minio.MakeBucketOptions{Region: config.Location})
		if err != nil {
			s.err = err
			return s
		}
	}
	return s
}

// s3 上传
// Upload 上传文件到S3存储
// bucketName: 桶名称
// objectName: 对象名称（存储在S3上的文件名）
// filePath: 本地文件路径
// contentType: 内容类型，如application/octet-stream
// 返回上传信息和错误
func (s *S3) Upload(bucketName, objectName, filePath, contentType string) (info minio.UploadInfo, err error) {
	if s.err != nil {
		return info, s.err
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return info, err
	}
	defer file.Close()

	// 获取文件信息
	fileStat, err := file.Stat()
	if err != nil {
		return info, err
	}

	// 上传文件
	return s.client.PutObject(context.Background(), bucketName, objectName, file, fileStat.Size(), minio.PutObjectOptions{
		ContentType: contentType,
	})
}

// UploadWithReader 使用io.Reader上传数据到S3存储
// bucketName: 桶名称
// objectName: 对象名称（存储在S3上的文件名）
// reader: 数据读取器
// objectSize: 对象大小
// contentType: 内容类型
// 返回上传信息和错误
func (s *S3) UploadWithReader(bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) (info minio.UploadInfo, err error) {
	if s.err != nil {
		return info, s.err
	}

	// 上传数据
	return s.client.PutObject(context.Background(), bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
}

// UploadWithProgress 上传文件到S3存储并跟踪进度
// bucketName: 桶名称
// objectName: 对象名称（存储在S3上的文件名）
// filePath: 本地文件路径
// contentType: 内容类型
// progressFunc: 进度回调函数，接收已上传字节数和总字节数
// 返回上传信息和错误
func (s *S3) UploadWithProgress(bucketName, objectName, filePath, contentType string, progressFunc func(bytesUploaded int64, totalBytes int64)) (info minio.UploadInfo, err error) {
	if s.err != nil {
		return info, s.err
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return info, err
	}
	defer file.Close()

	// 获取文件信息
	fileStat, err := file.Stat()
	if err != nil {
		return info, err
	}

	// 创建一个带进度的reader
	reader := &ProgressReader{
		reader:        file,
		totalBytes:    fileStat.Size(),
		progressFunc:  progressFunc,
		bytesUploaded: 0,
	}

	// 上传文件
	return s.client.PutObject(context.Background(), bucketName, objectName, reader, fileStat.Size(), minio.PutObjectOptions{
		ContentType: contentType,
	})
}
