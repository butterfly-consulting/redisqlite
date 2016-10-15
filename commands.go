package rxhash

import (
	"github.com/wenerme/go-rm/rm"
)

func init() {
	commands = append(commands,
		CreateCommand_HGETSET(),
		CreateCommand_HGETDEL(),
		CreateCommand_HSETM(),
		CreateCommand_HDELM(),
		CreateCommand_HSETEX(),
	)
}

func CreateCommand_HGETSET() rm.Command {
	return rm.Command{
		Usage: "HGETSET key field value",
		Desc: `Sets the 'field' in Hash 'key' to 'value' and returns the previous value, if any.
Reply: String, the previous value or NULL if 'field' didn't exist. `,
		Name:   "hgetset",
		Flags:  "write fast deny-oom",
		FirstKey:1, LastKey:1, KeyStep:1,
		Action: func(cmd rm.CmdContext) int {
			ctx, args := cmd.Ctx, cmd.Args
			if len(cmd.Args) != 4 {
				return ctx.WrongArity()
			}
			ctx.AutoMemory()
			key, ok := openHashKey(ctx, args[1])
			if !ok {
				return rm.ERR
			}
			// get the current value of the hash element
			var val rm.String;
			key.HashGet(rm.HASH_NONE, cmd.Args[2], (*uintptr)(&val))
			// set the element to the new value
			key.HashSet(rm.HASH_NONE, cmd.Args[2], cmd.Args[3])
			if val.IsNull() {
				ctx.ReplyWithNull()
			} else {
				ctx.ReplyWithString(val)
			}
			return rm.OK
		},
	}
}
func CreateCommand_HGETDEL() rm.Command {
	return rm.Command{
		Usage: "HGETDEL key field",
		Desc: `Delete field and return value`,
		Name:   "hgetdel",
		Flags:  "write fast deny-oom",
		FirstKey:1, LastKey:1, KeyStep:1,
		Action: func(cmd rm.CmdContext) int {
			ctx, args := cmd.Ctx, cmd.Args
			if len(cmd.Args) != 3 {
				return ctx.WrongArity()
			}
			ctx.AutoMemory()
			key, ok := openHashKey(ctx, args[1])
			if !ok {
				return rm.ERR
			}
			// get the current value of the hash element
			var val rm.String;
			key.HashGet(rm.HASH_NONE, cmd.Args[2], (*uintptr)(&val))
			if val.IsNull() {
				ctx.ReplyWithNull()
			} else {
				key.HashDel(args[2])
				ctx.ReplyWithString(val)
			}
			return rm.OK
		},
	}
}

func CreateCommand_HSETM() rm.Command {
	return rm.Command{
		Usage: "HSETM key field old-value new-value",
		Desc: "Set when value match old",
		Name:   "hsetm",
		Flags:  "write fast deny-oom",
		FirstKey:1, LastKey:1, KeyStep:1,
		Action: func(cmd rm.CmdContext) int {
			ctx, args := cmd.Ctx, cmd.Args
			if len(cmd.Args) != 5 {
				return ctx.WrongArity()
			}
			ctx.AutoMemory()
			key, ok := openHashKey(ctx, args[1])
			if !ok {
				return rm.ERR
			}
			// get the current value of the hash element
			var val rm.String;
			key.HashGet(rm.HASH_NONE, args[2], &val)

			if val.IsNull() || val.Compare(args[3]) != 0 {
				ctx.ReplyWithLongLong(0)
			} else {
				// set the element to the new value
				key.HashSet(rm.HASH_NONE, args[2], args[4])
				ctx.ReplyWithLongLong(1)
			}
			return rm.OK
		},
	}
}

func CreateCommand_HDELM() rm.Command {
	return rm.Command{
		Usage: "HDELM key field old-value",
		Desc: "Delete when value match old",
		Name:   "hdelm",
		Flags:  "write fast deny-oom",
		FirstKey:1, LastKey:1, KeyStep:1,
		Action: func(cmd rm.CmdContext) int {
			ctx, args := cmd.Ctx, cmd.Args
			if len(cmd.Args) != 4 {
				return ctx.WrongArity()
			}
			ctx.AutoMemory()

			key, ok := openHashKey(ctx, args[1])
			if !ok {
				return rm.ERR
			}
			// get the current value of the hash element
			var val rm.String;
			key.HashGet(rm.HASH_NONE, args[2], &val)

			if val.IsNull() || val.Compare(args[3]) != 0 {
				ctx.ReplyWithLongLong(0)
			} else {
				ctx.ReplyWithLongLong(int64(key.HashDel(args[2])))
			}
			return rm.OK
		},
	}
}
func CreateCommand_HSETEX() rm.Command {
	return rm.Command{
		Usage: "HSETEX key field value",
		Desc: "Set field to value ony if field is already exists",
		Name:   "hsetex",
		Flags:  "write fast deny-oom",
		FirstKey:1, LastKey:1, KeyStep:1,
		Action: func(cmd rm.CmdContext) int {
			ctx, args := cmd.Ctx, cmd.Args
			if len(args) != 4 {
				return ctx.WrongArity()
			}
			ctx.AutoMemory()
			key, ok := openHashKey(ctx, args[1])
			if !ok {
				return rm.ERR
			}
			ctx.ReplyWithLongLong(int64(key.HashSet(rm.HASH_XX, args[2], args[3])))
			return rm.OK
		},
	}
}
