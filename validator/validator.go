package validator

import (
	"regexp"
)

const (
	alphaRegStr       = `^[a-zA-Z]+$`
	alphaNumberRegStr = `^[a-zA-Z0-9]+$`
	numericRegStr     = `^[-+]?[0-9]+(?:\.[0-9]+)?$`
	numberRegStr      = `^[0-9]+$`
	emailRegStr       = `^[a-z0-9][.a-z0-9_-]*@[a-z0-9][a-z0-9-]*(\.[a-z0-9]{2,10})+$`
	mobileRegStr      = `^1[3-9][0-9]{9}$`
)

var (
	items = make(map[string]*regexp.Regexp)
)

func testStr(key, ptn, str string) bool {
	reg, bl := items[key]
	if !bl {
		reg = regexp.MustCompile(ptn)
		items[key] = reg
	}
	return reg.MatchString(str)
}

func IsNumber(str string) bool {
	return testStr("IsNumber", numberRegStr, str)
}

func IsNumeric(str string) bool {
	return testStr("IsNumeric", numericRegStr, str)
}

func IsAlpha(str string) bool {
	return testStr("IsAlpha", alphaRegStr, str)
}

func IsAlphaNum(str string) bool {
	return testStr("IsAlphaNum", alphaNumberRegStr, str)
}

func IsMobile(str string) bool {
	return testStr("IsMobile", mobileRegStr, str)
}

func IsEmail(str string) bool {
	return testStr("IsEmail", emailRegStr, str)
}

func StrLen(str string, minLength, maxLength int) bool {
	strLen := len(str)
	return minLength <= strLen && strLen <= maxLength
}
