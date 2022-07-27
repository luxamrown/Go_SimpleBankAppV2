package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GenerateAccountNumber() string {
	timeStamp, timeYear := GetTime()
	c1, c2 := GetTwoRandStr()
	accountNumber := fmt.Sprintf("%s%s%s%s", c1, timeStamp, timeYear, c2)
	return accountNumber
}

func GetTime() (string, string) {
	timeStamp := time.Now().Format("15:04:05")
	timeYear := time.Now().Format("2006")
	timeStamp = strings.ReplaceAll(timeStamp, ":", "")
	timeYear = timeYear[len(timeYear)-2:]
	return timeStamp, timeYear
}

func GetTwoRandStr() (string, string) {
	charset1 := "12345"
	charset2 := "6789"
	c1 := charset1[rand.Intn(len(charset1))]
	c2 := charset2[rand.Intn(len(charset2))]
	return string(c1), string(c2)
}
