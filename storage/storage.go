package storage

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"cnb.cool/liey/liey-go/contracts"
	"cnb.cool/liey/liey-go/storage/drivers"
)

type DiskConfig interface {
	DriverName() string
}

type Config struct {
	Default string                `json:"default" yaml:"default" mapstructure:"default"`
	Disks   map[string]DiskConfig `json:"disks" yaml:"disks" mapstructure:"disks"`
}

type Manager struct {
	drivers     map[string]contracts.Storage
	defaultDisk string
	mu          sync.RWMutex
}

func NewManager(config *Config) (*Manager, error) {
	m := &Manager{
		drivers:     make(map[string]contracts.Storage),
		defaultDisk: config.Default,
	}

	for name, disk := range config.Disks {
		if err := m.initDriver(name, disk); err != nil {
			return nil, fmt.Errorf("failed to initialize driver %s: %w", name, err)
		}
	}

	return m, nil
}

func (m *Manager) initDriver(name string, disk DiskConfig) error {
	switch d := disk.(type) {
	case *drivers.LocalDisk:
		driver, err := drivers.NewLocal(&d.Config)
		if err != nil {
			return err
		}
		m.drivers[name] = driver
	case *drivers.S3Disk:
		driver, err := drivers.NewS3(&d.Config)
		if err != nil {
			return err
		}
		m.drivers[name] = driver
	case *drivers.R2Disk:
		driver, err := drivers.NewR2(&d.Config)
		if err != nil {
			return err
		}
		m.drivers[name] = driver
	case *drivers.NilDisk:
		// 未配置的磁盘，跳过初始化
		return nil
	default:
		return fmt.Errorf("unsupported disk type: %T", disk)
	}
	return nil
}

func (m *Manager) Register(name string, driver contracts.Storage) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.drivers[name] = driver
}

func (m *Manager) Disk(name ...string) (contracts.Storage, error) {
	diskName := m.defaultDisk
	if len(name) > 0 && name[0] != "" {
		diskName = name[0]
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	driver, ok := m.drivers[diskName]
	if !ok {
		return nil, fmt.Errorf("disk %s not found", diskName)
	}

	return driver, nil
}

func (m *Manager) Put(ctx context.Context, path string, content io.Reader, options ...contracts.StorageOption) error {
	driver, err := m.Disk()
	if err != nil {
		return err
	}
	return driver.Put(ctx, path, content, options...)
}

func (m *Manager) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	driver, err := m.Disk()
	if err != nil {
		return nil, err
	}
	return driver.Get(ctx, path)
}

func (m *Manager) Delete(ctx context.Context, path string) error {
	driver, err := m.Disk()
	if err != nil {
		return err
	}
	return driver.Delete(ctx, path)
}

func (m *Manager) Exists(ctx context.Context, path string) (bool, error) {
	driver, err := m.Disk()
	if err != nil {
		return false, err
	}
	return driver.Exists(ctx, path)
}

func (m *Manager) Size(ctx context.Context, path string) (int64, error) {
	driver, err := m.Disk()
	if err != nil {
		return 0, err
	}
	return driver.Size(ctx, path)
}

func (m *Manager) Url(ctx context.Context, path string) string {
	driver, err := m.Disk()
	if err != nil {
		return ""
	}
	return driver.Url(ctx, path)
}

func (m *Manager) DriverName() string {
	driver, err := m.Disk()
	if err != nil {
		return ""
	}
	return driver.DriverName()
}

func (m *Manager) PreSign(ctx context.Context, path string, expire time.Duration, options ...contracts.StorageOption) (string, error) {
	driver, err := m.Disk()
	if err != nil {
		return "", err
	}
	return driver.PreSign(ctx, path, expire, options...)
}
