package util

import "fmt"

func PanicFmt(format string, vals ...any) {
	panic(fmt.Sprintf(format, vals...))
}
