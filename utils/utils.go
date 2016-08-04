package utils

import (
	"crypto/rand"
	mathRand "math/rand"
	"time"
)

func ToJSONTime(t time.Time) string {
	r, _ := t.MarshalJSON()
	return string(r[1 : len(r)-1])
}

func StrCpy(c string) *string {
	return &c
}

func BoolCpy(c bool) *bool {
	return &c
}

func IntCpy(c int) *int {
	return &c
}

func TimeCpy(c time.Time) *time.Time {
	return &c
}

func GenToken(strSize int) string {
	dictionary := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	var bytes = make([]byte, strSize)
	rand.Read(bytes)

	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}

	return string(bytes)
}

func RandInt(min int, max int) int {
	mathRand.Seed(time.Now().UTC().UnixNano())
	return min + mathRand.Intn(max-min)
}
