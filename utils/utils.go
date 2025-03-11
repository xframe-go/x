package utils

func SliceToAny[T any](slice []T) (res []any) {
	return SliceTo[T, any](slice, func(t T) any {
		return t
	})
}

func SliceTo[T any, V any](slice []T, converter func(t T) V) (res []V) {
	for i := range slice {
		res = append(res, converter(slice[i]))
	}
	return
}
