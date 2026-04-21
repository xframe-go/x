package contracts

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	DriverName() string
	Put(ctx context.Context, path string, content io.Reader, options ...StorageOption) error
	Get(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
	Exists(ctx context.Context, path string) (bool, error)
	Size(ctx context.Context, path string) (int64, error)
	Url(ctx context.Context, path string) string
	PreSign(ctx context.Context, path string, expire time.Duration, options ...StorageOption) (string, error)
}

type StorageOption func(*StorageOptions)

type StorageOptions struct {
	ContentType string
	Visibility  string
}

func WithContentType(contentType string) StorageOption {
	return func(o *StorageOptions) {
		o.ContentType = contentType
	}
}

func WithVisibility(visibility string) StorageOption {
	return func(o *StorageOptions) {
		o.Visibility = visibility
	}
}
