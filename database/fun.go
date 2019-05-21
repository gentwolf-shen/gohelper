package database

import (
	"database/sql"
	"regexp"
	"strings"
)

var (
	ptnCamelCase = regexp.MustCompile(`_([a-z])`)
)

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

	_ = rows.Close()

	return lists[0:index]
}

func toCamelCase(str string) string {
	return ptnCamelCase.ReplaceAllStringFunc(str, func(a string) string {
		return strings.Title(a[1:2])
	})
}
