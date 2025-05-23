// Copyright (c) 2024, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runes

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

var abcd = "abcd"
var faces = "☺☻☹"
var commas = "1,2,3,4"
var dots = "1....2....3....4"

func eq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func sliceOfString(s [][]rune) []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[i] = string(v)
	}
	return result
}

func TestEqualFold(t *testing.T) {
	tests := []struct {
		s        []rune
		t        []rune
		expected bool
	}{
		{[]rune("hello"), []rune("hello"), true},
		{[]rune("Hello"), []rune("hello"), true},
		{[]rune("hello"), []rune("HELLO"), true},
		{[]rune("world"), []rune("word"), false},
		{[]rune("abc"), []rune("def"), false},
		{[]rune(""), []rune(""), true},
		{[]rune("abc"), []rune(""), false},
		{[]rune(""), []rune("def"), false},
	}

	for _, test := range tests {
		result := EqualFold(test.s, test.t)
		assert.Equal(t, test.expected, result)
	}
}

func TestIndex(t *testing.T) {
	tests := []struct {
		txt      []rune
		find     []rune
		expected int
	}{
		{[]rune("hello"), []rune("el"), 1},
		{[]rune("Hello"), []rune("l"), 2},
		{[]rune("world"), []rune("or"), 1},
		{[]rune("abc"), []rune("def"), -1},
		{[]rune(""), []rune("def"), -1},
		{[]rune("abc"), []rune(""), -1},
		{[]rune(""), []rune(""), -1},
	}

	for _, test := range tests {
		result := Index(test.txt, test.find)
		assert.Equal(t, test.expected, result)
	}
}

func TestIndexFold(t *testing.T) {
	tests := []struct {
		txt      []rune
		find     []rune
		expected int
	}{
		{[]rune("hello"), []rune("el"), 1},
		{[]rune("Hello"), []rune("l"), 2},
		{[]rune("world"), []rune("or"), 1},
		{[]rune("abc"), []rune("def"), -1},
		{[]rune(""), []rune("def"), -1},
		{[]rune("abc"), []rune(""), -1},
		{[]rune(""), []rune(""), -1},
		{[]rune("hello"), []rune("EL"), 1},
		{[]rune("Hello"), []rune("L"), 2},
		{[]rune("world"), []rune("OR"), 1},
		{[]rune("abc"), []rune("DEF"), -1},
		{[]rune(""), []rune("DEF"), -1},
		{[]rune("abc"), []rune(""), -1},
		{[]rune(""), []rune(""), -1},
	}

	for _, test := range tests {
		result := IndexFold(test.txt, test.find)
		assert.Equal(t, test.expected, result)
	}
}

type IndexFuncTest struct {
	in          string
	f           predicate
	first, last int
}

var indexFuncTests = []IndexFuncTest{
	{"", isValidRune, -1, -1},
	{"abc", isDigit, -1, -1},
	{"0123", isDigit, 0, 3},
	{"a1b", isDigit, 1, 1},
	{space, isSpace, 0, len([]rune(space)) - 1},
	{"\u0e50\u0e5212hello34\u0e50\u0e51", isDigit, 0, 12},
	{"\u2C6F\u2C6F\u2C6F\u2C6FABCDhelloEF\u2C6F\u2C6FGH\u2C6F\u2C6F", isUpper, 0, 20},
	{"12\u0e50\u0e52hello34\u0e50\u0e51", not(isDigit), 4, 8},

	// tests of invalid UTF-8
	{"\x801", isDigit, 1, 1},
	{"\x80abc", isDigit, -1, -1},
	{"\xc0a\xc0", isValidRune, 1, 1},
	{"\xc0a\xc0", not(isValidRune), 0, 2},
	{"\xc0☺\xc0", not(isValidRune), 0, 2},
	{"\xc0☺\xc0\xc0", not(isValidRune), 0, 3},
	{"ab\xc0a\xc0cd", not(isValidRune), 2, 4},
	{"a\xe0\x80cd", not(isValidRune), 1, 2},
}

func TestIndexFunc(t *testing.T) {
	for _, tc := range indexFuncTests {
		first := IndexFunc([]rune(tc.in), tc.f.f)
		if first != tc.first {
			t.Errorf("IndexFunc(%q, %s) = %d; want %d", tc.in, tc.f.name, first, tc.first)
		}
		last := LastIndexFunc([]rune(tc.in), tc.f.f)
		if last != tc.last {
			t.Errorf("LastIndexFunc(%q, %s) = %d; want %d", tc.in, tc.f.name, last, tc.last)
		}
	}
}

const space = "\t\v\r\f\n\u0085\u00a0\u2000\u3000"

type predicate struct {
	f    func(r rune) bool
	name string
}

var isSpace = predicate{unicode.IsSpace, "IsSpace"}
var isDigit = predicate{unicode.IsDigit, "IsDigit"}
var isUpper = predicate{unicode.IsUpper, "IsUpper"}
var isValidRune = predicate{
	func(r rune) bool {
		return r != utf8.RuneError
	},
	"IsValidRune",
}

func not(p predicate) predicate {
	return predicate{
		func(r rune) bool {
			return !p.f(r)
		},
		"not " + p.name,
	}
}

type ReplaceTest struct {
	in       string
	old, new string
	n        int
	out      string
}

var ReplaceTests = []ReplaceTest{
	{"hello", "l", "L", 0, "hello"},
	{"hello", "l", "L", -1, "heLLo"},
	{"hello", "x", "X", -1, "hello"},
	{"", "x", "X", -1, ""},
	{"radar", "r", "<r>", -1, "<r>ada<r>"},
	// {"", "", "<>", -1, "<>"},
	{"banana", "a", "<>", -1, "b<>n<>n<>"},
	{"banana", "a", "<>", 1, "b<>nana"},
	{"banana", "a", "<>", 1000, "b<>n<>n<>"},
	{"banana", "an", "<>", -1, "b<><>a"},
	{"banana", "ana", "<>", -1, "b<>na"},
	// {"banana", "", "<>", -1, "<>b<>a<>n<>a<>n<>a<>"},
	// {"banana", "", "<>", 10, "<>b<>a<>n<>a<>n<>a<>"},
	// {"banana", "", "<>", 6, "<>b<>a<>n<>a<>n<>a"},
	// {"banana", "", "<>", 5, "<>b<>a<>n<>a<>na"},
	// {"banana", "", "<>", 1, "<>banana"},
	{"banana", "a", "a", -1, "banana"},
	{"banana", "a", "a", 1, "banana"},
	// {"☺☻☹", "", "<>", -1, "<>☺<>☻<>☹<>"},
}

func TestReplace(t *testing.T) {
	for _, tt := range ReplaceTests {
		in := append([]rune(tt.in), []rune("<spare>")...)
		in = in[:len(tt.in)]
		out := Replace(in, []rune(tt.old), []rune(tt.new), tt.n)
		if s := string(out); s != tt.out {
			t.Errorf("Replace(%q, %q, %q, %d) = %q, want %q", tt.in, tt.old, tt.new, tt.n, s, tt.out)
		}
		if cap(in) == cap(out) && &in[:1][0] == &out[:1][0] {
			t.Errorf("Replace(%q, %q, %q, %d) didn't copy", tt.in, tt.old, tt.new, tt.n)
		}
		if tt.n == -1 {
			out := ReplaceAll(in, []rune(tt.old), []rune(tt.new))
			if s := string(out); s != tt.out {
				t.Errorf("ReplaceAll(%q, %q, %q) = %q, want %q", tt.in, tt.old, tt.new, s, tt.out)
			}
		}
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		r        []rune
		count    int
		expected []rune
	}{
		{[]rune("hello"), 0, []rune{}},
		{[]rune("hello"), 1, []rune("hello")},
		{[]rune("hello"), 2, []rune("hellohello")},
		{[]rune("world"), 3, []rune("worldworldworld")},
		{[]rune(""), 5, []rune("")},
	}

	for _, test := range tests {
		result := Repeat(test.r, test.count)
		assert.Equal(t, test.expected, result)
	}
}

type SplitTest struct {
	s   string
	sep string
	n   int
	a   []string
}

var splittests = []SplitTest{
	{abcd, "a", 0, nil},
	{abcd, "a", -1, []string{"", "bcd"}},
	{abcd, "z", -1, []string{"abcd"}},
	{commas, ",", -1, []string{"1", "2", "3", "4"}},
	{dots, "...", -1, []string{"1", ".2", ".3", ".4"}},
	{faces, "☹", -1, []string{"☺☻", ""}},
	{faces, "~", -1, []string{faces}},
	{"1 2 3 4", " ", 3, []string{"1", "2", "3 4"}},
	{"1 2", " ", 3, []string{"1", "2"}},
	{"bT", "T", math.MaxInt / 4, []string{"b", ""}},
	// {"\xff-\xff", "-", -1, []string{"\xff", "\xff"}},
}

func TestSplit(t *testing.T) {
	for _, tt := range splittests {
		a := SplitN([]rune(tt.s), []rune(tt.sep), tt.n)

		// Appending to the results should not change future results.
		var x []rune
		for _, v := range a {
			x = append(v, 'z')
		}

		result := sliceOfString(a)
		if !eq(result, tt.a) {
			t.Errorf(`Split(%q, %q, %d) = %v; want %v`, tt.s, tt.sep, tt.n, result, tt.a)
			continue
		}
		if tt.n == 0 || len(a) == 0 {
			continue
		}

		if want := tt.a[len(tt.a)-1] + "z"; string(x) != want {
			t.Errorf("last appended result was %s; want %s", string(x), want)
		}

		s := Join(a, []rune(tt.sep))
		if string(s) != tt.s {
			t.Errorf(`Join(Split(%q, %q, %d), %q) = %q`, tt.s, tt.sep, tt.n, tt.sep, s)
		}
		if tt.n < 0 {
			b := Split([]rune(tt.s), []rune(tt.sep))
			if !reflect.DeepEqual(a, b) {
				t.Errorf("Split disagrees withSplitN(%q, %q, %d) = %v; want %v", tt.s, tt.sep, tt.n, b, a)
			}
		}
		if len(a) > 0 {
			in, out := a[0], s
			if cap(in) == cap(out) && &in[:1][0] == &out[:1][0] {
				t.Errorf("Join(%#v, %q) didn't copy", a, tt.sep)
			}
		}
	}
}

var splitaftertests = []SplitTest{
	{abcd, "a", -1, []string{"a", "bcd"}},
	{abcd, "z", -1, []string{"abcd"}},
	{commas, ",", -1, []string{"1,", "2,", "3,", "4"}},
	{dots, "...", -1, []string{"1...", ".2...", ".3...", ".4"}},
	{faces, "☹", -1, []string{"☺☻☹", ""}},
	{faces, "~", -1, []string{faces}},
	{"1 2 3 4", " ", 3, []string{"1 ", "2 ", "3 4"}},
	{"1 2 3", " ", 3, []string{"1 ", "2 ", "3"}},
	{"1 2", " ", 3, []string{"1 ", "2"}},
}

func TestSplitAfter(t *testing.T) {
	for _, tt := range splitaftertests {
		a := SplitAfterN([]rune(tt.s), []rune(tt.sep), tt.n)

		// Appending to the results should not change future results.
		var x []rune
		for _, v := range a {
			x = append(v, 'z')
		}

		result := sliceOfString(a)
		if !eq(result, tt.a) {
			t.Errorf(`Split(%q, %q, %d) = %v; want %v`, tt.s, tt.sep, tt.n, result, tt.a)
			continue
		}

		if want := tt.a[len(tt.a)-1] + "z"; string(x) != want {
			t.Errorf("last appended result was %s; want %s", string(x), want)
		}

		s := Join(a, nil)
		if string(s) != tt.s {
			t.Errorf(`Join(Split(%q, %q, %d), %q) = %q`, tt.s, tt.sep, tt.n, tt.sep, s)
		}
		if tt.n < 0 {
			b := SplitAfter([]rune(tt.s), []rune(tt.sep))
			if !reflect.DeepEqual(a, b) {
				t.Errorf("SplitAfter disagrees withSplitAfterN(%q, %q, %d) = %v; want %v", tt.s, tt.sep, tt.n, b, a)
			}
		}
	}
}

type FieldsTest struct {
	s string
	a []string
}

var fieldstests = []FieldsTest{
	{"", []string{}},
	{" ", []string{}},
	{" \t ", []string{}},
	{"  abc  ", []string{"abc"}},
	{"1 2 3 4", []string{"1", "2", "3", "4"}},
	{"1  2  3  4", []string{"1", "2", "3", "4"}},
	{"1\t\t2\t\t3\t4", []string{"1", "2", "3", "4"}},
	{"1\u20002\u20013\u20024", []string{"1", "2", "3", "4"}},
	{"\u2000\u2001\u2002", []string{}},
	{"\n™\t™\n", []string{"™", "™"}},
	{faces, []string{faces}},
}

func TestFields(t *testing.T) {
	for _, tt := range fieldstests {
		b := []rune(tt.s)
		a := Fields(b)

		// Appending to the results should not change future results.
		var x []rune
		for _, v := range a {
			x = append(v, 'z')
		}

		result := sliceOfString(a)
		if !eq(result, tt.a) {
			t.Errorf("Fields(%q) = %v; want %v", tt.s, a, tt.a)
			continue
		}

		if string(b) != tt.s {
			t.Errorf("slice changed to %s; want %s", string(b), tt.s)
		}
		if len(tt.a) > 0 {
			if want := tt.a[len(tt.a)-1] + "z"; string(x) != want {
				t.Errorf("last appended result was %s; want %s", string(x), want)
			}
		}
	}
}

func TestFieldsFunc(t *testing.T) {
	for _, tt := range fieldstests {
		a := FieldsFunc([]rune(tt.s), unicode.IsSpace)
		result := sliceOfString(a)
		if !eq(result, tt.a) {
			t.Errorf("FieldsFunc(%q, unicode.IsSpace) = %v; want %v", tt.s, a, tt.a)
			continue
		}
	}
	pred := func(c rune) bool { return c == 'X' }
	var fieldsFuncTests = []FieldsTest{
		{"", []string{}},
		{"XX", []string{}},
		{"XXhiXXX", []string{"hi"}},
		{"aXXbXXXcX", []string{"a", "b", "c"}},
	}
	for _, tt := range fieldsFuncTests {
		b := []rune(tt.s)
		a := FieldsFunc(b, pred)

		// Appending to the results should not change future results.
		var x []rune
		for _, v := range a {
			x = append(v, 'z')
		}

		result := sliceOfString(a)
		if !eq(result, tt.a) {
			t.Errorf("FieldsFunc(%q) = %v, want %v", tt.s, a, tt.a)
		}

		if string(b) != tt.s {
			t.Errorf("slice changed to %s; want %s", string(b), tt.s)
		}
		if len(tt.a) > 0 {
			if want := tt.a[len(tt.a)-1] + "z"; string(x) != want {
				t.Errorf("last appended result was %s; want %s", string(x), want)
			}
		}
	}
}

var containsTests = []struct {
	b, subslice []rune
	want        bool
}{
	{[]rune("hello"), []rune("hel"), true},
	{[]rune("日本語"), []rune("日本"), true},
	{[]rune("hello"), []rune("Hello, world"), false},
	{[]rune("東京"), []rune("京東"), false},
}

func TestContains(t *testing.T) {
	for _, tt := range containsTests {
		if got := Contains(tt.b, tt.subslice); got != tt.want {
			t.Errorf("Contains(%q, %q) = %v, want %v", tt.b, tt.subslice, got, tt.want)
		}
	}
}

var ContainsRuneTests = []struct {
	b        []rune
	r        rune
	expected bool
}{
	{[]rune(""), 'a', false},
	{[]rune("a"), 'a', true},
	{[]rune("aaa"), 'a', true},
	{[]rune("abc"), 'y', false},
	{[]rune("abc"), 'c', true},
	{[]rune("a☺b☻c☹d"), 'x', false},
	{[]rune("a☺b☻c☹d"), '☻', true},
	{[]rune("aRegExp*"), '*', true},
}

func TestContainsRune(t *testing.T) {
	for _, ct := range ContainsRuneTests {
		if ContainsRune(ct.b, ct.r) != ct.expected {
			t.Errorf("ContainsRune(%q, %q) = %v, want %v",
				ct.b, ct.r, !ct.expected, ct.expected)
		}
	}
}

func TestContainsFunc(t *testing.T) {
	for _, ct := range ContainsRuneTests {
		if ContainsFunc(ct.b, func(r rune) bool {
			return ct.r == r
		}) != ct.expected {
			t.Errorf("ContainsFunc(%q, func(%q)) = %v, want %v",
				ct.b, ct.r, !ct.expected, ct.expected)
		}
	}
}

type TrimTest struct {
	f            string
	in, arg, out string
}

var trimTests = []TrimTest{
	{"Trim", "abba", "a", "bb"},
	{"Trim", "abba", "ab", ""},
	{"TrimLeft", "abba", "ab", ""},
	{"TrimRight", "abba", "ab", ""},
	{"TrimLeft", "abba", "a", "bba"},
	{"TrimLeft", "abba", "b", "abba"},
	{"TrimRight", "abba", "a", "abb"},
	{"TrimRight", "abba", "b", "abba"},
	{"Trim", "<tag>", "<>", "tag"},
	{"Trim", "* listitem", " *", "listitem"},
	{"Trim", `"quote"`, `"`, "quote"},
	{"Trim", "\u2C6F\u2C6F\u0250\u0250\u2C6F\u2C6F", "\u2C6F", "\u0250\u0250"},
	{"Trim", "\x80test\xff", "\xff", "test"},
	{"Trim", " Ġ ", " ", "Ġ"},
	{"Trim", " Ġİ0", "0 ", "Ġİ"},
	//empty string tests
	{"Trim", "abba", "", "abba"},
	{"Trim", "", "123", ""},
	{"Trim", "", "", ""},
	{"TrimLeft", "abba", "", "abba"},
	{"TrimLeft", "", "123", ""},
	{"TrimLeft", "", "", ""},
	{"TrimRight", "abba", "", "abba"},
	{"TrimRight", "", "123", ""},
	{"TrimRight", "", "", ""},
	// {"TrimRight", "☺\xc0", "☺", "☺\xc0"},
	{"TrimPrefix", "aabb", "a", "abb"},
	{"TrimPrefix", "aabb", "b", "aabb"},
	{"TrimSuffix", "aabb", "a", "aabb"},
	{"TrimSuffix", "aabb", "b", "aab"},
}

type TrimNilTest struct {
	f   string
	in  []rune
	arg string
	out []rune
}

var trimNilTests = []TrimNilTest{
	{"Trim", nil, "", nil},
	{"Trim", []rune{}, "", nil},
	{"Trim", []rune{'a'}, "a", nil},
	{"Trim", []rune{'a', 'a'}, "a", nil},
	{"Trim", []rune{'a'}, "ab", nil},
	{"Trim", []rune{'a', 'b'}, "ab", nil},
	{"Trim", []rune("☺"), "☺", nil},
	{"TrimLeft", nil, "", nil},
	{"TrimLeft", []rune{}, "", nil},
	{"TrimLeft", []rune{'a'}, "a", nil},
	{"TrimLeft", []rune{'a', 'a'}, "a", nil},
	{"TrimLeft", []rune{'a'}, "ab", nil},
	{"TrimLeft", []rune{'a', 'b'}, "ab", nil},
	{"TrimLeft", []rune("☺"), "☺", nil},
	{"TrimRight", nil, "", nil},
	{"TrimRight", []rune{}, "", []rune{}},
	{"TrimRight", []rune{'a'}, "a", []rune{}},
	{"TrimRight", []rune{'a', 'a'}, "a", []rune{}},
	{"TrimRight", []rune{'a'}, "ab", []rune{}},
	{"TrimRight", []rune{'a', 'b'}, "ab", []rune{}},
	{"TrimRight", []rune("☺"), "☺", []rune{}},
	{"TrimPrefix", nil, "", nil},
	{"TrimPrefix", []rune{}, "", []rune{}},
	{"TrimPrefix", []rune{'a'}, "a", []rune{}},
	{"TrimPrefix", []rune("☺"), "☺", []rune{}},
	{"TrimSuffix", nil, "", nil},
	{"TrimSuffix", []rune{}, "", []rune{}},
	{"TrimSuffix", []rune{'a'}, "a", []rune{}},
	{"TrimSuffix", []rune("☺"), "☺", []rune{}},
}

func TestTrim(t *testing.T) {
	toFn := func(name string) (func([]rune, string) []rune, func([]rune, []rune) []rune) {
		switch name {
		case "Trim":
			return Trim, nil
		case "TrimLeft":
			return TrimLeft, nil
		case "TrimRight":
			return TrimRight, nil
		case "TrimPrefix":
			return nil, TrimPrefix
		case "TrimSuffix":
			return nil, TrimSuffix
		default:
			t.Errorf("Undefined trim function %s", name)
			return nil, nil
		}
	}

	for _, tc := range trimTests {
		name := tc.f
		f, fb := toFn(name)
		if f == nil && fb == nil {
			continue
		}
		var actual string
		if f != nil {
			actual = string(f([]rune(tc.in), tc.arg))
		} else {
			actual = string(fb([]rune(tc.in), []rune(tc.arg)))
		}
		if actual != tc.out {
			t.Errorf("%s(%q, %q) = %q; want %q", name, tc.in, tc.arg, actual, tc.out)
		}
	}

	for _, tc := range trimNilTests {
		name := tc.f
		f, fb := toFn(name)
		if f == nil && fb == nil {
			continue
		}
		var actual []rune
		if f != nil {
			actual = f(tc.in, tc.arg)
		} else {
			actual = fb(tc.in, []rune(tc.arg))
		}
		report := func(s []rune) string {
			if s == nil {
				return "nil"
			} else {
				return fmt.Sprintf("%q", s)
			}
		}
		if len(actual) != 0 {
			t.Errorf("%s(%s, %q) returned non-empty value", name, report(tc.in), tc.arg)
		} else {
			actualNil := actual == nil
			outNil := tc.out == nil
			if actualNil != outNil {
				t.Errorf("%s(%s, %q) got nil %t; want nil %t", name, report(tc.in), tc.arg, actualNil, outNil)
			}
		}
	}
}

type TrimFuncTest struct {
	f        predicate
	in       string
	trimOut  []rune
	leftOut  []rune
	rightOut []rune
}

var trimFuncTests = []TrimFuncTest{
	{isSpace, space + " hello " + space,
		[]rune("hello"),
		[]rune("hello " + space),
		[]rune(space + " hello")},
	{isDigit, "\u0e50\u0e5212hello34\u0e50\u0e51",
		[]rune("hello"),
		[]rune("hello34\u0e50\u0e51"),
		[]rune("\u0e50\u0e5212hello")},
	{isUpper, "\u2C6F\u2C6F\u2C6F\u2C6FABCDhelloEF\u2C6F\u2C6FGH\u2C6F\u2C6F",
		[]rune("hello"),
		[]rune("helloEF\u2C6F\u2C6FGH\u2C6F\u2C6F"),
		[]rune("\u2C6F\u2C6F\u2C6F\u2C6FABCDhello")},
	{not(isSpace), "hello" + space + "hello",
		[]rune(space),
		[]rune(space + "hello"),
		[]rune("hello" + space)},
	{not(isDigit), "hello\u0e50\u0e521234\u0e50\u0e51helo",
		[]rune("\u0e50\u0e521234\u0e50\u0e51"),
		[]rune("\u0e50\u0e521234\u0e50\u0e51helo"),
		[]rune("hello\u0e50\u0e521234\u0e50\u0e51")},
	{isValidRune, "ab\xc0a\xc0cd",
		[]rune("\xc0a\xc0"),
		[]rune("\xc0a\xc0cd"),
		[]rune("ab\xc0a\xc0")},
	{not(isValidRune), "\xc0a\xc0",
		[]rune("a"),
		[]rune("a\xc0"),
		[]rune("\xc0a")},
	// The nils returned by TrimLeftFunc are odd behavior, but we need
	// to preserve backwards compatibility.
	{isSpace, "",
		nil,
		nil,
		[]rune("")},
	{isSpace, " ",
		nil,
		nil,
		[]rune("")},
}

func TestTrimFunc(t *testing.T) {
	for _, tc := range trimFuncTests {
		trimmers := []struct {
			name string
			trim func(s []rune, f func(r rune) bool) []rune
			out  []rune
		}{
			{"TrimFunc", TrimFunc, tc.trimOut},
			{"TrimLeftFunc", TrimLeftFunc, tc.leftOut},
			{"TrimRightFunc", TrimRightFunc, tc.rightOut},
		}
		for _, trimmer := range trimmers {
			actual := trimmer.trim([]rune(tc.in), tc.f.f)
			if actual == nil && trimmer.out != nil {
				t.Errorf("%s(%q, %q) = nil; want %q", trimmer.name, tc.in, tc.f.name, trimmer.out)
			}
			if actual != nil && trimmer.out == nil {
				t.Errorf("%s(%q, %q) = %q; want nil", trimmer.name, tc.in, tc.f.name, actual)
			}
			if !Equal(actual, trimmer.out) {
				t.Errorf("%s(%q, %q) = %q; want %q", trimmer.name, tc.in, tc.f.name, actual, trimmer.out)
			}
		}
	}
}
