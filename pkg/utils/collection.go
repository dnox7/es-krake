package utils

import (
	"github.com/rivo/uniseg"
	"golang.org/x/exp/utf8string"
)

func ToSet[T comparable](items []T) []T {
	set := make(map[T]struct{})
	for _, v := range items {
		set[v] = struct{}{}
	}
	res := []T{}
	for ele := range set {
		res = append(res, ele)
	}
	return res
}

func IsSubSet[T comparable](subset, full []T) bool {
	set := make(map[T]struct{}, len(full))
	for _, v := range full {
		set[v] = struct{}{}
	}
	for _, v := range subset {
		if _, ok := set[v]; !ok {
			return false
		}
	}
	return true
}

// SliceUTF8 gets the characters from the beginning to the specified position,
// in UTF-8-based characters.
func SliceUTF8(str string, pos int, addString string) string {
	s := utf8string.NewString(str)
	length := GetStringCount(str)
	if pos >= length {
		return s.Slice(0, length)
	}
	return s.Slice(0, pos) + addString
}

// GetStringCount gets the number of characters in a person's view.
func GetStringCount(str string) int {
	return uniseg.GraphemeClusterCount(str)
}

// Range create int array with item values from start to end
func Range(start, end int) []int {
	var result []int
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result
}

// RangeN create int array with item values from 0 to end
func RangeN(end int) []int {
	return Range(0, end)
}
