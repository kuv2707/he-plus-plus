package utils


func Backwards[T any](arr []T) func(func(i int, x T) bool) {
	return func(yield func(i int, x T) bool) {
		for i := len(arr) - 1; i>= 0; i-- {
			if !yield(i, arr[i]) {
				return
			}
		}
	}
}