package utils

func Filter[T any](arr []T, test func(T, int, []T) bool) (filtered []T) {
	for idx, elem := range arr {
		if test(elem, idx, arr) {
			filtered = append(filtered, elem)
		}
	}

	return
}

func Map[T any, R any](arr []T, fn func(T, int, []T) R) (mapped []R) {
	for idx, elem := range arr {
		mapped = append(mapped, fn(elem, idx, arr))
	}

	return
}
