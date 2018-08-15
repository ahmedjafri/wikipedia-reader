package main

import (
	"strings"
	"sync"
)

type Graph struct {
	sync.Mutex
	NodeMap map[string]*PageNode
}

type PageNode struct {
	Title string
	Links []*PageNode
}

// Convenient constructer to init the map
func NewGraph() *Graph {
	var g Graph
	g.NodeMap = make(map[string]*PageNode, 0)
	return &g
}

func (g *Graph) AddXMLPage(pxml XMLPage) {
	g.Lock()
	defer g.Unlock()

	stringLinks := pxml.links()
	pn := g.GetNodeOrCreate(pxml.Title())
	pn.Links = make([]*PageNode, len(stringLinks))
	// Add links to
	for i, l := range stringLinks {
		pn.Links[i] = g.GetNodeOrCreate(l)

	}
}

func (g *Graph) GetNodeOrCreate(title string) *PageNode {
	if pn, ok := g.NodeMap[title]; ok {
		return pn
	}

	pn := &PageNode{Title: title}
	g.NodeMap[title] = pn
	return pn
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
