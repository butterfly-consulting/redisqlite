# Redisqlite

A fresh implementation of the sqlite module for redis.

Truly open source, no string attached, no phone-home requirements, no licensing per instance/hour and so on.

*currently tested only on OSX!*

You need go 1.15, install it with `brew install go`

You need redis of course, install it with `brew install redis`

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
