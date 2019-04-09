package encoding

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/gentwolf-shen/gohelper/convert"
	"github.com/gentwolf-shen/gohelper/hashhelper"
)

func Decode(str, key string, expires int64) (string, bool) {
	return to("DECODE", str, key, expires)
}

func Encode(str, key string, expires int64) (string, bool) {
	return to("ENCODE", str, key, expires)
}

func to(mode, str, key string, expires int64) (string, bool) {
	if mode == "DECODE" {
		str = strings.Replace(str, " ", "+", -1)
	}

	keyA := hashhelper.Md5(key[0:16])
	keyB := hashhelper.Md5(key[16:])
	keyC := ""
	if mode == "DECODE" {
		keyC = str[0:4]
	} else {
		keyC = convert.ToStr(time.Now().Unix())[0:4]
	}

	cryptKey := keyA + hashhelper.Md5(keyA+keyC)
	cryptKeyLength := 64

	var b []byte
	if mode == "DECODE" {
		str = str[4:]
		b, _ = base64.RawStdEncoding.DecodeString(str)
	} else {
		str1 := hashhelper.Md5(str + keyB)
		str = convert.ToStr(time.Now().Unix()+expires) + str1[0:16] + str

		b = []byte(str)
	}

	strLength := len(b)

	rndKey := make([]uint8, 256)
	box := make([]uint8, 256)
	for i := 0; i <= 255; i++ {
		rndKey[i] = cryptKey[i%cryptKeyLength]
		box[i] = uint8(i)
	}

	i := 0
	j := 0
	for i = 0; i <= 255; i++ {
		j = (j + int(box[i]) + int(rndKey[i])) % 256
		box[i], box[j] = box[j], box[i]
	}

	i = 0
	j = 0
	k := 0
	rs := make([]uint8, strLength)
	for i = 0; i < strLength; i++ {
		k = (k + 1) % 256
		j = (j + int(box[k])) % 256
		box[k], box[j] = box[j], box[k]
		rs[i] = b[i] ^ box[int((box[k]+box[j]))%256]
	}

	result := ""
	bl := false

	if mode == "DECODE" {
		rawStr := string(rs)
		bl = true
		if expires > 1 {
			t := convert.ToInt64(rawStr[0:10])
			if time.Now().Unix() < t {
				bl = true
			}
		}
		if bl && len(rawStr) > 26 {
			result = string(rawStr[26:])
		} else {
			bl = false
		}
	} else {
		str := base64.RawStdEncoding.EncodeToString(rs)
		result = keyC + strings.Replace(str, "=", "", -1)
		bl = true
	}

	return result, bl
}
