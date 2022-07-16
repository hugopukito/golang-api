package db

import (
	"database/sql"
	"fmt"

	"module.com/webServer/user"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB
var err error

func init() {
	DB, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connection db success...")

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS bgs (id INT PRIMARY KEY NOT NULL AUTO_INCREMENT, nom VARCHAR(50))")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Create table if not exists success...")
}

func GetAll() []string {
	results, err := DB.Query("SELECT nom FROM bgs")
	if err != nil {
		panic(err.Error())
	}

	var users []string

	for results.Next() {
		var user user.User

		err = results.Scan(&user.Name)
		if err != nil {
			panic(err.Error())
		}

		users = append(users, user.Name)
	}

	return users
}
