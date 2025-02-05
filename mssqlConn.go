package main

import (
	"database/sql"
	"fmt"
	"net/url"
)

type DBConfig struct {
	Server   string
	Port     int
	Database string
	AppName  string
}

type DBConn struct {
	db *sql.DB
}

func NewDBConn(config DBConfig) (*DBConn, error) {
	query := url.Values{}
	query.Add("app name", config.AppName)
	query.Add("trusted_connection", "yes")
	query.Add("TrustServerCertificate", "true")
	query.Add("database", config.Database)

	u := &url.URL{
		Scheme: "sqlserver",
		// User:   url.UserPassword("LPESCSVCDIQST", "3+5Nfdrp"),
		User:     url.UserPassword("azrsqlqst", "nM&9V!!JXX#g2vK@&Y$d"),
		Host:     fmt.Sprintf("%s:%d", config.Server, config.Port),
		RawQuery: query.Encode(),
	}

	db, err := sql.Open("sqlserver", u.String())
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	return &DBConn{db: db}, nil
}

func (c *DBConn) Close() error {
	return c.db.Close()
}

func (c *DBConn) Ping() (string, error) {

	if err := c.db.Ping(); err != nil {
		c.db.Close()
		return "", fmt.Errorf("database ping error: %w", err)
	}

	return "Ping Ok", nil
}

func (c *DBConn) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return c.db.Query(query, args...)
}
