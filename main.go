package main

import (
	"database/sql"
	"log"

	"github.com/jingtian168/admin-api-go/api"
	db "github.com/jingtian168/admin-api-go/db/sqlc"
	"github.com/jingtian168/admin-api-go/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadCconfig("/Users/libin/Downloads/project/soybean-admin-go/")
	if err != nil {
		log.Fatal("connect load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("connot start server:", err)
	}
}
