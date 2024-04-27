package api

import (
	"context"
	"time"

	"cloud.google.com/go/civil"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"larrymyers.com/rk9api/rk9"
)

type Client struct {
	conn *pgxpool.Pool
}

func NewClient(conn *pgxpool.Pool) *Client {
	return &Client{conn: conn}
}

func (c *Client) GetEvents(ctx context.Context) ([]*rk9.Event, error) {
	query := `
	SELECT
	  id,
		name,
		location,
		start_date,
		end_date,
		url  
	FROM events
	ORDER BY start_date DESC;
	`

	rows, err := c.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var startDate time.Time
	var endDate time.Time

	events, err := pgx.CollectRows[*rk9.Event](rows, func(row pgx.CollectableRow) (*rk9.Event, error) {
		var event rk9.Event
		err := row.Scan(&event.ID, &event.Name, &event.Location, &startDate, &endDate, &event.URL)
		if err != nil {
			return nil, err
		}

		event.StartDate = civil.DateOf(startDate)
		event.EndDate = civil.DateOf(endDate)

		return &event, nil
	})

	return events, err
}

func (c *Client) UpsertEvent(ctx context.Context, event *rk9.Event) error {
	upsert := `
	INSERT INTO events (id, name, location, start_date, end_date, url)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (id) DO UPDATE SET name = $2, location = $3, start_date = $4, end_date = $5, url = $6;
	`

	_, err := c.conn.Exec(ctx, upsert, event.ID, event.Name, event.Location,
		event.StartDate, event.EndDate, event.URL)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpsertMatch(ctx context.Context, match *rk9.Match) error {
	upsertPlayer := `
	INSERT INTO players (id, name, country, wins, losses, ties, points, decklist_url, standing)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (id) DO UPDATE SET name = $2, country = $3, wins = $4, losses = $5, ties = $6, points = $7, decklist_url = $8, standing = $9;
	`

	upsertMatch := `
	INSERT INTO matches (id, pod, round_number, table_number, player1_id, player2_id, winner_id, is_tie)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (id) DO UPDATE SET player1_id = $5, player2_id = $6, winner_id = $7, is_tie = $8;
	`

	p1 := match.Player1
	p2 := match.Player2

	if p1 != nil {
		_, err := c.conn.Exec(ctx, upsertPlayer, p1.ID(), p1.Name, p1.Country, p1.Wins, p1.Losses, p1.Ties, p1.Points, p1.DecklistURL, p1.Standing)
		if err != nil {
			return err
		}
	}

	if p2 != nil {
		_, err := c.conn.Exec(ctx, upsertPlayer, p2.ID(), p2.Name, p2.Country, p2.Wins, p2.Losses, p2.Ties, p2.Points, p2.DecklistURL, p2.Standing)
		if err != nil {
			return err
		}
	}

	_, err := c.conn.Exec(ctx, upsertMatch, match.ID(), match.Pod, match.Round, match.Table, match.Player1.ID(), match.Player2.ID(), match.Winner.ID(), match.IsTie)
	if err != nil {
		return err
	}

	return nil
}
