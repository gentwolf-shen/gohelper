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

func getValueInfo(valuePtr interface{}) (reflect.Value, reflect.Type) {
	var value reflect.Value

	valueType := reflect.TypeOf(valuePtr).Elem()
	kind := valueType.Kind()

	if kind == reflect.Struct {
		value = reflect.ValueOf(valuePtr).Elem()
	} else if kind == reflect.Slice {
		valueType = valueType.Elem()
		value = reflect.New(valueType).Elem()
	}

	return value, valueType
}

func getValues(value reflect.Value, valueType reflect.Type, columns []string) ([]reflect.Value, []interface{}) {
	columnSize := len(columns)
	for i, v := range columns {
		columns[i] = toCamelCase(v)
	}

	dest := make([]reflect.Value, columnSize)
	elem := make([]interface{}, columnSize)

	empty := ""
	nilValue := reflect.Indirect(reflect.ValueOf(empty))
	fieldSize := valueType.NumField()

	for i := 0; i < columnSize; i++ {
		bl := false
		for j := 0; j < fieldSize; j++ {
			field := valueType.Field(j)
			if field.Tag.Get("db") == columns[i] || toLowFirst(field.Name) == columns[i] {
				dest[i] = reflect.Indirect(value.Field(j))
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

	return dest, elem
}

func fetchObjectRowsForMore(valuePtr interface{}, rows *sql.Rows, isSingleRow bool) error {
	if rows == nil {
		return nil
	}

	columns, _ := rows.Columns()
	for i, v := range columns {
		columns[i] = toCamelCase(v)
	}

	value, valueType := getValueInfo(valuePtr)
	dest, elem := getValues(value, valueType, columns)

	items := reflect.Indirect(reflect.ValueOf(valuePtr))

	for rows.Next() {
		if e := rows.Scan(elem...); e == nil {
			for k, v := range elem {
				if dest[k].CanSet() {
					dest[k].Set(reflect.ValueOf(v).Elem())
				}
			}

			if isSingleRow {
				break
			} else {
				items.Set(reflect.Append(items, value))
			}
		}
	}

	_ = rows.Close()

	return nil
}
