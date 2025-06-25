package utils

import (
	"math/rand"
	"time"
)

func GetRandom(options []string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len(options))
	return options[randomIndex]
}
