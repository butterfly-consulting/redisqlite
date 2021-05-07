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

import (
	"github.com/wenerme/go-rm/rm"
)

var commands []rm.Command
var dataTypes []rm.DataType

func CreateModule() *rm.Module {
	mod := rm.NewMod()
	mod.Name = "redisqlite"
	mod.Version = 1
	mod.SemVer = "0.1-alpha"
	mod.Author = "sciabarracom"
	mod.Website = "http://github.com/sciabarracom/redisqlite"
	mod.Desc = `This module implements sqlite embedding in redis`
	mod.Commands = commands
	mod.DataTypes = dataTypes
	return mod
}
