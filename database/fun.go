package database

import (
	"database/sql"
	"reflect"
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
				row[field] = values[i]
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

func toLowFirst(str string) string {
	return strings.ToLower(str[0:1]) + str[1:]
}

func fetchObjectRow(value interface{}, rows *sql.Rows) error {
	_, err := fetchObjectRowsForMore(value, rows, true)
	return err
}

func fetchObjectRows(value interface{}, rows *sql.Rows) ([]interface{}, error) {
	return fetchObjectRowsForMore(value, rows, false)
}

func fetchObjectRowsForMore(value interface{}, rows *sql.Rows, isSingleRow bool) ([]interface{}, error) {
	if rows == nil {
		return nil, nil
	}

	columns, _ := rows.Columns()
	for i, v := range columns {
		columns[i] = toCamelCase(v)
	}
	columnSize := len(columns)
	types := reflect.TypeOf(value).Elem()
	fieldSize := types.NumField()

	values := reflect.ValueOf(value).Elem()
	dest := make([]reflect.Value, columnSize)
	elem := make([]interface{}, columnSize)

	empty := ""
	nilValue := reflect.Indirect(reflect.ValueOf(empty))

	for i := 0; i < columnSize; i++ {
		bl := false
		for j := 0; j < fieldSize; j++ {
			field := types.Field(j)
			if field.Tag.Get("db") == columns[i] || toLowFirst(field.Name) == columns[i] {
				dest[i] = reflect.Indirect(values.Field(j))
				elem[i] = reflect.New(dest[i].Type()).Interface()
				bl = true
				break
			}
		}

		if !bl {
			dest[i] = nilValue
			elem[i] = reflect.New(dest[i].Type()).Interface()
		}
	}

	maxLength := 1
	if !isSingleRow {
		maxLength = 100
	}
	items := make([]interface{}, maxLength)

	index := 0
	for rows.Next() {
		if e := rows.Scan(elem...); e == nil {
			for k, v := range elem {
				if dest[k].CanSet() {
					dest[k].Set(reflect.ValueOf(v).Elem())
				}
			}

			if index < maxLength {
				items[index] = reflect.ValueOf(value).Elem().Interface()
			} else {
				items = append(items, reflect.ValueOf(value).Elem().Interface())
			}
			index++

			if isSingleRow {
				break
			}
		}
	}

	_ = rows.Close()

	return items[0:index], nil
}
