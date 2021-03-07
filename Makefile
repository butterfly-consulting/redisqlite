
redis: build
	redis-server --loadmodule ./redisqlite --loglevel debug
