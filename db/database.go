package db

import (
	"database/sql"
	"fmt"

	"module.com/webServer/entity"

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

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS messages (id INT PRIMARY KEY NOT NULL AUTO_INCREMENT, name VARCHAR(30), message VARCHAR(255))")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Create table 'messages' if not exists success...")

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY NOT NULL AUTO_INCREMENT, name VARCHAR(255), email VARCHAR(255), password VARCHAR(255))")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Create table 'users' if not exists success...")
}

func FindUser(email string) entity.User {
	result, err := DB.Query("select name, email, password from users where email = ? order by email asc limit 1;", email)
	if err != nil {
		panic(err.Error())
	}

	var user entity.User

	for result.Next() {
		err = result.Scan(&user.Name, &user.Email, &user.Password)
		if err != nil {
			panic(err.Error())
		}
	}

	return user
}

func InsertUser(user entity.User) {
	insert := "INSERT INTO users (name, email, password) values (?, ?, ?)"
	_, err := DB.Exec(insert, user.Name, user.Email, user.Password)
	if err != nil {
		panic(err.Error())
	}
}
