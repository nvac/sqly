package sqly

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"reflect"
)

func Init(config *Config) error {
	if config.DatabasesFile == "" {
		return errors.New("DatabasesFile is required")
	}

	if config.ScriptsGlobFiles == "" {
		return errors.New("ScriptsGlobFiles is required")
	}

	globalConfig = config

	err := loadDatabasesFile()
	if err != nil {
		return err
	}

	err = watchDatabasesFile()
	if err != nil {
		return err
	}

	err = loadScriptsGlobFiles()
	if err != nil {
		return err
	}

	return nil
}

func QueryRow(name string, dest interface{}, arg interface{}) error {
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("dest must be ptr struct")
	}

	database, script, err := getScriptByName(name)
	if err != nil {
		return err
	}

	rows, err := database.db.NamedQuery(script.content, arg)
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

func QueryRows(name string, dest interface{}, arg interface{}) error {
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Slice || t.Elem().Elem().Kind() != reflect.Struct {
		return errors.New("dest must be ptr slice struct")
	}

	database, script, err := getScriptByName(name)
	if err != nil {
		return err
	}

	rows, err := database.db.NamedQuery(script.content, arg)
	if err != nil {
		return err
	}

	defer rows.Close()

	err = sqlx.StructScan(rows, dest)
	return err
}

func Exec(name string, arg interface{}) (sql.Result, error) {
	database, script, err := getScriptByName(name)
	if err != nil {
		return nil, err
	}

	return database.db.NamedExec(script.content, arg)
}
