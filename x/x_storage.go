package x

import "github.com/xframe-go/x/storage"

func RegisterStorage(fn func() storage.Config) {
	cfg := fn()

	if cfg.Disks == nil {
		return
	}

	manager, err := storage.NewManager(&cfg)
	if err != nil {
		panic(err)
	}

	rocket.storage = manager
}

func Storage() *storage.Manager {
	return rocket.storage
}
