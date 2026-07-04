package db

import (
	"database/sql"
	"fmt"

	"github.com/dprio/redis-limiter/internal/infrastructure/config"
	_ "github.com/go-sql-driver/mysql"
)

type DBs struct {
	OrderRepository OrderRepository
}

func New(dbConfig *config.DB) *DBs {
	db, err := sql.Open(dbConfig.Driver, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name))
	if err != nil {
		panic(err)
	}

	return &DBs{
		OrderRepository: NewOrderRepository(db),
	}
}
