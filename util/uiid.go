package util

import (
	"math/rand"
	"time"
)

func GenerateUUID() string {
	chars := []rune("0123456789abcdef")
	uuid := make([]rune, 36)

	uuid[8] = '-'
	uuid[13] = '-'
	uuid[18] = '-'
	uuid[23] = '-'
	uuid[14] = '4'

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 36; i++ {
		if uuid[i] == 0 {
			r := rune(0)

			if i == 19 {
				r = rune((rnd.Intn(4) | 8))
			} else {
				r = rune(rnd.Intn(16))
			}

			uuid[i] = chars[r]
		}
	}

	return string(uuid)
}

