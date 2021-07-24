package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// ErrNotFound 没有找到数据
var ErrNotFound = errors.New("not found")

type QueryErr struct {
	RawSQL string
	Err    error
}

func (e *QueryErr) Unwrap() error {
	return e.Err
}

func (e *QueryErr) As(in interface{}) bool {
	if v, ok := in.(*(*QueryErr)); ok {
		*v = e
		return true
	}
	return errors.As(e.Err, in)
}

func (e *QueryErr) Is(target error) bool {
	if target == sql.ErrNoRows && errors.Is(e.Err, ErrNotFound) {
		// 恢复被特殊标记的错误
		return true
	}
	return errors.Is(e.Err, target)
}

func (e *QueryErr) Error() string {
	return e.Err.Error()
}

func GetPerson(id int) (record Person, err error) {
	query := fmt.Sprintf("select id, name from foo  WHERE id=%d", id)
	err = db.QueryRow(query).Scan(&record.ID, &record.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 特殊值ErrNoRows，特殊标记
			err = fmt.Errorf("%w,no such id:%d , ", ErrNotFound, id)
		}
		err = &QueryErr{RawSQL: query, Err: fmt.Errorf("Query Sql:%s ;error:%w", query, err)}
	}
	return
}

func main() {
	defer db.Close()

	_, err := GetPerson(0)
	if errors.Is(err, sql.ErrTxDone) { // 会调用 Unwrap
		panic(err)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	if !errors.Is(err, ErrNotFound) {
		panic(err)
	}

	var qerr *QueryErr
	if !errors.As(err, &qerr) {
		panic(err)
	} else {
		println(qerr.RawSQL, qerr.Err.Error())
	}
}

var (
	//
	db *sql.DB = initDB()
)

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		log.Fatal(err)
	}
	sqlStmt := `
	create table foo (id integer not null primary key, name text);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

type Person struct {
	ID   int
	Name string
}
