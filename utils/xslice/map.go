package xslice

import ()

func MapFn[SliceA ~[]A, SliceB []B, A any, B any](
	sa SliceA,
	fn func(A) B,
) SliceB {
	sb := make([]B, 0, len(sa))
	for _, a := range sa {
		sb = append(sb, fn(a))
	}
	return sb
}

func MapIdxFn[SliceA ~[]A, SliceB []B, A any, B any](
	sa SliceA,
	fn func(int, A) B,
) SliceB {
	sb := make([]B, 0, len(sa))
	for i, a := range sa {
		sb = append(sb, fn(i, a))
	}
	return sb
}

func Parallel[A any, B any, SliceA ~[]A, SliceB ~[]B](
	sa SliceA, sb SliceB, fn func(int),
) {
	for i := range min(len(sa), len(sb)) {
		fn(i)
	}
}
