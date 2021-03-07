// +build cgo

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Init initialize db
func Init() (err error) {
	db, err = sql.Open("sqlite3", "./simple.db")
	return err
}

// Exec executes a query - returns the error message or nil if all ok
func Exec(query string, args ...interface{}) error {
	_, err := db.Exec(query)
	return err
}

// Query query string
func Query(query string) (res string, err error) {
	rows, err := db.Query(query)
	defer rows.Close()

	out := make([]map[string]interface{}, 0)

	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return "", err
		}
		record := make(map[string]interface{})
		for i, v := range values {
			record[columns[i]] = v
		}
		out = append(out, record)
	}
	bytes, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	err = rows.Err()
	return string(bytes), err
}

func main() {
	err := Init()
	if err != nil {
		log.Fatal(err)
	}

	// ignore error
	Exec("drop table person")

	fmt.Println("Creating table")
	sqlStmt := `
	create table person (id int not null primary key, name varchar(20), height real, data blob);
	delete from person;
	`
	err = Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	fmt.Println(err)

	fmt.Println("SQL error:", Exec("blabla"))

	fmt.Println("Transaction")
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into person(id, name, height, data) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 5; i++ {
		f := 1 + 0.01*float64(i)
		_, err = stmt.Exec(i, fmt.Sprintf("user%03d", i), f, "xxxxx")
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	fmt.Println("Query")
	res, err := Query("select * from person")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(res)
	}

	/*
		fmt.Println("Prepared query")
		stmt, err = db.Prepare("select name from person where id = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		var name string
		err = stmt.QueryRow("3").Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(name)

		fmt.Println("Query All")
		_, err = db.Exec("delete from person")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
		if err != nil {
			log.Fatal(err)
		}

		res, err = Query("select id, name from foo")
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(res)
		}
	*/
}
