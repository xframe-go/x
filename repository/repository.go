package repository

import (
	"context"

	"github.com/xframe-go/x/requests"
	"github.com/xframe-go/x/x"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository[M any, C ToModel[M], U ToUpdateModel, K comparable] struct {
	primaryKeyName   string
	primaryKeyGetter PrimaryKeyGetter[M, K]
}

func New[M any, C ToModel[M], U ToUpdateModel, K comparable](primaryKeyGetter PrimaryKeyGetter[M, K]) *Repository[M, C, U, K] {
	return &Repository[M, C, U, K]{
		primaryKeyName:   "id",
		primaryKeyGetter: primaryKeyGetter,
	}
}

func (repo *Repository[M, C, U, K]) List(ctx context.Context, params requests.QueryParams) (data []M, total int64, err error) {
	var mo M
	tx := x.DB().WithContext(ctx).Model(mo)

	tx = repo.attachQuery(tx, params)

	if err = tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return make([]M, 0), 0, nil
	}

	if err = repo.attachPagination(tx, params).Find(&data).Error; err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (repo *Repository[M, C, U, K]) BatchList(ctx context.Context, params requests.QueryParams) (data []M, err error) {
	var mo M
	tx := x.DB().WithContext(ctx).Model(mo)

	tx = repo.attachQuery(tx, params)

	if err = tx.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (*Repository[M, C, U, K]) Create(tx *gorm.DB, m *M) error {
	return tx.Create(m).Error
}

func (repo *Repository[M, C, U, K]) Show(ctx context.Context, key K, params requests.QueryParams) (m M, err error) {
	tx := x.DB().WithContext(ctx).Where("id", key)

	repo.attachPreload(tx, params)

	err = tx.First(&m).Error

	return
}

func (repo *Repository[M, C, U, K]) Update(ctx context.Context, tx *gorm.DB, key K, m U) error {
	updates := m.ToModel()
	_, err := x.Model[M](tx).Where(repo.primaryKeyName, key).Set(updates...).Update(ctx)
	return err
}

func (repo *Repository[M, C, U, K]) Destroy(ctx context.Context, tx *gorm.DB, keys ...K) error {
	if len(keys) == 0 {
		return nil
	}
	return tx.Where("id", keys).Delete(ctx).Error
}

func (repo *Repository[M, C, U, K]) GetByPrimaryKey(ctx context.Context, key K) (M, error) {
	return x.Model[M]().Where("id", key).First(ctx)
}

func (repo *Repository[M, C, U, K]) attachQuery(tx *gorm.DB, params requests.QueryParams) *gorm.DB {
	for fe, sorter := range params.Sorter {
		tx = tx.Order(clause.OrderByColumn{
			Column: clause.Column{
				Name: fe,
			},
			Desc: sorter == "desc",
		})
	}

	if len(params.Sorter) == 0 {
		tx = tx.Order(clause.OrderByColumn{
			Column: clause.Column{
				Name: "id",
			},
			Desc: true,
		})
	}

	repo.attachPreload(tx, params)

	for _, filter := range params.Filters {
		switch filter.Operator {
		case requests.Equal:
			tx = tx.Where(filter.Field, filter.Value)
		case requests.NotEqual:
			tx = tx.Where(filter.Field+" != ?", filter.Value)
		case requests.Greater:
			tx = tx.Where(filter.Field+" > ?", filter.Value)
		case requests.Less:
			tx = tx.Where(filter.Field+" < ?", filter.Value)
		case requests.GreaterEq:
			tx = tx.Where(filter.Field+" >= ?", filter.Value)
		case requests.LessEq:
			tx = tx.Where(filter.Field+" <= ?", filter.Value)
		case requests.In:
			tx = tx.Where(filter.Field+" IN ?", filter.Value)
		case requests.NotIn:
			tx = tx.Where(filter.Field+" NOT IN ?", filter.Value)
		case requests.Contains:
			tx = tx.Where(filter.Field+" LIKE ?", "%"+filter.Value.(string)+"%")
		case requests.Between:
			if values, ok := filter.Value.([]string); ok && len(values) == 2 {
				tx = tx.Where(filter.Field+" BETWEEN ? AND ?", values[0], values[1])
			}
		}
	}

	return tx
}

func (*Repository[M, C, U, K]) attachPreload(tx *gorm.DB, params requests.QueryParams) *gorm.DB {
	for _, relation := range params.Preload {
		tx = tx.Preload(relation)
	}
	return tx
}

func (*Repository[M, C, U, K]) attachPagination(tx *gorm.DB, params requests.QueryParams) *gorm.DB {
	offset := (params.Page - 1) * params.PageSize
	return tx.Limit(params.PageSize).Offset(offset)
}
