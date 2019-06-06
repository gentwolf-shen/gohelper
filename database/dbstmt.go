package database

import (
	"database/sql"
	"strings"
)

type DbStmt struct {
	query string
	stmt  *sql.Stmt
}

func NewDbStmt(stmt *sql.Stmt, query string) *DbStmt {
	return &DbStmt{query, stmt}
}

func (this *DbStmt) Query(args ...interface{}) ([]map[string]string, error) {
	rows, err := this.stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	return fetchRows(rows), nil
}

func (this *DbStmt) QueryRow(args ...interface{}) (map[string]string, error) {
	rows, err := this.Query(args...)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return rows[0], nil
}

func (this *DbStmt) QueryScalar(args ...interface{}) (string, error) {
	row, err := this.QueryRow(args...)
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

func (this *DbStmt) Insert(args ...interface{}) (int64, error) {
	var id int64
	var err error

	if strings.Contains(this.query, " RETURNING ") {
		if rows, err1 := this.stmt.Query(args...); err1 != nil {
			return 0, err1
		} else {
			rows.Next()
			err = rows.Scan(&id)
		}
	} else {
		if result, err1 := this.stmt.Exec(args...); err1 != nil {
			return 0, err1
		} else {
			id, err = result.LastInsertId()
		}
	}

	return id, err
}

func (this *DbStmt) Update(args ...interface{}) (int64, error) {
	var n int64

	result, err := this.stmt.Exec(args...)
	if err == nil {
		n, err = result.RowsAffected()
	}

	return n, err
}

func (this *DbStmt) Delete(args ...interface{}) (int64, error) {
	return this.Update(args...)
}
