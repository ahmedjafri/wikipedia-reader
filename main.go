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
)

var inputFile = flag.String("infile", defaultFile, "Input XML dump file path")
var from = flag.String("from", "guitar", "Wikipedia page to start the search at")
var to = flag.String("to", "laundry", "Wikipedia page to search for")

var Log = log15.New()

func main() {
	flag.Parse()
	filePath := *inputFile
	fromPage := *from
	toPage := *to

	g := readFile(filePath)
	path := SearchPathBFS(g, fromPage, toPage)

	Log.Debug("Found Path", "From", fromPage, "To", toPage, "Path", strings.Join(path, ","))
	Log.Debug("Total unique articles", "Articles", len(g.NodeMap))
}

func readFile(filePath string) *Graph {
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

	return g
}
