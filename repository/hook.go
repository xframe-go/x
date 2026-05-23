package repository

import (
	"context"

	"gorm.io/gorm"
)

type UpdatedHook[M any, K comparable] func(ctx context.Context, tx *gorm.DB, key K, m M) error

func WithUpdatedHook[M any, K comparable](hook UpdatedHook[M, K]) OptionFn[M, K] {
	return func(c *options[M, K]) {
		c.UpdatedHook = hook
	}
}

type CreatedHook[M any, K comparable] func(ctx context.Context, tx *gorm.DB, m *M) error

func WithCreatedHook[M any, K comparable](hook CreatedHook[M, K]) OptionFn[M, K] {
	return func(c *options[M, K]) {
		c.CreatedHook = hook
	}
}

type DeletedHook[M any, K comparable] func(ctx context.Context, tx *gorm.DB, deleted []M, keys ...K) error

func WithDeletedHook[M any, K comparable](hook DeletedHook[M, K]) OptionFn[M, K] {
	return func(c *options[M, K]) {
		c.DeletedHook = hook
	}
}
