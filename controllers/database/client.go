package database

import "fmt"

type Connect struct {
	Host     string
	Port     string
	Password string
	Username string
	Database string
}

// GenDatabaseUrl returns database connection url
func (c *Connect) GenDatabaseUrl() string {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.Username, c.Password, c.Host, c.Port, c.Database)
	return databaseURL
}
