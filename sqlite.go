package redisqlite

import "C"
import (
	"database/sql"
	"encoding/json"

	// sqlite database driver
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Open() (err error) {
	db, err = sql.Open("sqlite3", "./sqlite.db")
	return err
}

func Exec(sql string) (err error) {
	_, err = db.Exec(sql)
	return err
}

func Query(sql string, count int64) (string, error) {
	rows, err := db.Query(sql)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// output
	out := make([]map[string]interface{}, 0)
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	ncol := len(columns)
	values := make([]interface{}, ncol)
	scanArgs := make([]interface{}, ncol)
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
		// next
		count--
		if count == 0 {
			break
		}
	}
	err = rows.Err()
	if err != nil {
		return "", err
	}
	bytes, err := json.Marshal(out)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
