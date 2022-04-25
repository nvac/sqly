package sqly

import (
    "database/sql"
    "errors"
    "github.com/jmoiron/sqlx"
)

type DB struct {
    *sqlx.DB
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

    namedScript, namedArgs, err := sqlx.Named(script.content, args)
    if err != nil {
        return nil, err
    }

    row := db.QueryRowx(namedScript, namedArgs...)

    var t T
    err = row.StructScan(&t)
    if err != nil {
        return nil, err
    }

    return &t, nil
}

func QueryRows[T any](databaseName, scriptName string, args map[string]any) ([]T, error) {
    db, _ := connect(databaseName)

    script, err := getScriptByName(scriptName)
    if err != nil {
        return nil, err
    }

    rows, err := db.NamedQuery(script.content, args)
    if err != nil {
        return nil, err
    }

    var ts []T
    err = rows.StructScan(ts)
    if err != nil {
        return nil, err
    }

    return ts, err
}

func Exec(databaseName, scriptName string, args map[string]any) (sql.Result, error) {
    db, _ := connect(databaseName)

    script, err := getScriptByName(scriptName)
    if err != nil {
        return nil, err
    }

    return db.NamedExec(script.content, args)
}
