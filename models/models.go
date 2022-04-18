package models

import (
	guuid "github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type School struct {
	Id   guuid.UUID `gorm:"primary_key" json:"id"`
	Name string     `json:"name"`
}

func (base *School) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := guuid.New()
	base.Id = uuid
	return
}

type Students []Student
type PublicStudents []PublicStudent

type Student struct {
	Id        guuid.UUID `csv:"-" gorm:"primary_key" json:"id"`
	TeacherId guuid.UUID `csv:"-" json:"TeacherId"`
	Name      string     `csv:"Name" json:"name" bson:"name"`
	SignOut   string     `csv:"Signed Out" json:"signedOut" bson:"signedOut"`
	SignIn    string     `csv:"Signed In" json:"signedIn" bson:"signedIn"`
	Date      string     `csv:"Date" json:"date" bson:"date"`
}

func (base *Student) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := guuid.New()
	base.Id = uuid
	return
}

func (p Students) Len() int {
	return len(p)
}

func (p Students) Less(i, j int) bool {
	time1, _ := time.Parse("3:04 pm", p[i].SignOut)
	time2, _ := time.Parse("3:04 pm", p[j].SignOut)
	return time1.Before(time2)
}

func (p Students) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type PublicStudent struct {
	Name    string `csv:"Name" json:"name" bson:"name"`
	SignOut string `csv:"Signed Out" json:"signedOut" bson:"signedOut"`
	SignIn  string `csv:"Signed In" json:"signedIn" bson:"signedIn"`
	Date    string `csv:"Date" json:"date" bson:"date"`
}

func (p PublicStudents) Len() int {
	return len(p)
}

func (p PublicStudents) Less(i, j int) bool {
	time1, _ := time.Parse("3:04 pm", p[i].SignOut)
	time2, _ := time.Parse("3:04 pm", p[j].SignOut)
	return time1.Before(time2)
}

func (p PublicStudents) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func StudentsToPublicStudents(students Students) PublicStudents {
	var publicStudents PublicStudents
	for _, student := range students {
		publicStudents = append(publicStudents, PublicStudent{Name: student.Name, SignOut: student.SignOut, SignIn: student.SignIn, Date: student.Date})
	}
	return publicStudents
}

type User struct {
	Id              guuid.UUID `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	SchoolId        guuid.UUID `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email" gorm:"unique"`
	Password        string     `json:"-"`
	PermissionLevel uint       `json:"permissionLevel"`
}

func (base *User) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := guuid.New()
	base.Id = uuid
	return
}
