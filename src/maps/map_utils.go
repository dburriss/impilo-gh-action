package maps

type Pair[T, U any] struct {
	First  T
	Second U
}

func MapSlice[T any, U any](xs []T, f func(T) U) []U {
	var result []U
	for _, x := range xs {
		result = append(result, f(x))
	}
	return result
}

func Map[K1 comparable, T1 any, K2 comparable, T2 any](m map[K1]T1, f func(K1, T1) (K2, T2)) map[K2]T2 {
	var result map[K2]T2
	for k1, v1 := range m {
		k2, v2 := f(k1, v1)
		result[k2] = v2
	}
	return result
}

func ToSlice[K comparable, V any](m map[K]V) []Pair[K, V] {
	var result []Pair[K, V]
	for k1, v1 := range m {
		result = append(result, Pair[K, V]{k1, v1})
	}
	return result
}

func ItemExists[K comparable, V any](key K, m map[K]V) bool {
	_, exists := m[key]
	return exists
}

func Append[K comparable, V any](m1 map[K]V, m2 map[K]V) {
	for k, _ := range m2 {
		if !ItemExists(k, m1) {
			m1[k] = m2[k]
		}
	}
}
