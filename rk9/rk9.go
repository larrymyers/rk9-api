package rk9

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func getPage(pageURL string) (*html.Node, error) {
	reqURL, err := url.Parse(pageURL)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(reqURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("%d: %s", resp.StatusCode, body))
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return doc, nil
}
