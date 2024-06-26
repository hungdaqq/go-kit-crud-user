package db

import (
    // "database/sql"
    "log"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

var db *sqlx.DB

func InitDB(dataSourceName string) {
    var err error
    db, err = sqlx.Open("postgres", dataSourceName)
    if err != nil {
        log.Fatal(err)
    }

    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }

    log.Println("Connected to PostgreSQL!")
}

func GetDB() *sqlx.DB {
    return db
}
