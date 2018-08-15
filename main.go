package main

import (
	"flag"
	"io/ioutil"
	"path"

	"github.com/inconshreveable/log15"
)

const (
	defaultFile = "wikidumps/simplewiki-20170820-pages-meta-current.xml"
	outputFile  = "graph.txt"
	readAmount  = 1000000
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

	Log.Debug("Total unique articles", "Articles", len(g.NodeMap))
}
