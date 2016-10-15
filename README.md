# Redis hash type extension commands

rxhash is a redis module build on top of [go-rm](https://github.com/wenerme/go-rm/) implement a lot of missing commands for Redis [hash type](http://redis.io/commands/#hash).  

This module provides following commands

* `hgetset key field value`
    * Return old value
* `hgetdel key field`
    * Return removed value
* `hsetm key field old new`
    * Set when value match old
* `hsetex key field value`
    * Set when field exists
* `hdelm key field old new`
    * Delete when value match old

## Install

> Ensure you use latest redis build from source.

```bash
# Build module from source
go build -v -buildmode=c-shared github.com/redismodule/rxhash/cmd/rxhash
# Load module
redis-server --loadmodule ./rxhash --loglevel debug
# You can use these commands now.
```

## Test

After redis-server started you can run the test script

```bash
cd ~/go/src/github.com/redismodule/rxhash
./test.sh
```

If test failed will output something like 

```
$ hset a a 2
(integer) 1
FAILED: Expected (integer)
```
