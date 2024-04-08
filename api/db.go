package api

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	pgxslog "github.com/mcosta74/pgx-slog"
)

type ConnectionBuilder struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
	CertPath string
	Logger   *slog.Logger
}

func NewConnectionBuilder() *ConnectionBuilder {
	return &ConnectionBuilder{
		Host:    "localhost",
		Port:    5432,
		SSLMode: "disable",
	}
}

func (c *ConnectionBuilder) WithEnvVars() *ConnectionBuilder {
	if host := os.Getenv("DB_HOST"); len(host) > 0 {
		c.Host = host
	}

	if port := os.Getenv("DB_PORT"); len(port) > 0 {
		p, err := strconv.Atoi(port)
		if err != nil {
			panic(err)
		}

		c.Port = p
	}

	if user := os.Getenv("DB_USER"); len(user) > 0 {
		c.User = user
	}

	if pass := os.Getenv("DB_PASS"); len(pass) > 0 {
		c.Password = pass
	}

	if name := os.Getenv("DB_NAME"); len(name) > 0 {
		c.Name = name
	}

	return c
}

func (c *ConnectionBuilder) WithLogger(l *slog.Logger) *ConnectionBuilder {
	c.Logger = l

	return c
}

func (c *ConnectionBuilder) ConnectionString() string {
	connURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:   c.Name,
	}

	q := connURL.Query()

	if len(c.SSLMode) > 0 {
		q.Set("sslmode", c.SSLMode)
	}

	if len(c.CertPath) > 0 {
		q.Set("sslrootcert", c.CertPath)
	}

	connURL.RawQuery = q.Encode()

	return connURL.String()
}

func (c *ConnectionBuilder) GetConnection() (*pgxpool.Pool, error) {
	connURL := c.ConnectionString()
	pgConf, err := pgxpool.ParseConfig(connURL)
	if err != nil {
		return nil, err
	}

	if c.Logger != nil {
		adapter := pgxslog.NewLogger(c.Logger)
		pgConf.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   adapter,
			LogLevel: tracelog.LogLevelDebug,
		}
	}

	return pgxpool.NewWithConfig(context.Background(), pgConf)
}

//go:embed sql/schema.sql
var schema string

func ApplySchema(ctx context.Context, conn *pgxpool.Pool) error {
	_, err := conn.Exec(ctx, schema)
	if err != nil {
		return err
	}

	return nil
}
