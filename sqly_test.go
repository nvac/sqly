package sqly

import (
    "fmt"
    "os"
    "testing"

    _ "github.com/mattn/go-sqlite3"
)

type User struct {
    Username string `db:"username"`
    Password string `db:"password"`
}

func TestQueryRow(t *testing.T) {
    err := Init(&Config{
        DatabasesFile:    "config/databases.xml",
        ScriptsGlobFiles: "config/scripts/*.xml",
        Environment:      os.Getenv("Environment"),
        SourceDecryptFunc: func(source string) string {
            return source
        },
    })

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    db, err := Connect("ReadDb")
    if err != nil {
        panic(err)
    }

    var users User
    err = db.QueryRow("GetUser", &users, map[string]interface{}{
        "username": "root",
    })

    if err != nil {
        fmt.Printf("db.NamedQuery failed, err:%v\n", err)
        return
    }

    fmt.Println(users)
}

func TestQueryRows(t *testing.T) {
    err := Init(&Config{
        DatabasesFile:    "config/databases.xml",
        ScriptsGlobFiles: "config/scripts/*.xml",
    })

    if err != nil {
        panic(err)
    }

    db, err := Connect("ReadDb")
    if err != nil {
        panic(err)
    }
    var users []User
    if err := db.QueryRows("ListUser", &users, map[string]interface{}{}); err != nil {
        fmt.Printf("db.NamedQuery failed, err:%v\n", err)
        return
    }

    fmt.Println(users)
}

func TestExec(t *testing.T) {
    err := Init(&Config{
        DatabasesFile:    "config/databases.xml",
        ScriptsGlobFiles: "config/scripts/*.xml",
        SourceDecryptFunc: func(source string) string {
            return source
        },
    })

    if err != nil {
        panic(err)
    }

    db, err := Connect("WriteDb")

    if result, err := db.Exec("AddUser", map[string]interface{}{
        "username": "root",
        "password": "123456",
    }); err != nil {

    } else {
        fmt.Println(result.RowsAffected())
        fmt.Println(result.LastInsertId())
    }
}
