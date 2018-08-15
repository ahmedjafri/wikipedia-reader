package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
)

const defaultFile = "wikidumps/simplewiki-20170820-pages-meta-current.xml"

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
	readAmount := int64(1)
	var total int64
	var inElement string

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
				searchPage(p)

				total++
				if total > readAmount-1 {
					fmt.Printf("Parsed %d pages\n", total)
					break mainLoop
				}

			}
		default:
		}

	}

	fmt.Printf("Total articles: %d \n", total)
}

func searchPage(p XMLPage) {
	links := p.links()
	fmt.Printf("Page Title: %s\nLinks: %#v\n", p.Title, links)
}
