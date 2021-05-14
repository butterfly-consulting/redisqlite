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
	"github.com/wenerme/go-rm/rm"

	// sqlite database driver
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	if Open() == nil {
		commands = append(commands,
			CreateCommandSQLEXEC(),
			CreateCommandSQLMAP(),
			CreateCommandSQLARR(),
			//CreateCommandSQLPREP(),
			//CreateCommandSQLBEGIN(),
			//CreateCommandSQLCOMMIT(),
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

			// collect args
			params := make([]interface{}, 0)
			i := 2
			for i < len(args) {
				params = append(params, args[i].String())
				i++
			}
			count, lastId, err := Exec(sql, params)
			if err == nil {
				ctx.ReplyWithArray(2)
				ctx.ReplyWithLongLong(count)
				ctx.ReplyWithLongLong(lastId)
				return rm.OK
			} else {
				ctx.ReplyWithError(err.Error())
				return rm.ERR
			}
		},
	}
}

// CreateCommandSQL execute a sqlquery
func CreateCommandSQLMAP() rm.Command {
	return rm.Command{
		Usage:    "SQLMAP query count",
		Desc:     `Execute a query with SQLITE returning a list of maps`,
		Name:     "sqlmap",
		Flags:    "readonly random no-cluster",
		FirstKey: 1, LastKey: 1, KeyStep: 1,
		Action: func(cmd rm.CmdContext) int {
			ctx, args := cmd.Ctx, cmd.Args
			if len(cmd.Args) <= 3 {
				return ctx.WrongArity()
			}
			ctx.AutoMemory()
			var count int64
			if args[1].StringToLongLong(&count) == rm.ERR {
				ctx.ReplyWithError("invalid count")
				return rm.ERR
			}
			sql := args[2].String()
			ctx.Log(rm.LOG_DEBUG, sql)

			// collect args
			params := make([]interface{}, 0)
			i := 3
			for i < len(args) {
				params = append(params, args[i].String())
				i++
			}
			// query the database
			res, err := Query(sql, params, true, count)
			if err != nil {
				ctx.ReplyWithError(err.Error())
				return rm.ERR
			}
			ctx.ReplyWithArray(int64(len(res)))
			for _, v := range res {
				ctx.ReplyWithSimpleString(v)
			}
			return rm.OK
		},
	}
}

// CreateCommandSQL execute a sqlquery
func CreateCommandSQLARR() rm.Command {
	return rm.Command{
		Usage:    "SQLLST query count",
		Desc:     `Execute a query with SQLITE returning an array of json arrays`,
		Name:     "sqlarr",
		Flags:    "readonly random no-cluster",
		FirstKey: 1, LastKey: 1, KeyStep: 1,
		Action: func(cmd rm.CmdContext) int {
			ctx, args := cmd.Ctx, cmd.Args
			if len(cmd.Args) <= 3 {
				return ctx.WrongArity()
			}
			ctx.AutoMemory()
			var count int64
			if args[1].StringToLongLong(&count) == rm.ERR {
				ctx.ReplyWithError("invalid count")
				return rm.ERR
			}
			sql := args[2].String()
			ctx.Log(rm.LOG_DEBUG, sql)

			// collect args
			params := make([]interface{}, 0)
			i := 3
			for i < len(args) {
				params = append(params, args[i].String())
				i++
			}
			// query the database
			res, err := Query(sql, params, false, count)
			if err != nil {
				ctx.ReplyWithError(err.Error())
				return rm.ERR
			}
			ctx.ReplyWithArray(int64(len(res)))
			for _, v := range res {
				ctx.ReplyWithSimpleString(v)
			}
			return rm.OK
		},
	}
}
