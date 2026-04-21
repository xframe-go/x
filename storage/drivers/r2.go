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
	"github.com/xframe-go/x/contracts"
)

type R2 struct {
	client    *s3.Client
	bucket    string
	urlPrefix string
}

type R2Config struct {
	AccountId string `json:"account_id" yaml:"account_id" mapstructure:"account_id"`
	AccessKey string `json:"access_key" yaml:"access_key" mapstructure:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key" mapstructure:"secret_key"`
	Bucket    string `json:"bucket" yaml:"bucket" mapstructure:"bucket"`
	UrlPrefix string `json:"url_prefix" yaml:"url_prefix" mapstructure:"url_prefix"`
}

func NewR2(cfg *R2Config) (*R2, error) {
	if cfg.AccountId == "" || cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("r2 configuration is incomplete")
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountId)

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	_, err = client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: aws.String(cfg.Bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("r2 bucket %s does not exist or is not accessible: %w", cfg.Bucket, err)
	}

	return &R2{
		client:    client,
		bucket:    cfg.Bucket,
		urlPrefix: cfg.UrlPrefix,
	}, nil
}

func (r *R2) Put(ctx context.Context, path string, content io.Reader, options ...contracts.StorageOption) error {
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
		Bucket:        aws.String(r.bucket),
		Key:           aws.String(path),
		Body:          buf,
		ContentLength: aws.Int64(size),
	}

	if opts.ContentType != "" {
		putInput.ContentType = aws.String(opts.ContentType)
	}

	_, err = r.client.PutObject(ctx, putInput)
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}

	return nil
}

func (r *R2) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	output, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return output.Body, nil
}

func (r *R2) Delete(ctx context.Context, path string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

func (r *R2) Exists(ctx context.Context, path string) (bool, error) {
	_, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (r *R2) Size(ctx context.Context, path string) (int64, error) {
	output, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to stat object: %w", err)
	}

	return aws.ToInt64(output.ContentLength), nil
}

func (r *R2) Url(ctx context.Context, path string) string {
	if r.urlPrefix == "" {
		return fmt.Sprintf("r2://%s/%s", r.bucket, path)
	}
	return r.urlPrefix + "/" + strings.TrimPrefix(path, "/")
}

func (r *R2) PreSign(ctx context.Context, path string, expire time.Duration, options ...contracts.StorageOption) (string, error) {
	opts := &contracts.StorageOptions{}
	for _, opt := range options {
		opt(opts)
	}

	presignClient := s3.NewPresignClient(r.client)

	putInput := &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(strings.TrimPrefix(path, "/")),
	}

	if opts.ContentType != "" {
		putInput.ContentType = aws.String(opts.ContentType)
	}

	req, err := presignClient.PresignPutObject(ctx, putInput, s3.WithPresignExpires(expire))
	if err != nil {
		return "", fmt.Errorf("failed to presign url: %w", err)
	}

	return req.URL, nil
}

func (r *R2) DriverName() string {
	return "r2"
}
