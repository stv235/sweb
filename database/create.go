package database

import (
	"database/sql"
	"errors"
	sqlite "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var ErrCreateScriptNotFound = errors.New("database create script not found")
var ErrEmptyDriver = errors.New("empty driver")
var ErrEmptyDataSourceName = errors.New("empty data source name")

func Create(config Config) *sql.DB {
	if config.Driver == "" {
		log.Panicln(ErrEmptyDriver)
	}

	if config.DataSourceName == "" {
		log.Panicln(ErrEmptyDataSourceName)
	}

	sql.Register("sweb_sqlite3", &sqlite.SQLiteDriver{
		ConnectHook: func(conn *sqlite.SQLiteConn) error {
			if err := conn.RegisterFunc("match_or", matchOr, true); err != nil {
				return err
			}

			return nil
		},
	})

	db, err := sql.Open("sweb_sqlite3", config.DataSourceName)

	if err != nil {
		log.Panicln(err)
	}

	switch {
	case config.CreateScript != "":
		if _, err := os.Stat(config.CreateScript); os.IsNotExist(err) {
			panic(ErrCreateScriptNotFound)
		}

		Run(db, config.CreateScript)
	case config.CreateQuery != "":
		if _, err := db.Exec(config.CreateQuery); err != nil {
			panic(err)
		}
	}

	return db
}
