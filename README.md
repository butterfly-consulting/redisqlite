<!--
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
-->
# Redisqlite

A fresh implementation of the sqlite module for redis.

Truly open source, no string attached, no phone-home requirements, no licensing per instance/hour and so on.

# Installation

## How to build and run locally 

Instructions for OSX, it should work on any other unix-like environment with Go and Redis.

You need go 1.15, install it (OSX: `brew install go`)

You need redis of course, install it  (OSX: `brew install redis`)

Then run `make start`

On another terminal connect to redis with `redis-cli`  and try:

```
$ redis-cli
127.0.0.1:6379> sqlexec "create table t(i int)"
OK
127.0.0.1:6379> sqlexec "insert into t(i) values(1),(2),(3)"
OK
127.0.0.1:6379> sqlmap 0 "select * from t"
[{"i":1},{"i":2},{"i":3}]
127.0.0.1:6379> sqlmap 1 "select * from t"
[{"i":1}]
127.0.0.1:6379>
```

## Build and start a docker image

You need docker up and running, of course.

Use:

```
make image
make imagestart
```

It will build and start a local image called `redisqlite` with redis including the module.

# Command Reference

- `SQLPREP (<query>|<id>)`:<br> 
  if the argument is not numeric, it is sql to prepare a statement. 
  Returns a numeric `<id>` that identifies the statement in queryes.
  If the argument is numeric it is assumed to be a previously prepared statement that is then closed.

- `SQLEXEC (<statement>|<id>) [<args> ...]`: <br>
   execute a sql statement, either a prepared one or an sql statement.
   If the argument is numeric it is assumed to be an `<id>` of a prepared statement, otherwise it assumed to be an SQL statements.
   It returns the number of rows affected when relevant or -1, and the last id generated when relevant or -1. 

- `SQLMAP <limit> (<query>|<id>) [<args> ...]`:<br>
   execute a SQL query, either a prepared one or a sql query.
   If the argument is number it is assumed to be an `<id>` of a prepared statement, otherwiser it is assumed to be an SQL qyery.
   It returns an array of string representing a JSON maps.
   It limits the number of elements in the returned arrays by  `<limit>`; 0 means unlimited.
   Each object corresponds to a record, where each key corresponds to field names and each value is the corresponding field value.
 
- `SQLARR <limit> <query>|<number>) [<args> ...]`:<br>
   execute a SQL query, either a prepared one or a sql query.
   If the argument is number it is assumed to be an `<id>` of a prepared statement, otherwiser it is assumed to be an SQL qyery.
   It returns an array of string representing  JSON arrays.
   It limits the number of elements in the returned arrays by  `<limit>`; 0 means unlimited.
   Each array represents a record, with the values of the fields in the array.

