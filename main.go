package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
)

const (
	defaultFile = "wikidumps/simplewiki-20170820-pages-meta-current.xml"
	readAmount  = 1
)

var inputFile = flag.String("infile", defaultFile, "Input XML dump file path")

func main() {
	flag.Parse()

	xmlFile, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	var total int64
	var inElement string
	g := NewGraph()

mainLoop:
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			inElement = se.Name.Local
			// ...and its name is "page"
			if inElement == "page" {
				var p XMLPage
				// decode a whole chunk of following XML into the
				// variable p which is a Page
				decoder.DecodeElement(&p, &se)
				g.AddXMLPage(p)

				total++
				if total > readAmount-1 {
					break mainLoop
				}

			}
		default:
		}

	}

	fmt.Printf("Graph:\n%s", g.String())
	fmt.Printf("Total read articles: %d \n", total)
}
