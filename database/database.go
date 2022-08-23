package database

import (
	"brandonplank.org/neptune/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User models.User

var DB *gorm.DB

func Connect() {
	connection, err := gorm.Open(mysql.Open("root:********************@tcp(database.brandonplank.org)/neptune?parseTime=true"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = connection
	connection.AutoMigrate(&models.School{})
	connection.AutoMigrate(&models.User{})
	connection.AutoMigrate(&models.Student{})
}

const (
	Teacher = iota
	TeacherWhoCanAddTeacher
	SchoolAdmin
	SchoolIT
	DistrictAdmin
	_
	_
	_
	_
	_
	_
	SuperAdmin
)

func (u *User) SetPermission(flag uint) {
	DB.Model(&u).Update("permission_level", u.PermissionLevel|flag)
}
