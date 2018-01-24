package main

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"os"

	"github.com/golang/geo/r2"
	_ "github.com/mattn/go-sqlite3"
)

// Location allows us to implement the Value() function for writing r2.Point to the database
type Location struct {
	r2.Point
}

// Value does the conversion to a string
func (p *Location) Value() (driver.Value, error) {
	data := p.String()
	log.Printf("Converted r2.Point to a string: %s", data)
	return data, nil
}

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
	_, err = stmt.Exec(1, &Location{r2.Point{X: 1.0, Y: 2.0}})
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()

	rows, err := db.Query("select id, data from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var data Location
		err = rows.Scan(&id, &data)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Here's the data that we've read from the database: %s\n", data)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
