package db

import (
	"database/sql"
	"log"

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

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS messages (id INT PRIMARY KEY NOT NULL AUTO_INCREMENT, name VARCHAR(30), message VARCHAR(255))")
	if err != nil {
		panic(err.Error())
	}

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY NOT NULL AUTO_INCREMENT, name VARCHAR(255), email VARCHAR(255), password VARCHAR(255))")
	if err != nil {
		panic(err.Error())
	}
}

func FindUser(email string) entity.User {
	result, err := DB.Query("select name, email, password from users where email = ? order by email asc limit 1;", email)
	if err != nil {
		log.Println("error FindUser")
	}

	var user entity.User

	for result.Next() {
		err = result.Scan(&user.Name, &user.Email, &user.Password)
		if err != nil {
			log.Println("error FindUser")
		}
	}

	return user
}

func InsertUser(user entity.User) {
	insert := "INSERT INTO users (name, email, password) values (?, ?, ?)"
	_, err := DB.Exec(insert, user.Name, user.Email, user.Password)
	if err != nil {
		log.Println("error InsertUser")
	}
}
