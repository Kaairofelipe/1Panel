package main

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func main() {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	stmt := &gorm.Statement{DB: db}
	fmt.Println("Table quote:", stmt.Quote("users; drop table users;"))
}
