package main

import (
	"database/sql"
	"products/internal/connector/restserver/router"
	"products/pkg/restserver"
)

const (
	applicationName = "products"
	dbMigrationPath = "internal/connector/db/migrations"
)

type app struct {
	database *sql.DB
	server   *restserver.Server
	router   *router.Router
}
