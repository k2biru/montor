package models

import (
	"sort"
)

// create slice for sorting
type keyValue[T any] struct {
	Key   uint8
	Value T
}

func sortAscending[T any](list map[uint8]T) []keyValue[T] {
	// copy map to slice
	s := make([]keyValue[T], 0, len(list))
	for k, v := range list {
		s = append(s, keyValue[T]{k, v})
	}
	// sorting ascending
	sort.SliceStable(s, func(i, j int) bool {
		return s[i].Key < s[j].Key
	})
	return s
}
