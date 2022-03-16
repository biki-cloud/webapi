package random

import (
	"math/rand"
	"time"
)

// Choice choices random element of slice.
func Choice(l []string) string {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	return l[rand.Intn(len(l))]
}
