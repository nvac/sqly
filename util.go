package sqly

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"go/ast"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"time"
)

func getScriptByName(name string) (*scriptsCacheValue, error) {
	scriptsValue, ok := scriptsCache.Load(name)
	if !ok {
		return nil, errors.New("not found script")
	}

	script := scriptsValue.(scriptsCacheValue)

	return &script, nil
}

func getDatabaseByName(databaseName string) (*databasesCacheValue, error) {
	databasesCacheKey := getDatabasesCacheKey(databaseName)
	databasesValue, ok := databasesCache.Load(databasesCacheKey)
	if !ok {
		return nil, errors.New("not found database")
	}

	database := databasesValue.(databasesCacheValue)
	if database.ping {
		err := database.db.Ping()
		if err != nil {
			return nil, err
		}

		database.ping = true
	}

	return &database, nil
}

func loadDatabasesFile() error {
	if globalConfig.DatabasesFile == "" {
		return errors.New("miss DatabasesFile")
	}
	file, _ := ioutil.ReadFile(globalConfig.DatabasesFile)

	data := new(databases)
	if err := xml.Unmarshal(file, data); err != nil {
		return err
	}

	if data == nil || data.Databases == nil {
		return errors.New("no available database")
	}

	var cacheDatabase []string
	for _, database := range data.Databases {
		if globalConfig.Environment != "" && database.Environment != globalConfig.Environment {
			continue
		}

		source := database.Source
		if globalConfig.SourceDecryptFunc != nil {
			source = globalConfig.SourceDecryptFunc(database.Source)
		}

		db, err := sql.Open(database.Driver, source)
		if err != nil {
			return err
		}

		if database.MaxOpenConns != nil {
			db.SetMaxOpenConns(*database.MaxOpenConns)
		}

		if database.MaxIdleConns != nil {
			db.SetMaxIdleConns(*database.MaxIdleConns)
		}

		if database.ConnMaxLifetime != nil {
			seconds := time.Duration(*database.ConnMaxLifetime)
			db.SetConnMaxLifetime(seconds * time.Second)
		}

		if database.ConnMaxIdleTime != nil {
			seconds := time.Duration(*database.ConnMaxIdleTime)
			db.SetConnMaxIdleTime(seconds * time.Second)
		}

		databasesCacheKey := getDatabasesCacheKey(database.Name)
		value := databasesCacheValue{
			db:   db,
			name: database.Name,
			ping: false,
		}
		cacheDatabase = append(cacheDatabase, databasesCacheKey)
		databasesCache.Store(databasesCacheKey, value)
	}

	lintDatabasesCache(cacheDatabase)

	return nil
}

func loadScriptsGlobFile(path string) error {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	data := new(scripts)
	err = xml.Unmarshal(fileContent, data)
	if err != nil {
		return err
	}

	for _, script := range data.Scripts {
		key := script.Name
		value := scriptsCacheValue{
			name:    script.Name,
			content: script.Content,
			path:    path,
		}

		if loadedValue, loaded := scriptsCache.LoadOrStore(key, value); loaded {
			loadedValue := loadedValue.(scriptsCacheValue)
			message := fmt.Sprintf("The duplicate script name('%s') in the current file('%s') was found in another file('%s')",
				loadedValue.name, path, loadedValue.path)
			return errors.New(message)
		}
	}
	return nil
}

func loadScriptsGlobFiles() error {
	files, _ := filepath.Glob(globalConfig.ScriptsGlobFiles)
	for _, file := range files {
		if err := loadScriptsGlobFile(file); err != nil {
			return err
		}
	}
	return nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func argsToNamedArgs(args map[string]any) []any {
	var result []any
	for name, value := range args {
		result = append(result, sql.Named(name, value))
	}
	return result
}

func scanRowsToStruct[T any](rows *sql.Rows) ([]T, error) {
	var destSlice []T
	destType := reflect.ValueOf(destSlice).Type().Elem()
	fieldNames := parseTag(destType)
	for rows.Next() {
		var dest T
		var values []interface{}
		elem := reflect.ValueOf(&dest).Elem()
		for _, fieldName := range fieldNames {
			elem.FieldByName(fieldName).Addr().Interface()
		}

		if err := rows.Scan(values...); err != nil {
			return nil, err
		}
		destSlice = append(destSlice, dest)
	}
	return destSlice, nil
}

func parseTag(dest interface{}) []string {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	var result []string
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			v, ok := p.Tag.Lookup("sql")
			if !ok || v == "" {
				result = append(result, p.Name)
				continue
			}

			if v == "-" {
				continue
			}

			result = append(result, v)
		}
	}
	return result
}
