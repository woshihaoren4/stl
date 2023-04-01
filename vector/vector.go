package vector

import (
	"stl/util"
	"sync"
)

type Vec[T comparable] struct {
	data []T
	len  int
	lock sync.RWMutex
}

func NewVectorFromSlice[T comparable](list []T) Vec[T] {
	return Vec[T]{
		data: list,
		len:  len(list),
		lock: sync.RWMutex{},
	}
}
func NewVector[T comparable](vals ...T) Vec[T] {
	return NewVectorFromSlice(vals)
}
func (v *Vec[T]) Len() (l int) {
	v.read(func() {
		l = v.len
	})
	return
}
func (v *Vec[T]) Push(vals ...T) {
	v.write(func() {
		v.len += len(vals)
		v.data = append(v.data, vals...)
	})
}
func (v *Vec[T]) Insert(position int, val T) {
	v.write(func() {
		if position > v.len || position < 0 {
			util.PanicFmt("%T insert position[%d] over index range[0:%d]", v, position, v.len)
		}
		v.len += 1
		if position == v.len {
			v.data = append(v.data, val)
			return
		}
		buf := make([]T, v.len)
		buf[position] = val
		for i, _ := range v.data {
			if i < position {
				buf[i] = v.data[i]
			} else if i >= position {
				buf[i+1] = v.data[i]
			}
		}
		v.data = buf
	})
}
func (v *Vec[T]) Update(position int, val T) {
	v.write(func() {
		if position > v.len || position < 0 {
			util.PanicFmt("%T store position[%d] over index range[0:%d]", v, position, v.len)
		}
		if position == v.len {
			v.len += 1
			v.data = append(v.data, val)
			return
		}
		v.data[position] = val
	})
}
func (v *Vec[T]) Remove(position int) (val T, ok bool) {
	v.write(func() {
		if position >= v.len || position < 0 {
			return
		}
		ok = true
		val = v.data[position]
		v.len -= 1
		if position == 0 { //第一个
			if v.len == 0 {
				v.data = nil
			} else {
				v.data = v.data[1 : v.len+1]
			}
			return
		}
		if position+1 == v.len { //最后一个
			v.data = v.data[0:position]
			return
		}
		buf := v.data[0:position]
		buf = append(buf, v.data[position+1:v.len+1]...)
		v.data = buf
	})
	return
}
func (v *Vec[T]) Pop() (val T, ok bool) {
	v.write(func() {
		if v.len == 0 {
			return
		}
		ok = true
		v.len -= 1
		val = v.data[v.len]
		if v.len == 0 {
			v.data = nil
		} else {
			v.data = v.data[:v.len]
		}
	})
	return
}
func (v *Vec[T]) Def(position int, handle func(*T, bool)) {
	var ok bool
	v.read(func() {
		if position >= 0 && position < v.len {
			ok = true
			handle(&v.data[position], true)
		}
	})
	if !ok {
		handle(nil, false)
	}
	return
}
func (v *Vec[T]) UnsafeDef(position int) (val *T, ok bool) {
	v.Def(position, func(t *T, exist bool) {
		ok = exist
		val = t
	})
	return
}
func (v *Vec[T]) RangeDef(start, end int, handle func([]T, bool)) {
	var ok bool
	v.read(func() {
		if start >= 0 && start < end && end <= v.len {
			ok = true
			handle(v.data[start:end], true)
		}
	})
	if !ok {
		handle(nil, false)
	}
}
func (v *Vec[T]) Load(position int) (val T, ok bool) {
	v.Def(position, func(t *T, exist bool) {
		if exist {
			ok = exist
			val = *t
		}
	})
	return
}
func (v *Vec[T]) RangeLoad(start, end int) (buf []T, ok bool) {
	v.RangeDef(start, end, func(data []T, exist bool) {
		if exist {
			ok = true
			buf = make([]T, len(data))
			copy(buf, data)
		}
	})
	return
}
func (v *Vec[T]) UnsafeRangeDef(start, end int) (buf []T, ok bool) {
	v.RangeDef(start, end, func(data []T, exist bool) {
		ok = exist
		buf = data
	})
	return
}
func (v *Vec[T]) ContainMap(handle func(T) bool) (ok bool) {
	v.read(func() {
		for _, i := range v.data {
			if handle(i) {
				ok = true
				return
			}
		}
	})
	return
}
func (v *Vec[T]) Contain(des T) (ok bool) {
	return v.ContainMap(func(src T) bool {
		return des == src
	})
}
func (v *Vec[T]) FilterSlice(handle func([]T) []T) {
	v.write(func() {
		v.data = handle(v.data)
		v.len = len(v.data)
	})
}
func (v *Vec[T]) Filter(handle func(*T) bool) {
	v.FilterSlice(func(data []T) []T {
		var buf []T
		for i, _ := range data {
			if handle(&data[i]) {
				buf = append(buf, data[i])
			}
		}
		return buf
	})
}
func (v *Vec[T]) Foreach(foreach func(int, T)) {
	v.read(func() {
		for i, v := range v.data {
			foreach(i, v)
		}
	})
}
func (v *Vec[T]) Reverse() {
	v.FilterSlice(func(src []T) []T {
		l := len(src)
		des := make([]T, l)
		for i, _ := range src {
			des[l-i-1] = src[i]
		}
		return des
	})
}
func (v *Vec[T]) UnsafeIntoSlice() []T {
	return v.data
}
func (v *Vec[T]) Clean() {
	v.write(func() {
		v.len = 0
		v.data = nil
	})
}

func (v *Vec[T]) write(handle func()) {
	v.lock.Lock()
	defer v.lock.Unlock()
	handle()

}
func (v *Vec[T]) read(handle func()) {
	v.lock.RLock()
	defer v.lock.RUnlock()
	handle()

}
