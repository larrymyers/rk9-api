package rk9

import (
	"strings"

	"golang.org/x/net/html"
)

var BaseURL = "https://rk9.gg"

func innerText(n *html.Node) string {
	if n == nil || n.FirstChild == nil {
		return ""
	}

	n = n.FirstChild
	text := ""

	for n != nil {
		if n.Type == html.TextNode {
			text += n.Data
		}

		n = n.NextSibling
	}

	return strings.TrimSpace(text)
}

func attrVal(n *html.Node, key string) string {
	if n == nil {
		return ""
	}

	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}

	return ""
}

func hasClass(n *html.Node, className string) bool {
	if n == nil {
		return false
	}

	class := attrVal(n, "class")
	for _, name := range strings.Split(class, " ") {
		if name == className {
			return true
		}
	}

	return false
}
