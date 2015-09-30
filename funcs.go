package aeolus

import "strings"

// includes checks whether string s is a member of the ss string array
func includes(ss []string, s string) bool {
	for i := range ss {
		if s == ss[i] {
			return true
		}
	}

	return false
}

// Reverse returns a copy of s in reversed order
// Note: exported because it is used in ego
func Reverse(s []string) []string {
	rs := make([]string, len(s))
	j := 0
	for i := len(s) - 1; i >= 0; i-- {
		rs[j] = s[i]
		j++
	}
	return rs
}

// NormalizeName replaces hyphens and spaces in names with underscores
// i.e., this is my name => this_is_my_name
// normal-lize => normal_lize
// mis_match ed-stuff => mis_match_ed_stuff
func NormalizeName(s string) string {
	s = strings.Replace(s, "-", "_", -1)
	s = strings.Replace(s, " ", "_", -1)
	return s
}

var ProcessName = NormalizeName
