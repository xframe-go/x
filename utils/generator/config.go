package generator

import "gorm.io/gen"

type symbol string

var ConfigSymbol = symbol("generator")

type Config struct {
	DaoPath   string // query code path
	ModelPath string // generated model code's package name
	OutFile   string // query code file name, default: gen.go
	// 生成单元测试，默认值 false, 选项: false / true
	// generate unit test for query code
	WithUnitTest bool

	// 当字段允许空时用指针生成
	FieldNullable bool // generate pointer when field is nullable
	// generate pointer when field has default value, to fix problem zero value cannot be assign: https://gorm.io/docs/create.html#Default-Values
	FieldCoverable bool
	// detect integer field's unsigned type, adjust generated data type
	FieldSignable bool
	// 生成带有gorm index 标签的字段
	// generate with gorm index tag
	FieldWithIndexTag bool
	// 生成带有gorm type标签的字段
	// generate with gorm column type tag
	FieldWithTypeTag bool

	GenerateMode gen.GenerateMode

	Config func(g *gen.Generator)
}

func (*Config) Symbol() any {
	return ConfigSymbol
}

var DefaultConfig = &Config{
	DaoPath: "app/http/dao",
}
