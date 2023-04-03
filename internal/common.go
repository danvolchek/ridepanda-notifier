package internal

func Map[T, K any](items []T, f func(T) K) []K {
	result := make([]K, len(items))
	for i, item := range items {
		result[i] = f(item)
	}

	return result
}

func Filter[T any](items []T, f func(T) bool) []T {
	var result []T
	for _, item := range items {
		if f(item) {
			result = append(result, item)
		}
	}

	return result
}