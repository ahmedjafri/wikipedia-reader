package main

import (
	"regexp"
	"strings"
)

type XMLPage struct {
	Title_ string `xml:"title"`
	Redir  struct {
		Title string `xml:"title,attr"`
	} `xml:"redirect"`
	Text string `xml:"revision>text"`
}

const linkLimit = 1000

// Limit links to only have two words
const pageLinkRegexConstraints = "(?P<link>[a-zA-Z0-9]*\\s*[a-zA-Z0-9]*)"

var pageLinkRegex = regexp.MustCompile("\\[\\[" + pageLinkRegexConstraints + "\\]\\]")

func (p XMLPage) Title() string {
	return strings.ToLower(p.Title_)
}

func (p XMLPage) links() []string {
	r := pageLinkRegex.FindAllStringSubmatch(p.Text, linkLimit)
	if len(r) >= linkLimit {
		panic("Find all links matched the limit amount. It could be possible that we have clipped some")
	}

	rt := make([]string, len(r))
	for i, s := range r {
		// first group in the regex
		rt[i] = s[1]
	}

	return dedup(rt, false)
}
