package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type options[M any, K comparable] struct {
	keywordExpression func(tx *gorm.DB, keyword string) *gorm.DB

	conditions func(ctx context.Context) []clause.Expression

	UpdatedHook UpdatedHook[M, K]

	CreatedHook CreatedHook[M, K]

	DeletedHook DeletedHook[M, K]

	guardFunc GuardFunc[M]
}

type OptionFn[M any, K comparable] func(*options[M, K])

func WithKeywordExpression[M any, K comparable](
	keywordExp func(tx *gorm.DB, keyword string) *gorm.DB) OptionFn[M, K] {
	return func(opts *options[M, K]) {
		opts.keywordExpression = keywordExp
	}
}

func WithCustomConditions[M any, K comparable](conditions func(ctx context.Context) []clause.Expression) OptionFn[M, K] {
	return func(opts *options[M, K]) {
		opts.conditions = conditions
	}
}
