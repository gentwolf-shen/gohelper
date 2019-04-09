package util

import (
	"crypto/rand"
	"encoding/hex"
	"math"
	rnd "math/rand"
	"strings"
	"time"
)

var (
	chars map[string]string
)

func init() {
	chars = make(map[string]string, 5)
	chars["&"] = "&amp;"
	chars["\""] = "&quot;"
	chars["'"] = "&#039;"
	chars[">"] = "&gt;"
	chars["<"] = "&lt;"
}

func RndStr(length int) string {
	r := rnd.New(rnd.NewSource(time.Now().UnixNano()))
	rs := make([]string, length)
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := 0; i < length; i++ {
		index := r.Intn(62)

		rs[i] = string(str[index])
	}
	return strings.Join(rs, "")
}

func SubString(str string, start, length int) string {
	rs := []rune(str)
	strLen := len(rs)

	if start >= strLen {
		return ""
	}

	if start < 0 {
		start += strLen
	}

	if length < 0 {
		length = strLen + length - 1
	}

	end := start + length
	if end > strLen {
		end = strLen
	}

	s := rs[start:end]
	return string(s)
}

func Ceil(size, count int32) int32 {
	return int32(math.Ceil(float64(count) / float64(size)))
}

func Uuid() string {
	u := make([]byte, 16)
	if _, err := rand.Read(u); err != nil {
		panic(err)
	}
	u[6] = (u[6] & 0x0f) | (4 << 4)
	u[8] = (u[8] & 0xbf) | 0x80

	return hex.EncodeToString(u)
}

func Rnd(min, max int) int {
	max += 1
	r := rnd.New(rnd.NewSource(time.Now().UnixNano()))
	return min + r.Intn(max-min)
}

func FilterHtmlChars(str string) string {
	for k, v := range chars {
		str = strings.Replace(str, k, v, -1)
	}
	return str
}
