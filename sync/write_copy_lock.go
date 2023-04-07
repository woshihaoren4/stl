package sync

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type Clone[T any] interface {
	clone() T
}
type Default[T any] interface {
	def() T
}

type COWLock[T any] struct {
	inner           *T
	index           uint64
	lock            sync.Mutex
	forcedDelayDest bool
}

func NewCOWLock[T any](def *T) COWLock[T] {
	index := uint64(uintptr(unsafe.Pointer(def)))
	return COWLock[T]{
		inner: def,
		index: index,
		lock:  sync.Mutex{},
	}
}

func (c *COWLock[T]) ForceDelayDestruction() {
	c.forcedDelayDest = true
}

func (c *COWLock[T]) Read() *T {
	ptr := atomic.LoadUint64(&c.index)
	val := (*T)(unsafe.Pointer(uintptr(ptr)))
	return val
}

func (c *COWLock[T]) Set(t *T) {
	if t == nil {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	index := uint64(uintptr(unsafe.Pointer(t)))
	atomic.StoreUint64(&c.index, index)
	c.delayDestruction()
	c.inner = t
}

func (c *COWLock[T]) Write(handle func(*T) *T) { //保持一个更新的原子性
	c.lock.Lock()
	defer c.lock.Unlock()
	t := *c.Read()
	ptr := handle(&t)
	if ptr == nil {
		return
	}
	index := uint64(uintptr(unsafe.Pointer(ptr)))
	atomic.StoreUint64(&c.index, index)
	c.delayDestruction()
	c.inner = ptr
}
func (c *COWLock[T]) delayDestruction() { //保持一个更新的原子性
	if !c.forcedDelayDest {
		return
	}
	go func(ptr *T) *T {
		<-time.After(time.Second * 60 * 3) //三分钟后再释放老数据回收
		return ptr
	}(c.inner)
}
