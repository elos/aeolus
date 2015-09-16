package aeolus

// includes checks whether string s is a member of the ss string array
func includes(ss []string, s string) bool {
	for i := range ss {
		if s == ss[i] {
			return true
		}
	}

	return false
}

// reverse returns a copy of s in reversed order
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
