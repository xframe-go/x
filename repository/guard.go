package repository

import (
	"context"

	"gorm.io/gorm"
)

type GuardFunc[M any] func(ctx context.Context, tx *gorm.DB, m *M) error

func WithGuard[M any, K comparable](guard GuardFunc[M]) OptionFn[M, K] {
	return func(o *options[M, K]) {
		o.guardFunc = guard
	}
}
