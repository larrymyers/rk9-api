package rk9

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

type Event struct {
	ID        string
	Name      string
	Location  string
	StartDate time.Time
	EndDate   time.Time
	URL       string
}

func (evt *Event) DetailsURL() string {
	return fmt.Sprintf("/tournament/%s", evt.ID)
}

func (evt *Event) PairingsURL() string {
	return fmt.Sprintf("/pairings/%s", evt.ID)
}

func (evt *Event) RosterURL() string {
	return fmt.Sprintf("/roster/%s", evt.ID)
}

func (evt *Event) HasStarted() bool {
	return evt.StartDate.Before(time.Now())
}

func (evt *Event) HasEnded() bool {
	return evt.EndDate.Before(time.Now())
}

var EventsURL = "/events/pokemon"

func GetEvents() ([]*Event, error) {
	doc, err := getPage(BaseURL + EventsURL)
	if err != nil {
		return nil, err
	}

	rowSel, err := cascadia.Parse("table tbody tr")
	if err != nil {
		return nil, err
	}

	colSel, err := cascadia.Parse("td")
	if err != nil {
		return nil, err
	}

	linkSel, err := cascadia.Parse("a")
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

		start, end, err := parseDateRange(strings.TrimSpace(cols[0].FirstChild.Data))
		if err != nil {
			return nil, err
		}

		event.StartDate = start
		event.EndDate = end

		event.Location = strings.TrimSpace(cols[3].FirstChild.Data)

		links := cascadia.QueryAll(cols[4], linkSel)
		for _, link := range links {
			text, href := parseAnchor(link)
			if text == "TCG" {
				parts := strings.Split(href, "/")
				event.ID = parts[len(parts)-1]
				break
			}
		}

		events = append(events, event)
	}

	return events, nil
}

func parseAnchor(node *html.Node) (string, string) {
	if node == nil || node.Data != "a" {
		return "", ""
	}

	href := attrVal(node, "href")
	text := innerText(node)

	return text, href
}

func parseDateRange(s string) (time.Time, time.Time, error) {
	s = strings.ReplaceAll(s, "–", "-")
	parts := strings.Split(s, " ")

	// same month: April 1-10, 2023
	if len(parts) == 3 {
		month := parts[0]
		days := parts[1]
		year := parts[2]

		days = strings.Trim(days, ",")
		dayParts := strings.Split(days, "-")
		startDay := dayParts[0]
		endDay := dayParts[1]

		start, err := time.Parse("January 2 2006", fmt.Sprintf("%s %s %s", month, startDay, year))
		if err != nil {
			return time.Now(), time.Now(), err
		}

		end, err := time.Parse("January 2 2006", fmt.Sprintf("%s %s %s", month, endDay, year))
		if err != nil {
			return time.Now(), time.Now(), err
		}

		return start, end, nil
	}

	// multiple months: January 4-July 18, 2024
	if len(parts) == 4 {
		parts = strings.Split(s, ",")
		year := parts[1]
		rest := parts[0]
		parts = strings.Split(rest, "-")

		start, err := time.Parse("January 2 2006", fmt.Sprintf("%s %s", parts[0], year))
		if err != nil {
			return time.Now(), time.Now(), err
		}

		end, err := time.Parse("January 2 2006", fmt.Sprintf("%s %s", parts[1], year))
		if err != nil {
			return time.Now(), time.Now(), err
		}

		return start, end, nil
	}

	return time.Now(), time.Now(), errors.New(fmt.Sprintf("unrecognized date range: %s", s))
}
