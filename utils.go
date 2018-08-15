package main

import (
	"io/ioutil"
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

func ReadGraphFile(filename string) *Graph {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	g := NewGraph()
	dataString := string(data)
	lines := strings.Split(dataString, "\n")
	for _, line := range lines {
		lineParts := strings.Split(line, "=")
		if len(lineParts) != 2 {
			Log.Error("Line does not contain two parts", "Line", line)
			continue
		}

		links := lineParts[1]
		g.addPageStrings(lineParts[0], strings.Split(links, ","))
	}

	return g
}
