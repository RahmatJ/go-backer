package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-backer/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	//dsn := "root:@tcp(127.0.0.1:3306)/backer_db?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	//
	//fmt.Println("Connected to DB")
	//
	//var users []user.User
	//
	////gorm automatically detect database and table name
	//db.Find(&users)
	//
	//for _, userData := range users {
	//	fmt.Println(userData.Name)
	//	fmt.Println(userData.Email)
	//	fmt.Println("=============")
	//}

	router := gin.Default()
	router.GET("/handler", handler)

	err := router.Run()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func handler(c *gin.Context) {
	dsn := "root:@tcp(127.0.0.1:3306)/backer_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connected to DB")

	var users []user.User

	//gorm automatically detect database and table name
	db.Find(&users)

	c.JSON(http.StatusOK, users)
}
