package sqly

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

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

	err = loadScriptsGlobFiles()
	if err != nil {
		return err
	}

	globalInit = true
	return nil
}

func Connect(databaseName string) (*DB, error) {
	if !globalInit {
		return nil, errors.New("sqly has not been initialized, please use 'sqly.Init(&sqly.Config{})' to initialize")
	}

	databaseCache, err := getDatabaseByName(databaseName)
	if err != nil {
		return nil, err
	}

	if databaseCache.ping {
		err := databaseCache.db.Ping()
		if err != nil {
			return nil, err
		}
		databaseCache.ping = true
	}

	return &DB{DB: databaseCache.db}, nil
}

func (db *DB) QueryRow(scriptName string, dest interface{}, arg interface{}) error {
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("dest must be ptr struct")
	}

	script, err := getScriptByName(scriptName)
	if err != nil {
		return err
	}

	rows, err := db.NamedQuery(script.content, arg)
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

func (db *DB) QueryRows(scriptName string, dest interface{}, arg interface{}) error {
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Slice || t.Elem().Elem().Kind() != reflect.Struct {
		return errors.New("dest must be ptr slice struct")
	}

	script, err := getScriptByName(scriptName)
	if err != nil {
		return err
	}

	rows, err := db.NamedQuery(script.content, arg)
	if err != nil {
		return err
	}

	defer rows.Close()

	err = sqlx.StructScan(rows, dest)
	return err
}

func (db *DB) Exec(scriptName string, arg interface{}) (sql.Result, error) {
	script, err := getScriptByName(scriptName)
	if err != nil {
		return nil, err
	}

	return db.NamedExec(script.content, arg)
}
