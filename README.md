# Golang api connection to a mysql database running on the host

## Install mysql on your host

sudo apt update

sudo apt install -y mysql-server

sudo mysql

ALTER USER 'root'@'localhost' IDENTIFIED WITH caching_sha2_password BY 'password';

FLUSH PRIVILEGES;

## Create database

mysql -u root -p

CREATE DATABASE golang;

## Install redis key_value database to store chat messages

sudo apt install redis

## Run api
go run main.go