package util

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid"
)

func ULID() string {
	t := time.Now()
	generatedUlid := ulid.MustNew(ulid.Timestamp(t), rand.Reader) // Let it panic on purpose

	return generatedUlid.String()
}
