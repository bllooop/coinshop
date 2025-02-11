package main

import (
	running "github.com/bllooop/coinshop/internal/server"

	_ "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	running.Run()
}
