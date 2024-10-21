package db

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "log"
	
)

var DB *gorm.DB 

func InitDB() *gorm.DB {
    var err error
    DB, err = gorm.Open(sqlite.Open("mydatabase.db"), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect to the database: %v", err)
    }
    err = DB.AutoMigrate(&Video{}, &User{})
    if err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }

    return DB // You can still return DB if needed
}
