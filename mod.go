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
	mod.Desc = `This module implements sqlite embedding`
	mod.Commands = commands
	mod.DataTypes = dataTypes
	return mod
}
