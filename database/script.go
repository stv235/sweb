package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
)

func Run(db *sql.DB, path string) {
	buf, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}

	tx, err := db.Begin()

	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			if err := tx.Rollback(); err != nil {
				panic(err)
			}

			panic(r)
		}
	}()

	/*for i, q := range strings.Split(string(buf), ";") {
		q = strings.TrimSpace(q)

		if _, err := tx.Exec(q); err != nil {
			panic(errors.New(err.Error() + " \"" + q + "\" [" + strconv.Itoa(1 + i) + "]"))
		}
	}*/



	if _, err := tx.Exec(string(buf)); err != nil {
		panic(errors.New(fmt.Sprintf("In %v error %v", path, err.Error())))
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
