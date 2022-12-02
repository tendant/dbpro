package dbpro

import (
	"math/rand"
	"time"
)

func RandomBit() int {
	return rand.Intn(2)
}

func Now() string {
	return time.Now().Format(time.RFC3339)
}
