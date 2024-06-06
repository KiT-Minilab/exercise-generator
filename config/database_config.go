package config

import (
	"fmt"
	"net/url"
)

type DBConfig interface {
	String() string
	DSN() string
}

type DatabaseConfig struct {
	Host     string `json:"host" yaml:"host"`
	Database string `json:"database" yaml:"database"`
	Port     int    `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Options  string `json:"options" yaml:"options"`
}

func (c DatabaseConfig) DSN() string {
	options := c.Options
	if options != "" {
		if options[0] != '?' {
			options = "?" + options
		}
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s",
		c.Username,
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.Database,
		options)
}

type PostgreSQLConfig struct {
	DatabaseConfig
}

func (c PostgreSQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@%s:%d/%s%s", c.Username, url.QueryEscape(c.Password), c.Host, c.Port, c.Database, c.Options)
}

func (c PostgreSQLConfig) String() string {
	return fmt.Sprintf("postgres://%s", c.DSN())
}

func PostgresSQLDefaultConfig() PostgreSQLConfig {
	return PostgreSQLConfig{DatabaseConfig{
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "exercise_generator",
		Username: "kiet13",
		Password: "123456",
		Options:  "?sslmode=disable",
	}}
}
