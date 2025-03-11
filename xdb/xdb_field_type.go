package xdb

import (
	"gorm.io/gorm"
)

var FieldDatetime = func(col gorm.ColumnType) (dataType string) {
	return "carbon.Carbon"
}
