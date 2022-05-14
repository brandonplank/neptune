package global

import (
	"brandonplank.org/neptune/database"
	"brandonplank.org/neptune/models"
	"bytes"
	"errors"
	"fmt"
	csv "github.com/gocarina/gocsv"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	emailClient "github.com/jordan-wright/email"
	"log"
	"net/http"
	"net/smtp"
	"reflect"
	"sort"
)

func GetUserFromToken(ctx *fiber.Ctx) (*models.User, error) {
	cookie := ctx.Cookies("token")
	if len(cookie) < 5 {
		return nil, errors.New("please set jwt token")
	}
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SignedJWTKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid jwt token")
	}

	claims := token.Claims.(*jwt.StandardClaims)
	var user models.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)
	return &user, nil
}

func ReverseSlice(data interface{}) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		panic(errors.New("data must be a slice type"))
	}
	valueLen := value.Len()
	if valueLen < 1 {
		return
	}
	for i := 0; i <= (valueLen-1)/2; i++ {
		reverseIndex := valueLen - 1 - i
		tmp := value.Index(reverseIndex).Interface()
		value.Index(reverseIndex).Set(value.Index(i))
		value.Index(i).Set(reflect.ValueOf(tmp))
	}
}

func IsStudentOut(name string, students []models.Student) bool {
	if students == nil {
		return false
	}
	for _, stu := range students {
		if stu.Name == name {
			if stu.SignIn == "Signed Out" {
				return true
			}
		}
	}
	return false
}

func CraftReturnStatus(ctx *fiber.Ctx, status int, message string) error {
	return ctx.Status(status).JSON(fiber.Map{
		"statusString": http.StatusText(status),
		"status":       status,
		"message":      message,
	})
}

func CleanStudents() {
	var students models.Students
	database.DB.Find(&students).Delete(&students)
}

func GetStudentsFromUserID(id string) models.Students {
	var students models.Students
	database.DB.Where("teacher_id = ?", id).Find(&students)
	return students
}

func GetSchoolFromUser(id string) models.School {
	var school models.School
	database.DB.Where("id = ?", id).First(&school)
	return school
}

func GetNormalUsers() []models.User {
	var users []models.User
	database.DB.Where("permission_level < ?", 2).Find(&users)
	return users
}

func GetAdmins() []models.User {
	var admins []models.User
	database.DB.Where("permission_level > ? AND permission_level <= ?", 1, 5).Find(&admins)
	return admins
}

func DoesUserHaveStudents(id string) bool {
	students := GetStudentsFromUserID(id)
	if len(students) > 0 {
		return true
	}
	return false
}

func GetSchoolSignoutsFromSchoolID(id string) models.Students {
	var students models.Students
	var teachers []models.User
	database.DB.Where("school_id = ?", id).Find(&teachers)
	for _, teacher := range teachers {
		students = append(students, GetStudentsFromUserID(teacher.Id.String())...)
	}
	return students
}

func DoesSchoolHaveSignouts(id string) bool {
	students := GetSchoolSignoutsFromSchoolID(id)
	if len(students) > 0 {
		return true
	}
	return false
}

func EmailUsers() {
	for _, admin := range GetAdmins() {
		if !DoesSchoolHaveSignouts(admin.SchoolId.String()) {
			continue
		}
		allStudents := GetSchoolSignoutsFromSchoolID(admin.SchoolId.String())
		sort.Sort(allStudents)
		ReverseSlice(allStudents)
		content, _ := csv.MarshalBytes(allStudents)
		csvReader := bytes.NewReader(content)

		school := GetSchoolFromUser(admin.Id.String())

		schoolEmail := emailClient.NewEmail()
		schoolEmail.From = "Neptune <planksprojects@gmail.com>"
		schoolEmail.Subject = "Neptune Sign-Outs"
		schoolEmail.To = []string{admin.Email}
		schoolEmail.Text = []byte("This is an automated email to " + school.Name)
		schoolEmail.Attach(csvReader, fmt.Sprintf("%s.csv", school.Name), "text/csv")
		err := schoolEmail.Send("smtp.gmail.com:587", smtp.PlainAuth("", "planksprojects@gmail.com", EmailPassword, "smtp.gmail.com"))
		if err != nil {
			log.Println(err)
		}
	}

	for _, user := range GetNormalUsers() {
		if !DoesUserHaveStudents(user.Id.String()) {
			log.Println(fmt.Sprintf("%s has no students, not sending email", user.Name))
			continue
		}
		students := GetStudentsFromUserID(user.Id.String())

		csvClass, err := csv.MarshalBytes(students)
		if err != nil {
			log.Println(err)
		}
		if len(csvClass) < 1 {
			continue
		}
		csvReader := bytes.NewReader(csvClass)
		classroomEmail := emailClient.NewEmail()
		classroomEmail.From = "Neptune <planksprojects@gmail.com>"
		classroomEmail.Subject = "Neptune Sign-Outs"
		classroomEmail.To = []string{user.Email}
		classroomEmail.Text = []byte("This is an automated email to " + user.Name)
		classroomEmail.Attach(csvReader, fmt.Sprintf("%s.csv", user.Name), "text/csv")
		err = classroomEmail.Send("smtp.gmail.com:587", smtp.PlainAuth("", "planksprojects@gmail.com", EmailPassword, "smtp.gmail.com"))
		if err != nil {
			log.Println(err)
		}
	}
}
