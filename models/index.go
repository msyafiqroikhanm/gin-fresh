package models

import "gorm.io/gorm"

var DB *gorm.DB // Variabel DB untuk menyimpan koneksi database

// Inisialisasi DB
func InitDB(db *gorm.DB) {
	DB = db
}
