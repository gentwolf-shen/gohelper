package gomybatis

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"

	"github.com/gentwolf-shen/gohelper/hashhelper"
	"github.com/gentwolf-shen/gohelper/logger"
)

var (
	dbConns      map[string]*sql.DB
	mappers      map[string]map[string]SqlItem
	ptnParam     = regexp.MustCompile(`#\{(.*?)\}`)
	ptnParamVar  = regexp.MustCompile(`\$\{(.*?)\}`)
	ptnCamelCase = regexp.MustCompile(`_([a-z0-9])`)
	formatSql    = "\n%s\n    %s\n -> %s\n => %v"
	stmts        map[string]*sql.Stmt
)

func initMapper() {
	logger.InitDefault()

	if mappers == nil {
		mappers = make(map[string]map[string]SqlItem)
	}

	if dbConns == nil {
		dbConns = make(map[string]*sql.DB)
	}

	if stmts == nil {
		stmts = make(map[string]*sql.Stmt)
	}
}

func SetMapper(dbConn *sql.DB, name, xml string) {
	initMapper()

	mappers[name] = parseXmlFromStr(xml)
	dbConns[name] = dbConn
}

func SetMapperPath(dbConn *sql.DB, mapperPath string) {
	initMapper()

	if !strings.HasSuffix(mapperPath, "/") {
		mapperPath += "/"
	}

	files, err := ioutil.ReadDir(mapperPath)
	if err != nil {
		logger.Error("read mapper path error: " + mapperPath)
		panic(err)
	}

	for _, file := range files {
		filename := strings.ToLower(file.Name())
		if strings.HasSuffix(filename, ".xml") {
			basename := strings.Split(filename, ".xml")[0]
			mappers[basename] = parseXmlFromFile(mapperPath + filename)
			dbConns[basename] = dbConn
		}
	}
}

func Query(selector string, args map[string]interface{}) ([]map[string]string, error) {
	rows, err := query(selector, args)
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
	return UpdateTrans(nil, selector, args)
}

func Delete(selector string, args map[string]interface{}) (int64, error) {
	return DeleteTrans(nil, selector, args)
}

func Insert(selector string, args map[string]interface{}) (int64, error) {
	return InsertTrans(nil, selector, args)
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
		return "", ""
	}

	return tmp[0], tmp[1]
}

func parseSql(tsql string, args map[string]interface{}) (string, []interface{}) {
	tsql = ptnParamVar.ReplaceAllStringFunc(tsql, func(a string) string {
		return args[a[2:len(a)-1]].(string)
	})

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

func toLowFirst(str string) string {
	return strings.ToLower(str[0:1]) + str[1:]
}

func selectorNotExists(selector string) error {
	return errors.New("selector \"" + selector + "\" not exists!")
}

func Transaction(filenameForDb string, fun func(tx *sql.Tx) error) error {
	tx, err := dbConns[filenameForDb].Begin()

	if err == nil {
		if err = fun(tx); err == nil {
			err = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}

	return err
}

func UpdateTrans(tx *sql.Tx, selector string, args map[string]interface{}) (int64, error) {
	filename, id := parseSelector(selector)
	sqlItem, ok := mappers[filename][id]
	if !ok {
		return -1, selectorNotExists(selector)
	}

	rawSql := buildUpdate(&sqlItem, args)
	tsql, values := parseSql(rawSql, args)

	logger.Debugf(formatSql, selector, rawSql, tsql, values)

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.Exec(tsql, values...)
	} else {
		result, err = dbConns[filename].Exec(tsql, values...)
	}

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func DeleteTrans(tx *sql.Tx, selector string, args map[string]interface{}) (int64, error) {
	filename, id := parseSelector(selector)
	sqlItem, ok := mappers[filename][id]
	if !ok {
		return -1, selectorNotExists(selector)
	}

	rawSql := buildDelete(&sqlItem, args)
	tsql, values := parseSql(rawSql, args)

	logger.Debugf(formatSql, selector, rawSql, tsql, values)

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.Exec(tsql, values...)
	} else {
		result, err = dbConns[filename].Exec(tsql, values...)
	}

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func InsertTrans(tx *sql.Tx, selector string, args map[string]interface{}) (int64, error) {
	filename, id := parseSelector(selector)
	sqlItem, ok := mappers[filename][id]
	if !ok {
		return -1, selectorNotExists(selector)
	}

	rawSql := buildInsert(&sqlItem, args)
	tsql, values := parseSql(rawSql, args)

	logger.Debugf(formatSql, selector, rawSql, tsql, values)

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.Exec(tsql, values...)
	} else {
		result, err = dbConns[filename].Exec(tsql, values...)
	}

	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}

func getStmt(filename, tsql string) (*sql.Stmt, error) {
	var err error

	key := hashhelper.Md5(filename + ":" + tsql)
	stmt, ok := stmts[key]
	if !ok {
		stmt, err = dbConns[filename].Prepare(tsql)
		stmts[key] = stmt
	}

	return stmt, err
}

func Close() {
	for name, stmt := range stmts {
		_ = stmt.Close()
		delete(stmts, name)
	}
}

func query(selector string, args map[string]interface{}) (*sql.Rows, error) {
	filename, id := parseSelector(selector)
	sqlItem, ok := mappers[filename][id]
	if !ok {
		return nil, selectorNotExists(selector)
	}

	rawSql := buildSelect(&sqlItem, args)
	tsql, values := parseSql(rawSql, args)
	logger.Debugf(formatSql, selector, rawSql, tsql, values)

	return dbConns[filename].Query(tsql, values...)
}

func QueryObject(value interface{}, selector string, args map[string]interface{}) error {
	rows, err := query(selector, args)
	if err != nil {
		return nil
	}

	return fetchObjectRow(value, rows)
}

func QueryObjects(value interface{}, selector string, args map[string]interface{}) error {
	rows, err := query(selector, args)
	if err != nil {
		return nil
	}

	return fetchObjectRows(value, rows)
}

func fetchObjectRow(value interface{}, rows *sql.Rows) error {
	return fetchObjectRowsForMore(value, rows, true)
}

func fetchObjectRows(value interface{}, rows *sql.Rows) error {
	return fetchObjectRowsForMore(value, rows, false)
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
