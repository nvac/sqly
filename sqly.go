package sqly

import (
    "context"
    "database/sql"
    "errors"
)

type DB struct {
    *sql.DB
}

func Init(config Config) error {
    if config.DatabasesFile == "" {
        return errors.New("DatabasesFile is required")
    }

    if config.ScriptsGlobFiles == "" {
        return errors.New("ScriptsGlobFiles is required")
    }

    globalConfig = &config

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

func QueryRow[T any](databaseName, scriptName string, args map[string]any) (*T, error) {
    db, _ := connect(databaseName)

    script, err := getScriptByName(scriptName)
    if err != nil {
        return nil, err
    }

    namedArgs := argsToNamedArgs(args)
    rows, err := db.Query(script.content, namedArgs...)
    if err != nil {
        return nil, err
    }

    result, err := scanRowsToStruct[T](rows)
    if err != nil {
        return nil, err
    }

    return &result[0], nil
}

func QueryRows[T any](databaseName, scriptName string, arg map[string]any) (*[]T, error) {
    db, _ := connect(databaseName)

    script, err := getScriptByName(scriptName)
    if err != nil {
        return nil, err
    }

    rows, err := db.Query(script.content, arg)
    if err != nil {
        return nil, err
    }

    result, err := scanRowsToStruct[T](rows)

    return &result, err
}

func Exec(databaseName, scriptName string, args map[string]any) (sql.Result, error) {
    db, _ := connect(databaseName)

    script, err := getScriptByName(scriptName)
    if err != nil {
        return nil, err
    }

    namedArgs := argsToNamedArgs(args)
    return db.Exec(script.content, namedArgs)
}

func ExecContext(ctx context.Context, databaseName, scriptName string, args map[string]any) (sql.Result, error) {
    db, _ := connect(databaseName)

    script, err := getScriptByName(scriptName)
    if err != nil {
        return nil, err
    }

    namedArgs := argsToNamedArgs(args)
    return db.ExecContext(ctx, script.content, namedArgs)
}

func connect(databaseName string) (*DB, error) {
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
