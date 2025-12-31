package xslice_test

import (
	"fmt"
	"strconv"
	"testing"

	"play-ddd/utils/xslice"
)

func TestMap(t *testing.T) {
	var a = []int{1, 2, 3}

	var b = xslice.MapFn(a, func(a int) int { return a * 3 })
	if fmt.Sprint(b) != `[3 6 9]` {
		t.Fatal(fmt.Sprint(b))
	}

	var c = xslice.MapFn(a, func(a int) string { return strconv.Itoa(a) })
	if fmt.Sprint(c) != `[1 2 3]` {
		t.Fatal(fmt.Sprint(c))
	}
}
