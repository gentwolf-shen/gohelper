package database

import (
	"database/sql"
	"strings"

	"github.com/gentwolf-shen/gohelper/hashhelper"
)

type Database struct {
	dbType string
	dbConn *sql.DB
	stmts  map[string]*sql.Stmt
}

func New() *Database {
	return &Database{}
}

func (this *Database) Open(dbType, dsn string, maxOpenConnections, maxIdleConnections int) error {
	this.dbType = dbType
	var err error

	this.dbConn, err = sql.Open(this.dbType, dsn)
	if err == nil {
		this.dbConn.SetMaxOpenConns(maxOpenConnections)
		this.dbConn.SetMaxIdleConns(maxIdleConnections)

		this.stmts = make(map[string]*sql.Stmt, 10)
	}

	return err
}

func (this *Database) Close() error {
	this.CloseAllStmt()
	return this.dbConn.Close()
}

func (this *Database) GetConn() *sql.DB {
	return this.dbConn
}

func (this *Database) Query(query string, args ...interface{}) ([]map[string]string, error) {
	rows, err := this.dbConn.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return fetchRows(rows), nil
}

func (this *Database) QueryRow(query string, args ...interface{}) (map[string]string, error) {
	rows, err := this.Query(query, args...)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return rows[0], nil
}

func (this *Database) QueryScalar(query string, args ...interface{}) (string, error) {
	row, err := this.QueryRow(query, args...)
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

func (this *Database) Insert(query string, args ...interface{}) (int64, error) {
	var id int64
	var err error

	if strings.Contains(query, " RETURNING ") {
		row := this.dbConn.QueryRow(query, args...)
		err = row.Scan(&id)
	} else {
		if result, err1 := this.dbConn.Exec(query, args...); err1 != nil {
			return 0, err1
		} else {
			id, err = result.LastInsertId()
		}
	}
	return id, err
}

func (this *Database) Update(query string, args ...interface{}) (int64, error) {
	var n int64

	result, err := this.dbConn.Exec(query, args...)
	if err == nil {
		n, err = result.RowsAffected()
	}

	return n, err
}

func (this *Database) Delete(query string, args ...interface{}) (int64, error) {
	return this.Update(query, args...)
}

func (this *Database) CreateStmt(query string, name ...string) (*DbStmt, error) {
	key := ""
	if len(name) > 0 {
		key = name[0]
	} else {
		key = hashhelper.Md5(query)
	}

	stmt, ok := this.stmts[key]
	if !ok || stmt == nil {
		var err error
		stmt, err = this.dbConn.Prepare(query)
		if err != nil {
			return nil, err
		}

		this.stmts[key] = stmt
	}

	return NewDbStmt(stmt, query), nil
}

func (this *Database) CloseStmt(name string) {
	if stmt, ok := this.stmts[name]; ok {
		_ = stmt.Close()
		delete(this.stmts, name)
	}
}

func (this *Database) CloseAllStmt() {
	for name, stmt := range this.stmts {
		_ = stmt.Close()
		delete(this.stmts, name)
	}
}
