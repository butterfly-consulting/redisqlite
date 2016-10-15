package rxhash

import "github.com/wenerme/go-rm/rm"

// open the key and make sure it is indeed a Hash and not empty
func openHashKey(ctx rm.Ctx, k rm.String) (rm.Key, bool) {
	key := ctx.OpenKey(k, rm.READ | rm.WRITE)
	if key.KeyType() != rm.KEYTYPE_EMPTY && key.KeyType() != rm.KEYTYPE_HASH {
		ctx.ReplyWithError(rm.ERRORMSG_WRONGTYPE)
		return rm.Key(0), false
	}
	return key, true
}
