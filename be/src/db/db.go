package db

import (
	"api/src/config"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func Conn() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.DbConnectionStr)

	if err != nil {
		fmt.Println("sql don't openned")
		return nil, err
	}

	if err = db.Ping(); err != nil {
		fmt.Println("doesn't ping")
		db.Close()
		return nil, err
	}

	return db, nil
}
