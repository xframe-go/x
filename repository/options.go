package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type options[M any, K comparable] struct {
	keywordExpression func(tx *gorm.DB, keyword string) *gorm.DB

	conditions []clause.Expression

	UpdatedHook UpdatedHook[M, K]

	CreatedHook CreatedHook[M, K]

	DeletedHook DeletedHook[M, K]
}

type OptionFn[M any, K comparable] func(*options[M, K])

func WithKeywordExpression[M any, K comparable](
	keywordExp func(tx *gorm.DB, keyword string) *gorm.DB) OptionFn[M, K] {
	return func(opts *options[M, K]) {
		opts.keywordExpression = keywordExp
	}
}

func WithCustomConditions[M any, K comparable](conditions ...clause.Expression) OptionFn[M, K] {
	return func(opts *options[M, K]) {
		opts.conditions = conditions
	}
}
