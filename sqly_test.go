package sqly

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"testing"
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

	var users User
	err = QueryRow("GetUser", &users, map[string]interface{}{
		"username": "doovac",
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
		fmt.Println(err.Error())
	}

	var users []User
	if err := QueryRows("ListUser", &users, map[string]interface{}{}); err != nil {
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
		fmt.Println(err.Error())
	}

	if result, err := Exec("AddUser", map[string]interface{}{
		"username": "root",
		"password": "123456",
	}); err != nil {

	} else {
		fmt.Println(result.RowsAffected())
		fmt.Println(result.LastInsertId())
	}
}
