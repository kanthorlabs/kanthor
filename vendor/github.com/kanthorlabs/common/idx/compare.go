package idx

import (
	"time"

	"github.com/segmentio/ksuid"
)

var (
	Differ = time.Second * 10
)

// AfterTime uses SafeUnixDiff as factor to make sure we can get an id that is always less than the given time
func BeforeTime(t time.Time) string {
	id, err := ksuid.FromParts(t.Add(-Differ), []byte("0000000000000000"))
	if err != nil {
		panic(err)
	}
	return id.Prev().String()
}

// AfterTime uses SafeUnixDiff as factor to make sure we can get an id that is always greater than the given time
func AfterTime(t time.Time) string {
	id, err := ksuid.FromParts(t.Add(Differ), []byte("0000000000000000"))
	if err != nil {
		panic(err)
	}
	return id.Next().Next().String()
}
