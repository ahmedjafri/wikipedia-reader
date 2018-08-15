package main

import (
	"fmt"
	"strings"
	"sync"
)

type Graph struct {
	sync.Mutex
	NodeMap map[string]*PageNode
}

type PageNode struct {
	Title string
	Links []PageNode
}

func (g *Graph) AddXMLPage(pxml XMLPage) {
	g.Lock()
	defer g.Unlock()

	links := pxml.links()
	fmt.Printf("Page Title: %s\nLinks: %#v\n", pxml.Title, links)
}

func (p PageNode) String() string {
	links := make([]string, len(p.Links))
	for i, lp := range p.Links {
		links[i] = lp.Title
	}

	return p.Title + "=" + strings.Join(links, ",")
}

// Serialize to string
// The format is:
// page1=page2,page3,page4
// page2=page3,page1
func (g *Graph) String() string {
	var s strings.Builder
	for _, v := range g.NodeMap {
		// WriteString returns a nil error so no need to check it
		s.WriteString(v.String())
		s.WriteByte('\n')
	}

	return s.String()
}
