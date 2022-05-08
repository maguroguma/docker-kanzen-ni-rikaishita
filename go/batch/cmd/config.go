package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var host string
var port int
var db *gorm.DB

func init() {
	setDeps()

	user := "hoge"
	password := "passw0rd"
	dbname := "testdb"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Tokyo", host, user, password, dbname, port)
	tempDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db = tempDB
	fmt.Println("[Success] initialize")
}

func setDeps() {
	envName := os.Getenv("ENV_NAME")
	fmt.Println("[DEBUG]:", envName)
	if envName == "docker-containers" {
		host = "postgres_db"
		port = 5432
	} else {
		host = "localhost"
		port = 15432
	}
}
