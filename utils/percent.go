package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateCardNumber() int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	var result string
	for i := 0; i < 15; i++ {
		result += strconv.Itoa(r.Intn(10))
	}
	n, _ := strconv.Atoi(result)
	return n
}
func GenerateExpirationDate() string {
	currentDate := time.Now()
	expirationDate := currentDate.AddDate(6, 0, 0)
	return expirationDate.Format("01.06")
}
func ChooseRandomCard() string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	if r.Intn(2) == 0 {
		return "VISA"
	}
	return "MASTERCARD"
}
