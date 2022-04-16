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

	args := map[string]interface{}{
		"username": "root",
	}
	row, err := QueryRow[User]("ReadDb", "GetUser", args)

	if err != nil {
		fmt.Printf("db.NamedQuery failed, err:%v\n", err)
		return
	}

	fmt.Println(row)
}

func TestQueryRows(t *testing.T) {
	err := Init(&Config{
		DatabasesFile:    "config/databases.xml",
		ScriptsGlobFiles: "config/scripts/*.xml",
	})

	if err != nil {
		panic(err)
	}

	users, err := QueryRows[User]("ReadDb", "ListUser", nil)

	if err != nil {
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

	if result, err := Exec("WriteDb", "AddUser", map[string]interface{}{
		"username": "root",
		"password": "123456",
	}); err != nil {

	} else {
		fmt.Println(result.RowsAffected())
		fmt.Println(result.LastInsertId())
	}
}
