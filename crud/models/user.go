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

// Struct Penjualan dengan foreign key UserID yang mereference ke User.ID
type Penjualan struct {
    gorm.Model
    UserID uint   `gorm:"null"` // Foreign key ke User
    User   User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Relasi ke User
    Produk string
    Jumlah int
    Harga  float64
}

func main() {
    dsn := "root:@tcp(127.0.0.1:3306)/gocrud?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect to DB:", err)
    }

    // Cek apakah tabel users sudah ada
    if db.Migrator().HasTable(&User{}) {
        // Jika tabel User sudah ada, migrasi hanya tabel Penjualan
        err = db.AutoMigrate(&Penjualan{})
        if err != nil {
            log.Fatal("failed to migrate Penjualan:", err)
        }
        log.Println("Tabel Penjualan berhasil dimigrasi!")
    } else {
        // Jika tabel User belum ada, migrasi kedua tabel
        err = db.AutoMigrate(&User{}, &Penjualan{})
        if err != nil {
            log.Fatal("failed to migrate:", err)
        }
        log.Println("Migrasi User dan Penjualan berhasil!")
    }
}
