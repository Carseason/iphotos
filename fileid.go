package iphotos

import (
	"hash/fnv"

	"github.com/sqids/sqids-go"
)

const (
	alphabetString = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

var (
	sid *sqids.Sqids
)

func init() {
	sid, _ = sqids.New(sqids.Options{
		Alphabet:  alphabetString,
		MinLength: 16,
	})

}

func GenFnvID(v string) uint64 {
	f := fnv.New64a()
	f.Write([]byte(v))
	return f.Sum64()
}

func GenFileID(vs ...string) (string, error) {
	n := len(vs)
	fids := make([]uint64, 0, n)
	for i := 0; i < n; i++ {
		fids = append(fids, GenFnvID(vs[i]))
	}
	return sid.Encode(fids)
}
