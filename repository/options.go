package repository

import (
	"gorm.io/gorm"
)

type options struct {
	keywordExpression func(tx *gorm.DB, keyword string) *gorm.DB
}

type OptionFn func(*options)

func WithKeywordExpression(keywordExp func(tx *gorm.DB, keyword string) *gorm.DB) OptionFn {
	return func(opts *options) {
		opts.keywordExpression = keywordExp
	}
}
