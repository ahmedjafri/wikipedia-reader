package main

import (
	"encoding/xml"
	"os"
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
// The second word cannot contain numbers
const pageLinkRegexConstraints = "[a-zA-Z0-9]*\\s*[a-zA-Z]*"

var pageLinkRegex = regexp.MustCompile("\\[\\[(?P<link>" + pageLinkRegexConstraints + ")\\]\\]")
var pageTitleRegex = regexp.MustCompile(pageLinkRegexConstraints)

func ReadXMLFile(filename string) *Graph {
	xmlFile, err := os.Open(filename)
	if err != nil {
		Log.Error("Error opening file", "Error", err)
		return nil
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

				// skip pages whose titles dont match the constraints
				if !pageTitleRegex.MatchString(p.Title_) {
					continue
				}

				g.AddXMLPage(p)

				total++
				if total > readAmount-1 {
					break mainLoop
				}

			}
		default:
		}

	}

	Log.Debug("Total read articles", "Articles", total)
	return g
}

func (p XMLPage) Title() string {
	return strings.ToLower(p.Title_)
}

func (p XMLPage) links() []string {
	r := pageLinkRegex.FindAllStringSubmatch(p.Text, linkLimit)
	if len(r) >= linkLimit {
		Log.Crit("Find all links matched the limit amount. It could be possible that we have clipped some",
			"Page", p.Title_)
	}

	rt := make([]string, len(r))
	for i, s := range r {
		// first group in the regex
		rt[i] = s[1]
	}

	return dedup(rt, false)
}
