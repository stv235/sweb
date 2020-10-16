package database

import (
	"database/sql"
	"errors"
)

var ErrNotInDb = errors.New("not a db object (already in transaction?)")
var ErrNotInTransaction = errors.New("not in a transaction")


func Begin(q Q) Q {
	db, ok := q.(*sql.DB)

	if !ok {
		panic(ErrNotInDb)
	}

	tx, err := db.Begin()

	if err != nil {
		panic(err)
	}

	return tx
}

func End(q Q) {
	tx, ok := q.(*sql.Tx)

	if !ok {
		panic(ErrNotInTransaction)
	}

	if err := recover(); err != nil {
		tx.Rollback()
		panic(err)
	} else {
		if err := tx.Commit(); err != nil {
			panic(err)
		}
	}
}
