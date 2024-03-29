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

## Create dir with secret for jwt tokens

cd && mkdir secrets && cd secrets && touch secret_jwt.txt

Then add your long password/secret in it

## Install redis key_value database to store chat messages

sudo apt install redis

## Create dir of imgs

mkdir imgs-back

put img name in this dir as img id you pass in path variables in request

jpg only

## Run api
go run main.go

## Systemd

[Unit]
Description=Go back-end

[Service]
User=pukito
WorkingDirectory=/home/pukito/back-go
ExecStart=/snap/bin/go run main.go
Restart=always
Environment=GOMODCACHE=/home/pukito/go/pkg/mod
Environment=GOPATH=/home/pukito/go
Environment=GOCACHE=/home/pukito/go/pkg/cache
StandardOutput=file:/var/log/go-back-service.log
StandardError=file:/var/log/go-back-service.log

[Install]
WantedBy=multi-user.target
