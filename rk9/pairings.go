package rk9

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

const (
	Juniors = 0
	Seniors = 1
	Masters = 2
)

type Match struct {
	Player1 *Player
	Player2 *Player
	Table   string
	Winner  *Player
}

type Player struct {
	Name    string
	Country string
	Wins    int
	Losses  int
	Ties    int
	Points  int
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
		match := Match{}
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

			player := Player{
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

		points, err := strconv.Atoi(string(matches[3]))
		if err != nil {
			return 0, 0, 0, 0, err
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