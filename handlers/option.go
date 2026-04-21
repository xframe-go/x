package handlers

import (
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type (
	ParseHook[M any, C any, U any, K comparable] func(ctx *Context) (M, error)

	Hook[M any, C any, U any, K comparable] func(ctx *Context, tx *gorm.DB, primaryKey []K) error

	BeforeCreateHook[M any, C any, U any, K comparable] func(ctx *Context, tx *gorm.DB, c *C, m *M) error

	AfterCreatedHook[M any, C any, U any, K comparable] func(ctx *Context, tx *gorm.DB, c *C, m *M) error

	BeforeUpdateHook[M any, C any, U any, K comparable] func(ctx *Context, tx *gorm.DB, u *U, primaryKey K) error

	AfterUpdatedHook[M any, C any, U any, K comparable] func(ctx *Context, tx *gorm.DB, u *U, primaryKey K) error

	PrimaryKeyConverter[K comparable] func(pk any) K
)

type Option[M any, C any, U any, K comparable] struct {
	beforeCreate        BeforeCreateHook[M, C, U, K]
	afterCreated        AfterCreatedHook[M, C, U, K]
	beforeUpdate        BeforeUpdateHook[M, C, U, K]
	afterUpdated        AfterUpdatedHook[M, C, U, K]
	beforeDestroy       Hook[M, C, U, K]
	afterDestroyed      Hook[M, C, U, K]
	primaryKeyName      string
	primaryKeyConverter PrimaryKeyConverter[K]
}

type WithOption[M any, C any, U any, K comparable] func(opt *Option[M, C, U, K])

func defaultOption[M any, C any, U any, K comparable]() *Option[M, C, U, K] {
	return &Option[M, C, U, K]{
		beforeDestroy: func(ctx *Context, tx *gorm.DB, primaryKey []K) error {
			return nil
		},
		primaryKeyName: "id",
		primaryKeyConverter: func(pk any) K {
			var k any
			k = cast.ToUint64(pk)
			return k.(K)
		},
	}
}

func BeforeCreate[M any, C any, U any, K comparable](hook BeforeCreateHook[M, C, U, K]) WithOption[M, C, U, K] {
	return func(opt *Option[M, C, U, K]) {
		opt.beforeCreate = hook
	}
}

func AfterCreated[M any, C any, U any, K comparable](hook AfterCreatedHook[M, C, U, K]) WithOption[M, C, U, K] {
	return func(opt *Option[M, C, U, K]) {
		opt.afterCreated = hook
	}
}

func BeforeUpdate[M any, C any, U any, K comparable](hook BeforeUpdateHook[M, C, U, K]) WithOption[M, C, U, K] {
	return func(opt *Option[M, C, U, K]) {
		opt.beforeUpdate = hook
	}
}

func AfterUpdated[M any, C any, U any, K comparable](hook AfterUpdatedHook[M, C, U, K]) WithOption[M, C, U, K] {
	return func(opt *Option[M, C, U, K]) {
		opt.afterUpdated = hook
	}
}

func BeforeDestroy[M any, C any, U any, K comparable](hook Hook[M, C, U, K]) WithOption[M, C, U, K] {
	return func(opt *Option[M, C, U, K]) {
		opt.beforeDestroy = hook
	}
}

func AfterDestroyed[M any, C any, U any, K comparable](hook Hook[M, C, U, K]) WithOption[M, C, U, K] {
	return func(opt *Option[M, C, U, K]) {
		opt.afterDestroyed = hook
	}
}
