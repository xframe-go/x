package drivers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/xframe-go/x/contracts"
)

type S3 struct {
	client    *s3.Client
	bucket    string
	urlPrefix string
}

type S3Config struct {
	Endpoint  string `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
	AccessKey string `json:"access_key" yaml:"access_key" mapstructure:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key" mapstructure:"secret_key"`
	Bucket    string `json:"bucket" yaml:"bucket" mapstructure:"bucket"`
	Region    string `json:"region" yaml:"region" mapstructure:"region"`
	UseSSL    bool   `json:"use_ssl" yaml:"use_ssl" mapstructure:"use_ssl"`
	UrlPrefix string `json:"url_prefix" yaml:"url_prefix" mapstructure:"url_prefix"`
}

func NewS3(cfg *S3Config) (*S3, error) {
	if cfg.Endpoint == "" || cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("s3 configuration is incomplete")
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if cfg.Endpoint != "" {
			return aws.Endpoint{
				URL:           cfg.Endpoint,
				SigningRegion: cfg.Region,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	ctx := context.Background()
	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(cfg.Bucket),
	})
	if err != nil {
		_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(cfg.Bucket),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(cfg.Region),
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &S3{
		client:    client,
		bucket:    cfg.Bucket,
		urlPrefix: cfg.UrlPrefix,
	}, nil
}

func (s *S3) Put(ctx context.Context, path string, content io.Reader, options ...contracts.StorageOption) error {
	opts := &contracts.StorageOptions{}
	for _, opt := range options {
		opt(opts)
	}

	buf := new(bytes.Buffer)
	size, err := buf.ReadFrom(content)
	if err != nil {
		return fmt.Errorf("failed to read content: %w", err)
	}

	putInput := &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(path),
		Body:          buf,
		ContentLength: aws.Int64(size),
	}

	if opts.ContentType != "" {
		putInput.ContentType = aws.String(opts.ContentType)
	}

	if opts.Visibility == "public" {
		putInput.ACL = types.ObjectCannedACLPublicRead
	}

	_, err = s.client.PutObject(ctx, putInput)
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}

	return nil
}

func (s *S3) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return output.Body, nil
}

func (s *S3) Delete(ctx context.Context, path string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

func (s *S3) Exists(ctx context.Context, path string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		var notFound *types.NotFound
		if ok := fmt.Sprintf("%T", err) == "*types.NotFound"; ok || notFound != nil {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *S3) Size(ctx context.Context, path string) (int64, error) {
	output, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to stat object: %w", err)
	}

	return aws.ToInt64(output.ContentLength), nil
}

func (s *S3) Url(ctx context.Context, path string) string {
	if s.urlPrefix == "" {
		return fmt.Sprintf("s3://%s/%s", s.bucket, path)
	}
	return s.urlPrefix + "/" + strings.TrimPrefix(path, "/")
}

func (s *S3) PreSign(ctx context.Context, path string, expire time.Duration, options ...contracts.StorageOption) (string, error) {
	opts := &contracts.StorageOptions{}
	for _, opt := range options {
		opt(opts)
	}

	presignClient := s3.NewPresignClient(s.client)

	putInput := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	}

	if opts.ContentType != "" {
		putInput.ContentType = aws.String(opts.ContentType)
	}

	if opts.Visibility == "public" {
		putInput.ACL = types.ObjectCannedACLPublicRead
	}

	req, err := presignClient.PresignPutObject(ctx, putInput, s3.WithPresignExpires(expire))
	if err != nil {
		return "", fmt.Errorf("failed to presign url: %w", err)
	}

	return req.URL, nil
}

func (s *S3) DriverName() string {
	return "s3"
}
