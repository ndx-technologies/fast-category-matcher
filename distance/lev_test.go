package distance_test

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/ndx-technologies/fast-category-matcher/distance"
)

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected int
	}{
		{"kitten", "sitting", 3},
		{"book", "back", 2},
		{"", "abc", 3},
		{"abc", "", 3},
		{"", "", 0},
		{"a", "a", 0},
		{"a", "b", 1},
		{"ab", "abc", 1},
		{"abc", "ab", 1},
		{"abc", "acb", 2},
		{"abc", "def", 3},
		{"hello", "helo", 1},
		{"world", "word", 1},
		{"distance", "difference", 5},
	}

	for _, test := range tests {
		d := distance.Lev(test.s1, test.s2, nil)
		if d != test.expected {
			t.Error(d, test.expected)
		}
	}
}

func BenchmarkLevenshteinDistance(b *testing.B) {
	testCases := []struct {
		name string
		s1   string
		s2   string
	}{
		{"empty", "", ""},
		{"one empty", "hello", ""},
		{"equal", "kitten", "kitten"},
		{"different", "kitten", "sitting"},
		{"long different", strings.Repeat("a", 256), strings.Repeat("b", 256)},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for b.Loop() {
				distance.Lev(tc.s1, tc.s2, nil)
			}
		})

		b.Run("buff/"+tc.name, func(b *testing.B) {
			buf := make([]byte, 512)
			for b.Loop() {
				distance.Lev(tc.s1, tc.s2, buf)
			}
		})
	}
}

func FuzzLevenshteinDistance(f *testing.F) {
	f.Add("", "")
	f.Add("hello", "")
	f.Add("", "world")
	f.Add("kitten", "sitting")
	f.Add("distance", "difference")
	f.Add(strings.Repeat("a", 100), strings.Repeat("b", 100))

	f.Fuzz(func(t *testing.T, s1, s2 string) {
		if !utf8.ValidString(s1) || !utf8.ValidString(s2) {
			t.Skip("Invalid UTF-8")
		}

		d := distance.Lev(s1, s2, nil)

		if d < 0 {
			t.Error("d < 0")
		}
		if s1 == s2 && d != 0 {
			t.Error(d)
		}
		if s1 == "" && d != len(s2) {
			t.Error(d)
		}
		if s2 == "" && d != len(s1) {
			t.Error(d)
		}
	})
}
