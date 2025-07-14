package helpers

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateReference() string {
	now := time.Now()
	formatted := now.Format("20060102150405.000")
	nowFormat := formatted[:14] + formatted[15:]
	randomNumber := rand.Intn(100)
	reference := fmt.Sprintf("%s%d", nowFormat, randomNumber)
	return reference
}
