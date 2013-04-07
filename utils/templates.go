package utils

import (
	"reflect"
	"html/template"
)

var FuncMap = template.FuncMap{
	"any": Any,
}

// any reports whether the first argument is equal to
// any of the remaining arguments.
func Any(args ...interface{}) bool {
	if len(args) == 0 {
		return false
	}
	x := args[0]
	switch x := x.(type) {
	case string, int, int64, byte, float32, float64:
		for _, y := range args[1:] {
			if x == y {
				return true
			}
		}
		return false
	}

	for _, y := range args[1:] {
		if reflect.DeepEqual(x, y) {
			return true
		}
	}
	return false
}
