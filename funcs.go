package aeolus

// Reverse returns a copy of s reversed order
func Reverse(s []string) []string {
	rs := make([]string, len(s))
	j := 0
	for i := len(s) - 1; i >= 0; i-- {
		rs[j] = s[i]
		j++
	}
	return rs
}
