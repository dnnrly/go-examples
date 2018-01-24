package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang/geo/r2"
	_ "github.com/mattn/go-sqlite3"
)

// This is shamelessly lifted from as a starting point:
// https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go
func main() {
	os.Remove("./foo.db")

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, data text);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, data) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(1, r2.Point{X: 1.0, Y: 2.0})
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}
