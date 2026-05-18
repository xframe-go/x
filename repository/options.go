package repository

import (
	"gorm.io/gorm"
)

type options[M any, C ToModel[M], U ToModel[M], K comparable] struct {
	keywordExpression func(tx *gorm.DB, keyword string) *gorm.DB

	UpdatedHook UpdatedHook[M, C, U, K]

	CreatedHook CreatedHook[M, C, U, K]

	DeletedHook DeletedHook[M, C, U, K]
}

type OptionFn[M any, C ToModel[M], U ToModel[M], K comparable] func(*options[M, C, U, K])

func WithKeywordExpression[M any, C ToModel[M], U ToModel[M], K comparable](
	keywordExp func(tx *gorm.DB, keyword string) *gorm.DB) OptionFn[M, C, U, K] {
	return func(opts *options[M, C, U, K]) {
		opts.keywordExpression = keywordExp
	}
}
