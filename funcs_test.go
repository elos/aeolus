package aeolus

import (
	"strings"
	"testing"
)

func TestIncludes(t *testing.T) {
	s := []string{"one", "two", "three"}
	one := "one"
	four := "four"

	if !includes(s, one) {
		t.Errorf("%+v should contain the string: %s", s, one)
	}

	if includes(s, four) {
		t.Errorf("%+v should not contain the string: %s", s, four)
	}

	if includes(s, "") {
		t.Errorf("%+v should not contain empty string")
	}
}

// checks equality of two slices
func eq(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

func TestReverse(t *testing.T) {
	// quickly sanity check our helper function
	if eq([]string{"1"}, []string{"1", "2"}) != false {
		t.Fatal("Helper function string array eq doesn't work")
	}
	hw := []string{"hello", "world"}
	if eq(hw, hw) != true {
		t.Fatal("Helper function string array eq doesn't work")
	}

	s := strings.Split("0 1 2 3 4 5 6 7 8 9", " ")
	palindrome := strings.Split("r a c e c a r", " ")

	if !eq(Reverse(s), strings.Split("9 8 7 6 5 4 3 2 1 0", " ")) {
		t.Errorf("Expected reverse of %s to be '9876543210'", s)
	}

	if !eq(Reverse(palindrome), palindrome) {
		t.Errorf("Expected reverse of %s to be %s", palindrome)
	}

	if !eq(Reverse([]string{}), []string{}) {
		t.Errorf("Expected reverse of the empty string to be the empty string")
	}
}

func TestNormalizeName(t *testing.T) {
	cases := map[string]string{
		"this is my name":    "this_is_my_name",
		"normal-lize":        "normal_lize",
		"mis_match ed-stuff": "mis_match_ed_stuff",
	}

	for example, solution := range cases {
		if NormalizeName(example) != solution {
			t.Errorf("Expected NoramlizeName(%s) to be %s, got %s", example, solution, NormalizeName(example))
		}
	}
}
