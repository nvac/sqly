package sqly

import (
	"database/sql"
	"errors"
	"reflect"
)

type Tx struct {
	*sql.Tx
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
	tx, err := database.db.Begin()

	return &Tx{Tx: tx, DatabaseName: databaseName}, err
}

func (tx *Tx) Rollback() error {
	return tx.Tx.Rollback()
}

func (tx *Tx) Commit() error {
	return tx.Tx.Rollback()
}

func (tx *Tx) QueryRow(scriptName string, dest interface{}, args map[string]any) error {
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("dest must be ptr struct")
	}

	script, err := getScriptByName(scriptName)
	if err != nil {
		return err
	}

	namedArgs := argsToNamedArgs(args)
	rows, err := tx.Tx.Query(script.content, namedArgs)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(dest)
		if err != nil {
			return err
		}
		break
	}

	return nil
}

func (tx *Tx) QueryRows(scriptName string, dest interface{}, args map[string]any) error {
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Slice || t.Elem().Elem().Kind() != reflect.Struct {
		return errors.New("dest must be ptr slice struct")
	}

	script, err := getScriptByName(scriptName)
	if err != nil {
		return err
	}

	namedArgs := argsToNamedArgs(args)
	rows, err := tx.Tx.Query(script.content, namedArgs)
	if err != nil {
		return err
	}

	defer rows.Close()
	err = rows.Scan(dest)
	return err
}

func (tx *Tx) Exec(scriptName string, args map[string]any) (sql.Result, error) {
	script, err := getScriptByName(scriptName)
	if err != nil {
		return nil, err
	}

	return tx.Tx.Exec(script.content, args)
}
