package database

import (
	"database/sql"
	"fmt"
	"log"
)

func SetDatabaseVersion(db *sql.DB, version int) {
	_, err := db.Exec(fmt.Sprintf("PRAGMA user_version = %v", version))

	if err != nil {
		log.Panicln(err)
	}
}

func GetDatabaseVersion(db *sql.DB) int {
	res := db.QueryRow("PRAGMA user_version")

	var version int

	if err := res.Scan(&version); err != nil {
		log.Panicln(err)
	}

	return version
}

func UpdateDatabase(db *sql.DB, updateFn func(db *sql.DB, version int) int) {
	userVersion := GetDatabaseVersion(db)

	for {
		log.Println("[DB]", "upgrading from version", userVersion)

		nextVersion := updateFn(db, userVersion)

		if nextVersion == userVersion {
			break
		}

		log.Println("[DB]", "upgraded from schema version", userVersion, "to", nextVersion)

		userVersion = nextVersion
	}

	SetDatabaseVersion(db, userVersion)
}