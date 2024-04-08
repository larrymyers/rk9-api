package api

import (
	"context"

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
	FROM events;
	`

	rows, err := c.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var events []*rk9.Event
	for rows.Next() {
		var event rk9.Event
		rows.Scan(&event.ID, &event.Name, &event.StartDate, &event.EndDate, &event.URL)
	}

	return events, nil
}

func (c *Client) UpsertEvent(ctx context.Context, event *rk9.Event) error {
	upsert := `
	INSERT INTO events (id, name, location, start_date, end_date, url)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (id) UPDATE SET name = $2, location = $3, start_date = $4, end_date = $5, url = $6;
	`

	_, err := c.conn.Exec(ctx, upsert, event.ID, event.Name, event.Location,
		event.StartDate, event.EndDate, event.Location)
	if err != nil {
		return err
	}

	return nil
}
