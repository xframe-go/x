package repository

import (
	"context"

	"github.com/xframe-go/x/requests"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Interface[M any, C ToModel[M], U ToUpdateModel, K comparable] interface {
	List(ctx context.Context, params requests.QueryParams) (data []M, total int64, err error)
	BatchList(ctx context.Context, params requests.QueryParams) (data []M, err error)
	Create(tx *gorm.DB, m *M) error
	Show(ctx context.Context, key K, params requests.QueryParams) (M, error)
	Update(ctx context.Context, tx *gorm.DB, key K, m U) error
	Destroy(ctx context.Context, tx *gorm.DB, keys ...K) error
}

type ToModel[M any] interface {
	ToModel() M
}

type ToUpdateModel interface {
	ToModel() []clause.Assigner
}
