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
FROM redis:5.0.3 as builder
RUN apt-get update &&\
    apt-get -y install curl gcc &&\
    curl -L https://golang.org/dl/go1.16.4.linux-amd64.tar.gz | tar xzvf - -C /usr
ENV PATH=/bin:/usr/bin:/usr/go/bin
COPY *.go go.* build/
COPY main/* build/main/
RUN  cd build/main &&\ 
     GO11MODULE=off go build -v -buildmode=c-shared -o /lib/redisqlite.so &&\
	chmod +x /lib/redisqlite.so
FROM redis:5.0.3
COPY --from=builder /lib/redisqlite.so /lib/redisqlite.so
RUN echo "loadmodule /lib/redisqlite.so" >/etc/redis.conf
ENTRYPOINT ["redis-server", "/etc/redis.conf"]
