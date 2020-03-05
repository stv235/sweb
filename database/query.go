package database

import (
	"database/sql"
	"log"
)

func Query(db Q, str string, args... interface{}) *sql.Rows {
	res, err := db.Query(str, args...)

	if err != nil {
		log.Output(2, err.Error())
		panic(err)
	}

	return res
}
