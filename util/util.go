package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	rnd "math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var (
	chars         map[string]string
	ptnCamelCase  = regexp.MustCompile(`_([a-z0-9])`)
	ptnUnderScore = regexp.MustCompile(`([A-Z])`)
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

func Struct2Map(obj interface{}) map[string]interface{} {
	types := reflect.TypeOf(obj)
	values := reflect.ValueOf(obj)

	if types.Kind() == reflect.Ptr {
		types = types.Elem()
		values = values.Elem()
	}

	size := types.NumField()
	var data = make(map[string]interface{})
	for i := 0; i < size; i++ {
		name := types.Field(i).Name
		name = strings.ToLower(name[0:1]) + name[1:]
		data[name] = values.Field(i).Interface()
	}
	return data
}

func ToCamelCase(str string) string {
	return ptnCamelCase.ReplaceAllStringFunc(str, func(a string) string {
		return strings.Title(a[1:2])
	})
}

func ToUnderScore(str string) string {
	str = ptnUnderScore.ReplaceAllStringFunc(str, func(a string) string {
		fmt.Println(a)
		return "_" + strings.ToLower(a)
	})
	return strings.Trim(str, "_")
}

func ToLowFirst(str string) string {
	return strings.ToLower(str[0:1]) + str[1:]
}
