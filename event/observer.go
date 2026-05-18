package event

import "gorm.io/gorm"

type Observer[M any] struct {
	TX     *gorm.DB
	Origin *M
	Model  *M
}
