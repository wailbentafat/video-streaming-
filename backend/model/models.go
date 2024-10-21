package db

import ()
type User struct{
	ID int `gorm:"primaryKey"`
	Password string `gorm:"not null"`
	Email string `gorm:"not null"`
	Videos []Video `gorm:"foreignkey:User_id"`
}

type Video struct {
	ID         int    `gorm:"primaryKey"`
	VideoName  string `gorm:"not null"`
	VideoPath  string `gorm:"not null"`
	HLSPath    string `gorm:"not null"`  // Path to the .m3u8 playlist file
	User_id    int    `gorm:"not null"`
	User       User   `gorm:"foreignKey:User_id"`
}