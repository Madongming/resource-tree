package tools

import (
	"math/rand"
	"time"

	"github.com/satori/go.uuid"
)

func GetUuid() string {
RETRY_UUID:
	id, err := uuid.NewV4()
	if err != nil {
		goto RETRY_UUID
	}
	return id.String()
}

func GetRandInt() int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return r.Int()
}
