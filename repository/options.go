package repository

import "gorm.io/gorm/clause"

type options struct {
	keywordExpression func(keyword string) clause.Expression
}

type OptionFn func(*options)

func WithKeywordExpression(keywordExp func(keyword string) clause.Expression) OptionFn {
	return func(opts *options) {
		opts.keywordExpression = keywordExp
	}
}
