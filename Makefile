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

redisqlite.so: *.go
	cd main && go build -v -buildmode=c-shared -o ../redisqlite.so
	chmod +x redisqlite.so

.PHONY: clean
clean:
	rm redisqlite.so

.PHONY: start
start: redisqlite.so
	redis-server --loadmodule ./redisqlite.so --loglevel debug

.PHONY: image
image:
	docker build . -t redisqlite

.PHONY: imagestart
imagestart: image
	docker run -p 6379:6379 -d --rm --name redisqlite redisqlite --requirepass password
	docker logs redisqlite

.PHONY: imagestop
imagestop:
	docker kill redisqlite
