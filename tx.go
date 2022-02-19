package sqly

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"reflect"
)

type Tx struct {
	*sqlx.Tx
	DatabaseName string
}

func Begin(databaseName string) (*Tx, error) {
	if !globalInit {
		return nil, errors.New("sqly has not been initialized, please use 'sqly.Init(&sqly.Config{})' to initialize")
	}

	database, err := getDatabaseByName(databaseName)
	if err != nil {
		return nil, err
	}
	tx, err := database.db.Beginx()

	return &Tx{Tx: tx, DatabaseName: databaseName}, err
}

func (tx *Tx) Rollback() error {
	return tx.Tx.Rollback()
}

func (tx *Tx) Commit() error {
	return tx.Tx.Rollback()
}

func (tx *Tx) QueryRow(scriptName string, dest interface{}, arg interface{}) error {
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("dest must be ptr struct")
	}

	script, err := getScriptByName(scriptName)
	if err != nil {
		return err
	}

	rows, err := tx.Tx.NamedQuery(script.content, arg)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(dest)
		if err != nil {
			return err
		}
		break
	}

	return nil
}

func (tx *Tx) QueryRows(scriptName string, dest interface{}, arg interface{}) error {
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Slice || t.Elem().Elem().Kind() != reflect.Struct {
		return errors.New("dest must be ptr slice struct")
	}

	script, err := getScriptByName(scriptName)
	if err != nil {
		return err
	}

	rows, err := tx.Tx.NamedQuery(script.content, arg)
	if err != nil {
		return err
	}

	defer rows.Close()

	err = sqlx.StructScan(rows, dest)
	return err
}

func (tx *Tx) Exec(scriptName string, arg interface{}) (sql.Result, error) {
	script, err := getScriptByName(scriptName)
	if err != nil {
		return nil, err
	}

	return tx.Tx.NamedExec(script.content, arg)
}
