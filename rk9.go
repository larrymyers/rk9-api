package rk9

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

type Event struct {
	Name        string
	Location    string
	StartDate   time.Time
	EndDate     time.Time
	URL         string
	DetailsURL  string
	PairingsURL string
}

var EventsURL = "https://rk9.gg/events/pokemon"

func GetEvents() ([]*Event, error) {
	resp, err := http.Get(EventsURL)
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

	rowSel, err := cascadia.Parse("table tbody tr")
	if err != nil {
		return nil, err
	}

	colSel, err := cascadia.Parse("td")
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	events := make([]*Event, 0)
	rows := cascadia.QueryAll(doc, rowSel)
	for _, row := range rows {
		cols := cascadia.QueryAll(row, colSel)

		if len(cols) != 5 {
			return nil, errors.New(fmt.Sprintf("expected 5 row columns, but got %d", len(cols)))
		}

		event := &Event{}

		n := cols[2].FirstChild
		for n != nil {
			if n.Data == "a" {
				break
			}

			n = n.NextSibling
		}
		text, href := parseAnchor(n)
		if text == "" {
			text = strings.TrimSpace(cols[2].FirstChild.Data)
		}

		event.Name = text
		event.URL = href

		// dateRange := strings.TrimSpace(cols[0].Data)
		event.Location = strings.TrimSpace(cols[3].FirstChild.Data)

		events = append(events, event)
	}

	return events, nil
}

func parseAnchor(node *html.Node) (string, string) {
	if node == nil || node.Data != "a" {
		return "", ""
	}

	href := ""
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			href = attr.Val
			continue
		}
	}

	text := strings.TrimSpace(node.FirstChild.Data)

	return text, href
}
