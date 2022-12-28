package slices

func SliceFind[K any](xs []K, predicate func(x K) bool) (K, bool) {
	for _, el := range xs {
		if predicate(el) {
			return el, true
		}
	}
	var empty K
	return empty, false
}

func SliceContains[K comparable](xs []K, x K) bool {
	_, contains := SliceFind(xs, func(y K) bool { return x == y })
	return contains
}
