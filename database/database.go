package database

import (
	"brandonplank.org/neptune/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	connection, err := gorm.Open(mysql.Open("root:HsgGSFXLh0ppvLTJwser@tcp(containers-us-west-30.railway.app:5942)/railway?parseTime=true"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = connection
	connection.AutoMigrate(&models.School{})
	connection.AutoMigrate(&models.User{})
	connection.AutoMigrate(&models.Student{})
}
