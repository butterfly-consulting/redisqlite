// +build cgo

package redisqlite

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"
import (
	"database/sql"
	"encoding/json"

	"github.com/wenerme/go-rm/rm"

	// sqlite database driver
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./sqlite.db")
	if err == nil {
		commands = append(commands,
			CreateCommandSQLEXEC(),
			CreateCommandSQL(),
		)
	}
}

// CreateCommandSQLEXEC is sql execute command
func CreateCommandSQLEXEC() rm.Command {
	return rm.Command{
		Usage:    "SQLEXEC sql",
		Desc:     `Execute a statement with SQLITE`,
		Name:     "sqlexec",
		Flags:    "readonly random no-cluster",
		FirstKey: 1, LastKey: 1, KeyStep: 1,
		Action: func(cmd rm.CmdContext) int {
			ctx, args := cmd.Ctx, cmd.Args
			if len(cmd.Args) != 2 {
				return ctx.WrongArity()
			}
			ctx.AutoMemory()
			// execute query
			sql := args[1].String()
			ctx.Log(rm.LOG_DEBUG, sql)
			_, err := db.Exec(sql)
			if err == nil {
				ctx.ReplyWithOK()
				return rm.OK
			} else {
				ctx.ReplyWithError(err.Error())
				return rm.ERR
			}
		},
	}
}

// CreateCommandSQL execute a sqlquery
func CreateCommandSQL() rm.Command {
	return rm.Command{
		Usage:    "SQL query",
		Desc:     `Execute a query with SQLITE`,
		Name:     "sql",
		Flags:    "readonly random no-cluster",
		FirstKey: 1, LastKey: 1, KeyStep: 1,
		Action: func(cmd rm.CmdContext) int {
			ctx, args := cmd.Ctx, cmd.Args
			if len(cmd.Args) != 2 {
				return ctx.WrongArity()
			}
			ctx.AutoMemory()
			sql := args[1].String()
			ctx.Log(rm.LOG_DEBUG, sql)

			// query the database
			rows, err := db.Query(sql)
			defer rows.Close()

			// output
			out := make([]map[string]interface{}, 0)
			columns, err := rows.Columns()
			if err != nil {
				ctx.ReplyWithError(err.Error())
				return rm.ERR
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
					ctx.ReplyWithError(err.Error())
					return rm.ERR
				}
				record := make(map[string]interface{})
				for i, v := range values {
					record[columns[i]] = v
				}
				out = append(out, record)
			}
			err = rows.Err()
			if err != nil {
				ctx.ReplyWithError(err.Error())
				return rm.ERR
			}
			bytes, err := json.Marshal(out)
			if err != nil {
				ctx.ReplyWithError(err.Error())
				return rm.ERR
			}
			res := string(bytes)

			ctx.Log(rm.LOG_DEBUG, res)
			ctx.ReplyWithSimpleString(res)
			return rm.OK
		},
	}
}
