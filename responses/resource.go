package responses

type Resource[M any] struct {
	Base
}

func NewResource[M any]() *Resource[M] {
	return &Resource[M]{}
}
