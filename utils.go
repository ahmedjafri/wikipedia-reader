package main

import (
	"strings"
)

func dedup(s []string, caseSensitive bool) []string {

	var rt []string
	m := map[string]struct{}{}

	for _, ss := range s {
		if !caseSensitive {
			ss = strings.ToLower(ss)
		}

		m[ss] = struct{}{}
	}

	for k, _ := range m {
		rt = append(rt, k)
	}

	return rt
}
