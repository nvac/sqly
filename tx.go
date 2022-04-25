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

func TxQueryRow[T any](tx *Tx, scriptName string, args map[string]any) (*T, error) {
    script, err := getScriptByName(scriptName)
    if err != nil {
        return nil, err
    }

    namedScript, namedArgs, err := sqlx.Named(script.content, args)
    if err != nil {
        return nil, err
    }

    row := tx.Tx.QueryRowx(namedScript, namedArgs...)

    var t T
    err = row.StructScan(&t)
    if err != nil {
        return nil, err
    }

    return &t, nil
}

func TxQueryRows[T any](tx *Tx, scriptName string, dest interface{}, args map[string]any) ([]T, error) {
    t := reflect.TypeOf(dest)
    if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Slice || t.Elem().Elem().Kind() != reflect.Struct {
        return nil, errors.New("dest must be ptr slice struct")
    }

    script, err := getScriptByName(scriptName)
    if err != nil {
        return nil, err
    }

    rows, err := tx.Tx.NamedQuery(script.content, args)
    if err != nil {
        return nil, err
    }

    var ts []T
    err = rows.StructScan(ts)
    if err != nil {
        return nil, err
    }

    return ts, nil
}

func (tx *Tx) Exec(scriptName string, args map[string]any) (sql.Result, error) {
    script, err := getScriptByName(scriptName)
    if err != nil {
        return nil, err
    }

    return tx.Tx.Exec(script.content, args)
}
