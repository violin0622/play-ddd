package xslice

import (
	"iter"
)

func Foreach[S iter.Seq[A], A any](sa S, fn func(A)) {
	for a := range sa {
		fn(a)
	}
}

func ForeachE[S iter.Seq[A], A any](sa S, fn func(A) error) error {
	for a := range sa {
		if err := fn(a); err != nil {
			return err
		}
	}

	return nil
}

func ForeachIdx[S iter.Seq2[K, V], K, V any](s S, fn func(K, V)) {
	for k, v := range s {
		fn(k, v)
	}
}

func ForeachIdxE[S iter.Seq2[K, V], K, V any](s S, fn func(K, V) error) error {
	for k, v := range s {
		if err := fn(k, v); err != nil {
			return err
		}
	}

	return nil
}
