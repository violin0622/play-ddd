package xslice

import (
	"iter"
)

// func MapFn[SliceA ~[]A, SliceB []B, A any, B any](
// 	sa SliceA,
// 	fn func(A) B,
// ) SliceB {
// 	sb := make([]B, 0, len(sa))
// 	for _, a := range sa {
// 		sb = append(sb, fn(a))
// 	}
// 	return sb
// }

func Map[S iter.Seq[A], A, B any](sa S, fn func(A) B) iter.Seq[B] {
	return func(yield func(B) bool) {
		for a := range sa {
			if !yield(fn(a)) {
				return
			}
		}
	}
}

func MapIdx[S iter.Seq2[K, V], K, V, A, B any](
	s S, fn func(K, V) (A, B),
) iter.Seq2[A, B] {
	return func(yeild func(A, B) bool) {
		for k, v := range s {
			if !yeild(fn(k, v)) {
				return
			}
		}
	}
}

func Collect[S iter.Seq2[int, E], E any](s S) []E {
	var arr []E
	for _, e := range s {
		arr = append(arr, e)
	}

	return arr
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
