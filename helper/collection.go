package helper

type Collection[T any] []T

func NewCollection[T any](items ...T) Collection[T] {
	return items
}

func (c Collection[T]) Map(f func(T) T) Collection[T] {
	var result Collection[T]
	for _, v := range c {
		result = append(result, f(v))
	}
	return result
}

func (c Collection[T]) Filter(f func(T) bool) Collection[T] {
	var result Collection[T]
	for _, v := range c {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func (c Collection[T]) Reduce(f func(T, T) T) T {
	var result T
	for _, v := range c {
		result = f(result, v)
	}
	return result
}

func (c Collection[T]) Find(f func(T) bool) *T {
	for _, v := range c {
		if f(v) {
			return &v
		}
	}
	return nil
}

func (c Collection[T]) FindIndex(f func(T) bool) int {
	for i, v := range c {
		if f(v) {
			return i
		}
	}
	return -1
}

func (c *Collection[T]) ForEach(f func(*T)) {
	for i := range *c {
		f(&(*c)[i])
	}
}

func (c Collection[T]) Contain(item T, equalFunc func(T, T) bool) bool {
	for _, v := range c {
		if equalFunc(v, item) {
			return true
		}
	}
	return false
}

func (c Collection[T]) ToMap(keyFunc func(T) string) map[string]T {
	result := make(map[string]T)
	for _, v := range c {
		result[keyFunc(v)] = v
	}
	return result
}

func (c Collection[T]) Join(sep string, stringFunc func(T) string) string {
	var result string
	for i, v := range c {
		if i == 0 {
			result += stringFunc(v)
		} else {
			result += sep + stringFunc(v)
		}
	}
	return result
}
