package util

import (
	"bytes"
	"context"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func TextWithoutSpaces(_ context.Context, selector *goquery.Selection) string {
	var buf bytes.Buffer

	// Slightly optimized vs calling Each: no single selection object created
	// Copied from how goquery handles finding of raw text except the trimming of \n and \t
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			// Remove \n and \t unlike jQuery
			trimmedStr := strings.Trim(n.Data, "\n\t")
			buf.WriteString(trimmedStr + " ")
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	for _, n := range selector.Nodes {
		f(n)
	}

	return buf.String()
}

func GetLinksFromSelection(_ context.Context, selector *goquery.Selection) []string {
	var links []string
	selector.Find("a").Each(func(i int, linkSelector *goquery.Selection) {
		link, exists := linkSelector.Attr("href")
		if exists {
			links = append(links, link)
		}
	})
	return links
}
