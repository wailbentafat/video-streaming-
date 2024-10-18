package main

import (
	"fmt"
	"net/http"
	"videoback/control" 
	"videoback/model"
)

func main() { 
	database := db.InitDB()
	defer func() {
		sqlDB, _ := database.DB()
		sqlDB.Close()
	}()

	fmt.Println("Server is starting...")

	http.HandleFunc("/upload", Upload.Upload)

	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
