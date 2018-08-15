package main

import (
	"flag"
	"io/ioutil"
	"path"
	"strings"

	"github.com/inconshreveable/log15"
)

const (
	defaultFile = "wikidumps/simplewiki-20170820-pages-meta-current.xml"
	outputFile  = "graph.txt"
	readAmount  = 1000000

	fromPage = "cat"
	toPage   = "dog"
)

var inputFile = flag.String("infile", defaultFile, "Input XML dump file path")
var Log = log15.New()

func main() {
	flag.Parse()
	filePath := *inputFile
	ext := path.Ext(filePath)

	var g *Graph
	if ext == "xml" {
		g = ReadXMLFile(filePath)
	} else {
		// try to load using the graph loader instead. It should be much faster
		g = ReadGraphFile(filePath)
	}

	err := ioutil.WriteFile(outputFile, []byte(g.String()), 0644)
	if err != nil {
		panic(err)
	}

	path := searchPath(g, fromPage, toPage)

	Log.Debug("Found Path",
		"From", fromPage,
		"To", toPage,
		"Path", strings.Join(path, ","))
	Log.Debug("Total unique articles", "Articles", len(g.NodeMap))
}

// searchPath finds one possible path from pageA
// to pageB. Since the graph is directed, the inverse
// paths need not be the same
func searchPath(g *Graph, pageA, pageB string) []string {
	seen := make(map[*PageNode]struct{}, len(g.NodeMap)) // the seen nodes

	return searchPathBFS(g, g.GetNodeOrCreate(pageA), g.GetNodeOrCreate(pageB), []string{}, seen)
}

// BFS implmentation. Finds the shortest path
func searchPathBFS(g *Graph, pageA, pageB *PageNode, path []string, seen map[*PageNode]struct{}) []string {
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

// DFS implementation. Finds one path quickly
func searchPathDFS(g *Graph, pageA, pageB *PageNode, path []string, seen map[*PageNode]struct{}) []string {
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
		p := searchPathDFS(g, link, pageB, path, seen)
		if p != nil {
			return p
		}
	}

	return nil
}
