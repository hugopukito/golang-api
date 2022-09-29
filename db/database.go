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

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS bgs (id INT PRIMARY KEY NOT NULL AUTO_INCREMENT, nom VARCHAR(30), message VARCHAR(255))")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Create table 'bgs' if not exists success...")

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS users (id INT PRIMARY KEY NOT NULL AUTO_INCREMENT, name VARCHAR(255), email VARCHAR(255), password VARCHAR(255), role VARCHAR(255))")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Create table 'users' if not exists success...")
}

func GetAllMessages() []entity.Message {
	results, err := DB.Query("SELECT nom, message FROM bgs")
	if err != nil {
		panic(err.Error())
	}

	var messages []entity.Message

	for results.Next() {
		var message entity.Message

		err = results.Scan(&message.Name, &message.Message)
		if err != nil {
			panic(err.Error())
		}
		messages = append(messages, message)
	}

	return messages
}

func InsertMessage(message entity.Message) {
	insert := "INSERT INTO bgs (nom, message) values (?, ?)"
	_, err := DB.Exec(insert, message.Name, message.Message)
	if err != nil {
		panic(err.Error())
	}
}

func FindUser(user entity.User) {

}

func InsertUser(user entity.User) {

}
