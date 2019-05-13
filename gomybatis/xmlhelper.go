package gomybatis

import (
	"encoding/xml"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/gentwolf-shen/gohelper/logger"
)

var (
	ptnOp    = regexp.MustCompile(`[\s]*(.*?)[\s]*(!=|==|>=|<=|>|<|=)[\s]*(.*)[\s]*`)
	ptnAndOr = regexp.MustCompile(`^(AND|and|OR|or) `)
	ptnTest  = regexp.MustCompile(`[\s]*(AND|and|OR|or)[\s]+`)
)

func parseXmlFromStr(str string) map[string]SqlItem {
	rs, err := parseXml([]byte(str))
	if err != nil {
		logger.Error("parse xml error:" + str)
		panic(err)
	}

	return rs
}

func parseXmlFromFile(filename string) map[string]SqlItem {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("read file error: " + filename)
		panic(err)
	}

	rs, err := parseXml(b)
	if err != nil {
		logger.Error("parse xml file error:" + filename)
		panic(err)
	}

	return rs
}

func parseXml(b []byte) (map[string]SqlItem, error) {
	mapper := &Mapper{}
	if err := xml.Unmarshal(b, mapper); err != nil {
		return nil, err
	}

	result := make(map[string]SqlItem)

	addResult := func(items []SqlItem) {
		for _, item := range items {
			item.Sql = strings.TrimSpace(item.Sql)
			result[item.Id] = item
		}
	}

	addResult(mapper.Selects)
	addResult(mapper.Updates)
	addResult(mapper.Deletes)
	addResult(mapper.Inserts)

	return result, nil
}

func buildSelect(sqlItem *SqlItem, args map[string]interface{}) string {
	return parseFilter(sqlItem.Sql + buildSqlWhere(sqlItem.Where.Ifs, args) + " " + strings.TrimSpace(sqlItem.Suffix))
}

func buildUpdate(sqlItem *SqlItem, args map[string]interface{}) string {
	return parseFilter(sqlItem.Sql + buildSqlSet(sqlItem.Set.Ifs, args) + " " + buildSqlWhere(sqlItem.Where.Ifs, args) + " " + strings.TrimSpace(sqlItem.Suffix))
}

func buildDelete(sqlItem *SqlItem, args map[string]interface{}) string {
	return parseFilter(sqlItem.Sql + buildSqlWhere(sqlItem.Where.Ifs, args) + " " + strings.TrimSpace(sqlItem.Suffix))
}

func buildInsert(sqlItem *SqlItem, args map[string]interface{}) string {
	return sqlItem.Sql
}

func buildSqlWhere(ifs []ItemIf, args map[string]interface{}) string {
	length := len(ifs)
	arr := make([]string, length)

	index := 0
	for i := 0; i < length; i++ {
		if parseTest(ifs[i].Test, args) {
			tmp := strings.TrimSpace(ifs[i].If)
			if ptnAndOr.MatchString(tmp) {
				arr[index] = tmp
			} else {
				arr[index] = "AND " + tmp
			}

			index++
		}
	}

	if index > 0 {
		tmp := strings.Join(arr[:index], " ")
		return " WHERE " + strings.Trim(strings.Trim(tmp, "OR "), "AND ")
	}

	return ""
}

func buildSqlSet(ifs []ItemIf, args map[string]interface{}) string {
	length := len(ifs)
	arr := make([]string, length)

	index := 0
	for i := 0; i < length; i++ {
		if parseTest(ifs[i].Test, args) {
			arr[index] = strings.TrimSpace(ifs[i].If)
			index++
		}
	}

	if index > 0 {
		return " SET " + strings.Join(arr[:index], ", ")
	}

	return ""
}

func parseTest(str string, args map[string]interface{}) bool {
	if str == "" {
		return true
	}

	bl := false
	str = ptnTest.ReplaceAllStringFunc(str, strings.ToUpper)
	ands := strings.Split(str, " AND ")
	for i := range ands {
		if strings.Contains(ands[i], " OR ") {
			if !bl {
				ors := strings.Split(ands[i], " OR ")
				for j := range ors {
					bl = bl || testVal(ors[j], args)
				}
			}
		} else {
			if i == 0 {
				bl = testVal(ands[i], args)
			} else {
				bl = bl && testVal(ands[i], args)
			}
		}
	}

	return bl
}

func parseFilter(str string) string {
	str = strings.Replace(str, "&lt;", "<", -1)
	str = strings.Replace(str, "&gt;", ">", -1)
	return str
}

func testVal(testStr string, args map[string]interface{}) bool {
	key, op, testValue := getTestSegment(testStr)
	value, ok := args[key]
	if !ok {
		return false
	}

	return compare(value, op, testValue)
}

func getTestSegment(testStr string) (string, string, string) {
	segments := ptnOp.FindStringSubmatch(testStr)
	if len(segments) != 4 {
		return "", "", ""
	}

	return segments[1], segments[2], segments[3]
}

type (
	Mapper struct {
		XMLName xml.Name `xml:"mapper"`
		Version string   `xml:"version,attr"`

		Selects []SqlItem `xml:"select"`
		Updates []SqlItem `xml:"update"`
		Deletes []SqlItem `xml:"delete"`
		Inserts []SqlItem `xml:"insert"`
	}

	SqlItem struct {
		XMLName xml.Name  `xml:""`
		Sql     string    `xml:",chardata"`
		Id      string    `xml:"id,attr"`
		Where   ItemWhere `xml:"where"`
		Set     ItemSet   `xml:"set"`
		Suffix  string    `xml:"suffix"`
	}

	ItemWhere struct {
		XMLName xml.Name `xml:""`
		Ifs     []ItemIf `xml:"if"`
	}

	ItemSet struct {
		XMLName xml.Name `xml:""`
		Ifs     []ItemIf `xml:"if"`
	}

	ItemIf struct {
		XMLName xml.Name `xml:""`
		If      string   `xml:",chardata"`
		Test    string   `xml:"test,attr"`
	}
)
