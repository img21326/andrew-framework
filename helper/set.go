package helper

type Set[T comparable] struct {
	Items map[T]bool
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		Items: make(map[T]bool),
	}
}

func (s *Set[T]) Add(item T) {
	s.Items[item] = true
}

func (s *Set[T]) Remove(item T) {
	delete(s.Items, item)
}

func (s *Set[T]) Contains(item T) bool {
	_, ok := s.Items[item]
	return ok
}

func (s *Set[T]) Size() int {
	return len(s.Items)
}

func (s *Set[T]) List() []T {
	var result []T
	for k := range s.Items {
		result = append(result, k)
	}
	return result
}
