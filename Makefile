redisqlite.so:
	cd main && go build -v -buildmode=c-shared -o ../redisqlite.so
	chmod +x  redisqlite.so

.PHONY: clean
clean:
	rm redisqlite.so

.PHONY: start
start: redisqlite.so
	redis-server --loadmodule ./redisqlite.so --loglevel debug
