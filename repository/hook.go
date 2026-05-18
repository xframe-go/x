package repository

import (
	"context"

	"gorm.io/gorm"
)

type UpdatedHook[M any, C ToModel[M], U ToModel[M], K comparable] func(ctx context.Context, tx *gorm.DB, key K, m M) error

func WithUpdatedHook[M any, C ToModel[M], U ToModel[M], K comparable](hook UpdatedHook[M, C, U, K]) OptionFn[M, C, U, K] {
	return func(c *options[M, C, U, K]) {
		c.UpdatedHook = hook
	}
}

type CreatedHook[M any, C ToModel[M], U ToModel[M], K comparable] func(ctx context.Context, tx *gorm.DB, m *M) error

func WithCreatedHook[M any, C ToModel[M], U ToModel[M], K comparable](hook CreatedHook[M, C, U, K]) OptionFn[M, C, U, K] {
	return func(c *options[M, C, U, K]) {
		c.CreatedHook = hook
	}
}

type DeletedHook[M any, C ToModel[M], U ToModel[M], K comparable] func(ctx context.Context, tx *gorm.DB, deleted []M, keys ...K) error

func WithDeletedHook[M any, C ToModel[M], U ToModel[M], K comparable](hook DeletedHook[M, C, U, K]) OptionFn[M, C, U, K] {
	return func(c *options[M, C, U, K]) {
		c.DeletedHook = hook
	}
}
