package drivers

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cnb.cool/liey/liey-go/contracts"
)

type Local struct {
	root      string
	urlPrefix string
}

type LocalConfig struct {
	Root      string `json:"root" yaml:"root" mapstructure:"root"`
	UrlPrefix string `json:"url_prefix" yaml:"url_prefix" mapstructure:"url_prefix"`
}

func NewLocal(config *LocalConfig) (*Local, error) {
	if config.Root == "" {
		return nil, fmt.Errorf("local storage root path is required")
	}

	if err := os.MkdirAll(config.Root, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &Local{
		root:      config.Root,
		urlPrefix: config.UrlPrefix,
	}, nil
}

func (l *Local) Put(ctx context.Context, path string, content io.Reader, options ...contracts.StorageOption) error {
	fullPath := filepath.Join(l.root, path)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (l *Local) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(l.root, path)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

func (l *Local) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(l.root, path)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (l *Local) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(l.root, path)

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (l *Local) Size(ctx context.Context, path string) (int64, error) {
	fullPath := filepath.Join(l.root, path)

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, fmt.Errorf("file not found: %s", path)
		}
		return 0, fmt.Errorf("failed to stat file: %w", err)
	}

	return info.Size(), nil
}

func (l *Local) Url(ctx context.Context, path string) string {
	if l.urlPrefix == "" {
		return path
	}
	return l.urlPrefix + "/" + strings.TrimPrefix(path, "/")
}

func (l *Local) PreSign(ctx context.Context, path string, expire time.Duration, options ...contracts.StorageOption) (string, error) {
	return "", fmt.Errorf("local storage does not support pre-sign urls")
}

func (l *Local) DriverName() string {
	return "local"
}
