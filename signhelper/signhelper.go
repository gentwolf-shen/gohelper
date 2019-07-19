package signhelper

import (
	"sort"
	"strings"

	"github.com/gentwolf-shen/gohelper/hashhelper"
)

const (
	MD5 = iota
	SHA1
	SHA256
)

func GetSignSimple(params map[string]string, secret string, signType int) string {
	return GetSign(params, secret, signType, "", "")
}

func GetSign(params map[string]string, secret string, signType int, joiner, separator string) string {
	str := BuildQuery(params, joiner, separator) + secret
	signStr := ""

	switch signType {
	case MD5:
		signStr = hashhelper.Md5(str)
	case SHA1:
		signStr = hashhelper.Sha1(str)
	case SHA256:
		signStr = hashhelper.Sha256(str)
	}

	return signStr
}

func BuildQuery(params map[string]string, joiner, separator string) string {
	paramsLength := len(params)
	keys := make([]string, paramsLength)
	i := 0
	for k := range params {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	tmp := make([]string, paramsLength)
	for i, k := range keys {
		tmp[i] = k + joiner + params[k]
	}

	return strings.Join(tmp, separator)
}
