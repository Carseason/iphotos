package iphotos

import (
	"hash/fnv"
)

const (
	alphabetString = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

func GenFnvID(v string) uint64 {
	f := fnv.New64a()
	f.Write([]byte(v))
	return f.Sum64()
}
