package redisqlite

import "C"
import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"

	// sqlite database driver
	_ "github.com/mattn/go-sqlite3"
)

// database access
var db *sql.DB

// prepared statement cache
const PREP_MAX = 1000000

var prep_index = 0
var prep_cache map[int]*sql.Stmt

// Open opens the sqlite database
func Open() (err error) {
	db, err = sql.Open("sqlite3", "./sqlite.db")
	prep_cache = make(map[int]*sql.Stmt)
	return err
}

// Exec execute a statement applying an array of arguments,
// returns the number of affected rows and the last id modified, when applicable
func Exec(stmtOrNumber string, args []interface{}) (count int64, lastId int64, err error) {
	var res sql.Result
	// select number or string
	if index, err := strconv.Atoi(stmtOrNumber); err == nil {
		stmt := prep_cache[index]
		if stmt == nil {
			return -1, -1, errors.New("no such prepared statement index")
		}
		res, err = stmt.Exec(args...)
		if err != nil {
			return -1, -1, err
		}
	} else {
		res, err = db.Exec(stmtOrNumber, args...)
		if err != nil {
			return -1, -1, err
		}
	}

	count, err = res.RowsAffected()
	if err != nil {
		count = -1
	}
	lastId, err = res.LastInsertId()
	if err != nil {
		lastId = -1
	}
	return lastId, count, nil
}

// Query execute a query applying an array of args
// query can be either an sql string or a number
// if it is a number then it will execute a prepared statement idenfied by the number returned by Prep
// returns an array of results, either as an array of maps or as an array of arrays
// according the `asMap` parameters, and returns up to `count` results (0 for everything)
func Query(queryOrNumber string, args []interface{}, asMap bool, count int64) (res []string, err error) {
	// execute a query
	var rows *sql.Rows

	// grab cached stmp
	if index, err := strconv.Atoi(queryOrNumber); err == nil {
		stmt := prep_cache[index]
		if stmt == nil {
			return nil, errors.New("no such prepared statement index")
		}
		rows, err = stmt.Query(args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
	} else {
		rows, err = db.Query(queryOrNumber, args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
	}

	// prepare output
	out := make([]string, 0)
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	ncol := len(columns)
	values := make([]interface{}, ncol)
	scanArgs := make([]interface{}, ncol)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// scan rows
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		var bytes []byte
		if asMap {
			// serialize as map
			record := make(map[string]interface{})
			for i, v := range values {
				record[columns[i]] = v
			}
			// serialize record
			bytes, err = json.Marshal(record)
			if err != nil {
				continue
			}
		} else {
			// serialize as array
			array := []interface{}{}
			for _, v := range values {
				array = append(array, v)
			}
			// serialize record
			bytes, err = json.Marshal(array)
			if err != nil {
				continue
			}
		}
		// collect
		out = append(out, string(bytes))

		// next until exhausted
		count--
		if count == 0 {
			break
		}
	}
	return out, rows.Err()
}

// Prep accepts prepares a sql statement and stores it in a table
// returning a number. It also accepts a number, and if it corresponds
// to the number returned by a previous statement, it closes the prepared stateent
// you can store up to one million statements, older will be automatically closed and dropper
func Prep(queryOrNumber string) (int, error) {
	if index, err := strconv.Atoi(queryOrNumber); err == nil {
		if index < PREP_MAX && prep_cache[index] != nil {
			prep_cache[index].Close()
			prep_cache[index] = nil
			return -1, nil
		} else {
			return -1, errors.New("invalid prepared statement index")
		}
	}

	// get next index and close very old statements if still unclosed
	prep_index = (prep_index + 1) % PREP_MAX
	if prep_cache[prep_index] != nil {
		prep_cache[prep_index].Close()
	}

	// prepare statement and store it
	stmt, err := db.Prepare(queryOrNumber)
	if err != nil {
		return -1, err
	}
	prep_cache[prep_index] = stmt
	return prep_index, nil
}
