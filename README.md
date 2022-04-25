## sqly

<hr/>

An easy-to-use extension for [sqlx](https://github.com/jmoiron/sqlx) ，base on xml files and named query/exec

<p style="color: orangered">this repo is under development, please do not use it in production.</p>

### install

``go get github.com/nvac/sqly``

### Usage

#### 1. set database config in xml file

* name: needs to be unique in same environment
* environment: custom string，runtime environment
* source: data source name
* connMaxLifetime(seconds): sets the maximum amount of time a connection may be reused. . if default values is required,
  remove the attr
* connMaxIdleTime(seconds): sets the maximum amount of time a connection may be idle. if default values is required,
  remove the attr
* maxIdleConns: sets the maximum number of connections in the idle connection pool. if default values is required,
  remove the attr
* maxOpenConns: sets the maximum number of open connections to the database. if default values is required, remove the
  attr

````xml
<?xml version="1.0" encoding="utf-8" ?>

<databases>
    <database name="ReadDb"
              environment="development"
              driver="mysql"
              source="user:password@tcp(127.0.0.1:3306)/test?charset=utf8mb4&amp;parseTime=True"
              connMaxLifetime="30"
              connMaxIdleTime="30"
              maxIdleConns="2"
              maxOpenConns="10"
    />
    
    <database name="WriteDb"
              environment="development"
              driver="mysql"
              source="user:password@tcp(127.0.0.1:3306)/test?charset=utf8mb4&amp;parseTime=True"
              connMaxLifetime="30"
              connMaxIdleTime="30"
              maxIdleConns="2"
              maxOpenConns="10"
    />
</databases>
````

#### 2. write sql script in xml file

* name: needs to be unique
* database: using the above configured database
* content: ensure in CDATA

````xml
<?xml version="1.0" encoding="utf-8" ?>

<scripts>
    <script name="GetUser">
        <![CDATA[
            SELECT username, password
            FROM `user`
            WHERE username = :username
        ]]>
    </script>

    <script name="ListUser">
        <![CDATA[
            SELECT username, password
            FROM `user`
            LIMIT 10 OFFSET 0
        ]]>
    </script>
    
    <script name="AddUser">
        <![CDATA[
            INSERT INTO user (username, password)
            VALUES (:username, :password)
        ]]>
    </script>
</scripts>
````

3. inti & use sqly

````go
package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nvac/sqly"
	"os"
)

type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

func main() {
	err := sqly.Init(&Config{
		DatabasesFile:    "config/databases.xml",
		ScriptsGlobFiles: "config/scripts/*.xml",
		Environment:      os.Getenv("Environment"),
		SourceDecryptFunc: func(source string) string {
			return source
		},
	})

	if err != nil {
		panic(err)
	}

	readDb, err :=sqly.Connect("ReadDb")
	if err != nil {
		panic(err)
	}
	
	var user User
	if err := readDb.QueryRow("GetUser", &user, map[string]interface{}{
		"username": "lisa",
	}); err != nil {
		panic(err)
	} else {
		fmt.Println(user)
	}

	var users []User
	if err := readDb.QueryRows("ListUser", &users, map[string]interface{}{}); err != nil {
		panic(err)
	} else {
		fmt.Println(users)
	}

	writeDb, err := sqly.Connect("WriteDb")
	if err != nil {
		panic(err)
	}
	if result, err := writeDb.Exec("AddUser", map[string]interface{}{
		"username": "root",
		"password": "123456",
	}); err != nil {
		panic(err)
	} else {
		fmt.Println(result.RowsAffected())
		fmt.Println(result.LastInsertId())
	}

	tx, err := sqly.Begin("WriteDb")

	if result, err := tx.Exec("AddUser", map[string]interface{}{
		"username": "root",
		"password": "123456",
	}); err != nil {
		if err := tx.Rollback(); err != nil {
			panic(err)
		}
	} else {
		fmt.Println(result.RowsAffected())
		fmt.Println(result.LastInsertId())
		tx.Commit()
	}
}
````

### License

[MIT](LICENSE) © nvac