package util

import (
	"crypto/hmac"
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GetAuthKey(deviceUUID string) string {
	encKey := uuid.New().String()
	authKey := fmt.Sprintf("%x", md5.Sum([]byte(encKey)))

	today := ""
	currentTime := time.Now()

	day := strconv.Itoa(currentTime.Day())
	month := strconv.Itoa(int(currentTime.Month()))
	year := strconv.Itoa(currentTime.Year())

	if len(day) == 1 {
		day = "0" + day
	}
	if len(month) == 1 {
		month = "0" + month
	}
	today = day + month + year

	getKey := fmt.Sprintf("%x", md5.Sum([]byte(strings.ToLower(encKey))))
	uiidKey := fmt.Sprintf("%x", md5.Sum([]byte(strings.ToLower(deviceUUID))))

	setKey := fmt.Sprintf("%x", hmacMD5([]byte(strings.ToLower(getKey)), []byte(strings.ToLower(uiidKey))))

	authKey = strings.ToUpper(getKey) +
		strings.ToUpper(fmt.Sprintf("%x", hmacMD5([]byte(setKey), []byte(deviceUUID)))) +
		strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(today))))

	return authKey
}

func hmacMD5(key, data []byte) []byte {
	hash := hmac.New(md5.New, key)
	hash.Write(data)
	finalSum := hash.Sum(nil)
	return finalSum
}
