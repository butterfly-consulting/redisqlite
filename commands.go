// +build cgo

/*
#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
*/

package redisqlite

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
