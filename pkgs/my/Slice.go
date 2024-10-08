package my

import "unsafe"

type Slice[T any] struct {
	array    unsafe.Pointer // 底层数组
	length   int            // 长度
	capacity int            // 容量
}

func NewSlice[T any](capacity int) *Slice[T] {
	arrPtr := unsafe.Pointer(&make([]T, capacity)[0])
	return &Slice[T]{
		array:    arrPtr,
		length:   0,
		capacity: capacity,
	}
}

func (s *Slice[T]) Len() int {
	return s.length
}

func (s *Slice[T]) Cap() int {
	return s.capacity
}

func (s *Slice[T]) Array() []T {
	return (*(*[]T)(unsafe.Pointer(&s.array)))[:s.length]
}

func (s *Slice[T]) Append(elems ...T) {
	if s.length >= s.capacity {
		s.grow(s.length + 1)
	}
	for i := 0; i < len(elems); i++ {
		elemPtr := unsafe.Pointer(uintptr(s.array) + uintptr(s.length)*unsafe.Sizeof(*new(T)))
		*(*T)(elemPtr) = elems[i]
		s.length += 1
	}
}

func (s *Slice[T]) grow(expCap int) {
	newCap := s.capacity
	doubleCap := newCap + newCap
	// 1.18 以前
	//if expCap > doubleCap {
	//	newCap = expCap
	//} else {
	//	if s.capacity < 1024 {
	//		newCap = doubleCap
	//	} else {
	//		for 0 < newCap && newCap < expCap {
	//			newCap += newCap / 4
	//		}
	//		if newCap <= 0 {
	//			newCap = expCap
	//		}
	//	}
	//}

	// 1.18 及以后
	if expCap > doubleCap {
		newCap = expCap
	} else {
		if s.capacity < 256 {
			newCap = doubleCap
		} else {
			for 0 < newCap && newCap < expCap {
				newCap += (newCap + 3*256) / 4
			}
			if newCap <= 0 {
				newCap = expCap
			}
		}
	}
	newArrPtr := unsafe.Pointer(&make([]T, newCap)[0])
	for i := 0; i < s.length; i++ {
		oldElemPtr := unsafe.Pointer(uintptr(s.array) + uintptr(i)*unsafe.Sizeof(*new(T)))
		newElemPtr := unsafe.Pointer(uintptr(newArrPtr) + uintptr(i)*unsafe.Sizeof(*new(T)))
		*(*T)(newElemPtr) = *(*T)(oldElemPtr)
	}
	s.array = newArrPtr
	s.capacity = newCap
}

func (s *Slice[T]) Get(idx int) T {
	if idx < 0 || idx >= s.length {
		panic("index out of range")
	}
	elemPtr := unsafe.Pointer(uintptr(s.array) + uintptr(idx)*unsafe.Sizeof(*new(T)))
	return *(*T)(elemPtr)
}

func (s *Slice[T]) Set(idx int, elem T) {
	if idx < 0 || idx >= s.length {
		panic("index out of range")
	}
	elemPtr := unsafe.Pointer(uintptr(s.array) + uintptr(idx)*unsafe.Sizeof(*new(T)))
	*(*T)(elemPtr) = elem
}

func (s *Slice[T]) Slice(leftBound, rightBound int) *Slice[T] {
	if leftBound < 0 || rightBound > s.length || leftBound > rightBound {
		panic("invalid slice bounds")
	}
	return &Slice[T]{
		array:    unsafe.Pointer(uintptr(s.array) + uintptr(leftBound)*unsafe.Sizeof(*new(T))),
		length:   rightBound - leftBound,
		capacity: s.capacity - leftBound,
	}
}

func (s *Slice[T]) Del(idx int) {
	if idx < 0 || idx >= s.length {
		panic("index out of range")
	}
	newArrPtr := s.Slice(0, idx)
	newArrPtr.Append(s.Slice(idx+1, s.length).Array()...)
	s.array = newArrPtr.array
	s.length -= 1
	s.capacity = newArrPtr.capacity
}
