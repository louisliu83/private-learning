package utils

import (
	"testing"
)

func TestNilSliceToCommaSeperatedString(t *testing.T) {
	var s []string
	if SliceToCommaSeperatedString(s) != "" {
		t.Fail()
	}
}

func TestEmptySliceToCommaSeperatedString(t *testing.T) {
	var s = make([]string, 0)
	if SliceToCommaSeperatedString(s) != "" {
		t.Fail()
	}
}

func TestEmptyStringToSlice(t *testing.T) {
	var s string = ""
	if CommaSeperatedStringToSlice(s) == nil {
		t.Fail()
	}
	if len(CommaSeperatedStringToSlice(s)) > 0 {
		t.Fail()
	}
}

func TestUnion(t *testing.T) {
	s1 := []string{"a", "b", "c"}
	s2 := []string{"b", "c", "d"}
	s := Union(s1, s2)
	t.Log(s)
	if len(s) != 4 {
		t.Fail()
	}
}

func TestIntersect(t *testing.T) {
	s1 := []string{"a", "b", "c"}
	s2 := []string{"b", "c", "d"}
	s := Intersect(s1, s2)
	t.Log(s)
	if len(s) != 1 {
		t.Fail()
	}
}
