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

## How to build and run locally 

Instructions for OSX, it works on any unix-like environment.

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
127.0.0.1:6379> sql "select * from t"
[{"i":1},{"i":2},{"i":3}]
127.0.0.1:6379>
```

## Build and start a docker image

You need docker up and running, of course.

Use:

```
make image
make imagestart
```
