package util

import (
	"reflect"
)

func IsPtr(val any) bool {
	ty := reflect.TypeOf(val)
	kind := ty.Kind()
	return kind == reflect.Pointer || kind == reflect.UnsafePointer
}

type GoSlice struct {
	array uintptr
	len   int
	cap   int
}

func DepthCopy[T any](t *T) T {
	var val T = *t
	return val
}
