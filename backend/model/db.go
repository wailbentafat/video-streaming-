package db

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "log"
	
)


func InitDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open("mydatabase.db"), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect to the database: %v", err)
    }
	err = db.AutoMigrate(&Video{}, &User{})
    if err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }

    return db
}
