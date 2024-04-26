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
