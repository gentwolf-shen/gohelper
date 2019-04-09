package gomybatis

import (
	"database/sql"
	"errors"
	"gohelper/logger"
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	dbConn       *sql.DB
	mappers      map[string]map[string]SqlItem
	ptnParam     = regexp.MustCompile(`#\{(.*?)\}`)
	formatSql    = "\n%s\n    %s\n -> %s\n => %v"
	ptnCamelCase = regexp.MustCompile(`_([a-z])`)
)

func SetDb(db *sql.DB) {
	dbConn = db
}

func Init(xmlMapperPath string) {
	logger.InitDefault()

	mappers = make(map[string]map[string]SqlItem)

	if !strings.HasSuffix(xmlMapperPath, "/") {
		xmlMapperPath += "/"
	}

	files, err := ioutil.ReadDir(xmlMapperPath)
	if err != nil {
		logger.Error("read mapper path error: " + xmlMapperPath)
		panic(err)
	}

	for _, file := range files {
		filename := strings.ToLower(file.Name())
		if strings.HasSuffix(filename, ".xml") {
			mappers[strings.Split(filename, ".xml")[0]] = parseXML(xmlMapperPath + filename)
		}
	}
}

func Query(selector string, args map[string]interface{}) ([]map[string]string, error) {
	filename, id := parseSelector(selector)
	sqlItem := mappers[filename][id]

	rawSql := buildSelect(&sqlItem, args)
	tsql, values := parseSql(rawSql, args)
	logger.Debugf(formatSql, selector, rawSql, tsql, values)

	rows, err := dbConn.Query(tsql, values...)
	if err != nil {
		return nil, err
	}

	return fetchRows(rows), nil
}

func QueryRow(selector string, args map[string]interface{}) (map[string]string, error) {
	rows, err := Query(selector, args)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return rows[0], nil
}

func QueryScalar(selector string, args map[string]interface{}) (string, error) {
	row, err := QueryRow(selector, args)
	if err != nil {
		return "", err
	}

	value := ""
	for _, val := range row {
		value = val
		break
	}

	return value, nil
}

func Update(selector string, args map[string]interface{}) (int64, error) {
	filename, id := parseSelector(selector)
	sqlItem, ok := mappers[filename][id]
	if !ok {
		return -1, errors.New(selector + " not exists")
	}

	rawSql := buildUpdate(&sqlItem, args)
	tsql, values := parseSql(rawSql, args)

	logger.Debugf(formatSql, selector, rawSql, tsql, values)

	result, err := dbConn.Exec(tsql, values...)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func Delete(selector string, args map[string]interface{}) (int64, error) {
	filename, id := parseSelector(selector)
	sqlItem, ok := mappers[filename][id]
	if !ok {
		return -1, errors.New(selector + " not exists")
	}

	rawSql := buildDelete(&sqlItem, args)
	tsql, values := parseSql(rawSql, args)

	logger.Debugf(formatSql, selector, rawSql, tsql, values)

	result, err := dbConn.Exec(tsql, values...)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func Insert(selector string, args map[string]interface{}) (int64, error) {
	filename, id := parseSelector(selector)
	sqlItem, ok := mappers[filename][id]
	if !ok {
		return -1, errors.New(selector + " not exists")
	}

	rawSql := buildInsert(&sqlItem, args)
	tsql, values := parseSql(rawSql, args)

	logger.Debugf(formatSql, selector, rawSql, tsql, values)

	result, err := dbConn.Exec(tsql, values...)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}

func fetchRows(rows *sql.Rows) []map[string]string {
	if rows == nil {
		return nil
	}

	fields, _ := rows.Columns()
	for k, v := range fields {
		fields[k] = toCamelCase(v)
	}
	columnsLength := len(fields)

	values := make([]string, columnsLength)
	args := make([]interface{}, columnsLength)
	for i := 0; i < columnsLength; i++ {
		args[i] = &values[i]
	}

	index := 0
	listLength := 100
	lists := make([]map[string]string, listLength, listLength)
	for rows.Next() {
		if e := rows.Scan(args...); e == nil {
			row := make(map[string]string, columnsLength)
			for i, field := range fields {
				row[field] = string(values[i])
			}

			if index < listLength {
				lists[index] = row
			} else {
				lists = append(lists, row)
			}
			index++
		}
	}

	rows.Close()

	return lists[0:index]
}

func parseSelector(selector string) (string, string) {
	tmp := strings.Split(selector, ".")
	if len(tmp) != 2 {
		panic(errors.New("gomybatis selector error: " + selector))
	}

	return tmp[0], tmp[1]
}

func parseSql(tsql string, args map[string]interface{}) (string, []interface{}) {
	rs := ptnParam.FindAllStringSubmatch(tsql, -1)
	values := make([]interface{}, len(rs))
	for i := range rs {
		tsql = strings.Replace(tsql, rs[i][0], "?", -1)
		values[i] = args[rs[i][1]]
	}

	return tsql, values
}

func toCamelCase(str string) string {
	return ptnCamelCase.ReplaceAllStringFunc(str, func(a string) string {
		return strings.Title(a[1:2])
	})
}
