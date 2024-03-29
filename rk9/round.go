package rk9

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

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
	Table   int
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
	reqURL, err := url.Parse(BaseURL + event.PairingsURL)
	if err != nil {
		return nil, err
	}

	reqURL.Query().Add("rnd", strconv.Itoa(round))
	reqURL.Query().Add("pod", strconv.Itoa(pod))

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

	rowSel, err := cascadia.Parse("div.row")
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	rows := cascadia.QueryAll(doc, rowSel)

	for _, row := range rows {
		log.Println(row.Data)
	}

	return nil, nil
}
