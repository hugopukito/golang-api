#!/bin/bash

git pull
go build main.go
sudo systemctl restart go-back-service.service
