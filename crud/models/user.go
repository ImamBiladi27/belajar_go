package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Struct User dengan field password
type User struct {
    gorm.Model
    Name     string
    Email    string
    Phone    string
    Password string
}

func main() {
    dsn := "root:@tcp(127.0.0.1:3306)/gocrud?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect to DB:", err)
    }

    err = db.AutoMigrate(&User{})
    if err != nil {
        log.Fatal("failed to migrate:", err)
    }

    log.Println("Migrasi berhasil!")
}
