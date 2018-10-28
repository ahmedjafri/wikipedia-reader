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

	// used by BFS
	Path []string
}

// Convenient constructer to init the map
func NewGraph() *Graph {
	var g Graph
	g.NodeMap = make(map[string]*PageNode, 0)
	return &g
}

func (g *Graph) AddXMLPage(pxml XMLPage) {
	g.addPageStrings(pxml.Title(), pxml.links())
}

func (g *Graph) addPageStrings(pageTitle string, pageLinks []string) {
	g.Lock()
	defer g.Unlock()

	stringLinks := pageLinks
	pn := g.GetNodeOrCreate(pageTitle)
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

//
// ***** Search *******
//

// searchPath finds one possible path from pageA
// to pageB. Since the graph is directed, the inverse
// paths need not be the same
func SearchPathBFS(g *Graph, pageA, pageB string) []string {
	seen := make(map[*PageNode]struct{}, len(g.NodeMap)) // the seen nodes

	return searchPathBFSi(g, g.GetNodeOrCreate(pageA), g.GetNodeOrCreate(pageB), []string{}, seen)
}

func SearchPathDFS(g *Graph, pageA, pageB string) []string {
	seen := make(map[*PageNode]struct{}, len(g.NodeMap)) // the seen nodes

	return searchPathDFSi(g, g.GetNodeOrCreate(pageA), g.GetNodeOrCreate(pageB), []string{}, seen)
}

// DFS implementation. Finds one path quickly
func searchPathDFSi(g *Graph, pageA, pageB *PageNode, path []string, seen map[*PageNode]struct{}) []string {
	if pageA == nil {
		return nil
	}

	if pageA == pageB {
		path = append(path, pageB.Title)
		return path
	}

	if _, ok := seen[pageA]; ok {
		// already seen this node
		return nil
	}

	seen[pageA] = struct{}{}
	path = append(path, pageA.Title)

	for _, link := range pageA.Links {
		p := searchPathDFSi(g, link, pageB, path, seen)
		if p != nil {
			return p
		}
	}

	return nil
}

// BFS implmentation. Finds the shortest path
func searchPathBFSi(g *Graph, pageA, pageB *PageNode, path []string, seen map[*PageNode]struct{}) []string {
	var q []*PageNode
	q = append(q, pageA)

	for len(q) > 0 {
		// pop
		n := q[0]
		q = q[1:len(q)]
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		path := append(n.Path, n.Title)

		if n == pageB {
			return path
		}

		links := make([]*PageNode, len(n.Links))
		copy(links, n.Links)
		for i, _ := range links {
			links[i].Path = path
		}

		q = append(q, n.Links...)
	}

	return nil
}
