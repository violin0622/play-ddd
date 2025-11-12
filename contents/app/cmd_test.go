package app

import (
	"slices"
	"testing"
)

func TestDiffStrings(t *testing.T) {
	tests := []struct {
		name  string
		a     []string
		b     []string
		wantA []string
		wantB []string
	}{
		{
			name:  "两个空切片",
			a:     []string{},
			b:     []string{},
			wantA: []string{},
			wantB: []string{},
		},
		{
			name:  "a 为空，b 有元素",
			a:     []string{},
			b:     []string{"x", "y"},
			wantA: []string{},
			wantB: []string{"x", "y"},
		},
		{
			name:  "a 有元素，b 为空",
			a:     []string{"a", "b"},
			b:     []string{},
			wantA: []string{"a", "b"},
			wantB: []string{},
		},
		{
			name:  "a 和 b 完全相同",
			a:     []string{"a", "b", "c"},
			b:     []string{"a", "b", "c"},
			wantA: []string{},
			wantB: []string{},
		},
		{
			name:  "a 和 b 完全不同",
			a:     []string{"a", "b"},
			b:     []string{"x", "y"},
			wantA: []string{"a", "b"},
			wantB: []string{"x", "y"},
		},
		{
			name:  "a 和 b 有部分重叠",
			a:     []string{"a", "b", "c"},
			b:     []string{"b", "c", "d"},
			wantA: []string{"a"},
			wantB: []string{"d"},
		},
		{
			name:  "a 包含 b 的所有元素",
			a:     []string{"a", "b", "c", "d"},
			b:     []string{"b", "c"},
			wantA: []string{"a", "d"},
			wantB: []string{},
		},
		{
			name:  "b 包含 a 的所有元素",
			a:     []string{"b", "c"},
			b:     []string{"a", "b", "c", "d"},
			wantA: []string{},
			wantB: []string{"a", "d"},
		},
		{
			name:  "a 有重复元素",
			a:     []string{"a", "a", "b"},
			b:     []string{"b", "c"},
			wantA: []string{"a"},
			wantB: []string{"c"},
		},
		{
			name:  "b 有重复元素",
			a:     []string{"a", "b"},
			b:     []string{"b", "b", "c"},
			wantA: []string{"a"},
			wantB: []string{"c"},
		},
		{
			name:  "单元素情况",
			a:     []string{"a"},
			b:     []string{"b"},
			wantA: []string{"a"},
			wantB: []string{"b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotA, gotB := diffStrings(tt.a, tt.b)

			// 排序以便比较（因为 map 遍历顺序不确定）
			slices.Sort(gotA)
			slices.Sort(gotB)
			slices.Sort(tt.wantA)
			slices.Sort(tt.wantB)

			if !slices.Equal(gotA, tt.wantA) {
				t.Errorf("diffStrings() onlyA = %v, want %v", gotA, tt.wantA)
			}
			if !slices.Equal(gotB, tt.wantB) {
				t.Errorf("diffStrings() onlyB = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}
