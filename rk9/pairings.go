package rk9

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/andybalholm/cascadia"
)

const (
	Juniors = 0
	Seniors = 1
	Masters = 2
)

type EventRounds struct {
	Juniors int
	Seniors int
	Masters int
}

type Match struct {
	EventID string
	Pod     int
	Round   int
	Table   string
	Player1 *EventPlayer
	Player2 *EventPlayer
	Winner  *EventPlayer
	IsTie   bool
}

func (m *Match) ID() string {
	return fmt.Sprintf("%d-%d-%s", m.Pod, m.Round, m.Table)
}

type EventPlayer struct {
	Name        string
	Country     string
	Wins        int
	Losses      int
	Ties        int
	Points      int
	DecklistURL string
	Standing    int
}

func (p *EventPlayer) ID() string {
	if p == nil {
		return ""
	}

	return fmt.Sprintf("%s-%s", p.Name, p.Country)
}

func GetRounds(event *Event) (EventRounds, error) {
	rounds := EventRounds{}

	doc, err := getPage(BaseURL + event.PairingsURL())
	if err != nil {
		return rounds, err
	}

	navLinkSel, err := cascadia.Parse(".nav-tabs .nav-link")
	if err != nil {
		return rounds, err
	}

	podRoundPattern := regexp.MustCompile(`P(?P<pod>\d)R(?P<round>\d+)`)

	navLinks := cascadia.QueryAll(doc, navLinkSel)

	// TODO do we need to iterate if we assume the links are ascending?
	// why not just grab the last one that isn't the standings link?
	for _, navLink := range navLinks {
		href := attrVal(navLink, "href")
		href = strings.TrimLeft(href, "#")
		matches := podRoundPattern.FindStringSubmatch(href)

		if len(matches) == 3 {
			pod := matches[1]
			round := matches[2]

			roundNum, err := strconv.Atoi(round)
			if err != nil {
				return rounds, err
			}

			switch pod {
			case "0":
				rounds.Juniors = roundNum
			case "1":
				rounds.Seniors = roundNum
			case "2":
				rounds.Masters = roundNum
			}
		}
	}

	return rounds, nil
}

func GetRound(event *Event, pod int, round int) ([]*Match, error) {
	reqURL, err := url.Parse(BaseURL + event.PairingsURL())
	if err != nil {
		return nil, err
	}

	vals := reqURL.Query()
	vals.Add("rnd", strconv.Itoa(round))
	vals.Add("pod", strconv.Itoa(pod))
	reqURL.RawQuery = vals.Encode()

	doc, err := getPage(reqURL.String())
	if err != nil {
		return nil, err
	}

	matchSel, err := cascadia.Parse(".match")
	if err != nil {
		return nil, err
	}

	playersSel, err := cascadia.Parse(".player")
	if err != nil {
		return nil, err
	}

	playerNameSel, err := cascadia.Parse(".name")
	if err != nil {
		return nil, err
	}

	tableSel, err := cascadia.Parse(".tablenumber")
	if err != nil {
		return nil, err
	}

	matches := make([]*Match, 0)

	matchesEl := cascadia.QueryAll(doc, matchSel)
	for _, m := range matchesEl {
		match := Match{
			EventID: event.ID,
			Pod:     pod,
			Round:   round,
		}

		table := cascadia.Query(m, tableSel)
		match.Table = innerText(table)

		players := cascadia.QueryAll(m, playersSel)
		for _, p := range players {
			record := innerText(p)
			wins, losses, ties, points, err := parseRecord(record)
			if err != nil {
				return nil, err
			}

			name, country := parseName(innerText(cascadia.Query(p, playerNameSel)))

			// TODO find a way to safely match name to FirstName + LastName and Country

			player := EventPlayer{
				Name:    name,
				Country: country,
				Wins:    wins,
				Losses:  losses,
				Ties:    ties,
				Points:  points,
			}

			if hasClass(p, "player1") {
				match.Player1 = &player
			} else if hasClass(p, "player2") {
				match.Player2 = &player
			}

			if hasClass(p, "winner") {
				match.Winner = &player
			}

			if hasClass(p, "tie") {
				match.IsTie = true
			}
		}

		matches = append(matches, &match)
	}

	return matches, nil
}

var recordExp = regexp.MustCompile(`\d+`)

func parseRecord(record string) (int, int, int, int, error) {
	matches := recordExp.FindAll([]byte(record), -1)
	if matches != nil {
		wins, err := strconv.Atoi(string(matches[0]))
		if err != nil {
			return 0, 0, 0, 0, err
		}

		losses, err := strconv.Atoi(string(matches[1]))
		if err != nil {
			return 0, 0, 0, 0, err
		}

		ties, err := strconv.Atoi(string(matches[2]))
		if err != nil {
			return 0, 0, 0, 0, err
		}

		// when the Top 8 cut is reached points are not reported in the record
		points := 0
		if len(matches) > 4 {
			points, err = strconv.Atoi(string(matches[3]))
			if err != nil {
				return 0, 0, 0, 0, err
			}
		}

		return wins, losses, ties, points, nil
	}

	return 0, 0, 0, 0, nil
}

func parseName(name string) (string, string) {
	parts := strings.Split(name, " ")
	playerName := []string{}
	country := ""

	for _, p := range parts {
		if strings.HasPrefix(p, "[") {
			country = strings.Trim(p, "[]")
		} else {
			playerName = append(playerName, p)
		}
	}
	return strings.Join(playerName, " "), country
}
