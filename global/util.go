package global

import (
	"brandonplank.org/neptune/database"
	"brandonplank.org/neptune/models"
	"bytes"
	"errors"
	"fmt"
	"github.com/Cryptolens/cryptolens-golang/cryptolens"
	csv "github.com/gocarina/gocsv"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	emailClient "github.com/jordan-wright/email"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/smtp"
	"reflect"
	"sort"
	"time"
)

func GetUserFromToken(ctx *fiber.Ctx) (*database.User, error) {
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
	var user database.User

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

func GetNormalUsers() []database.User {
	var users []database.User
	database.DB.Where("permission_level < ?", 2).Find(&users)
	return users
}

func GetAdmins() []database.User {
	var admins []database.User
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
	var teachers []database.User
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
	if len(EmailPassword) == 0 {
		log.Println("No email password set. Skipping.")
		return
	}
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
		// Skip teacher if they have no students.
		if !DoesUserHaveStudents(user.Id.String()) {
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

func GenerateJoinCode() string {
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, 6)
	for i := range s {
		rand.NewSource(time.Now().Unix())
		s[i] = letters[rand.Intn(len(letters))]
	}
	log.Printf("%s-%s", string(s)[:3], string(s)[3:])
	return string(s)
}

func VerifyLicense(license string) bool {
	// Feature 1 (F1) is unlimited time in our case
	log.Printf("Verifying key: %s", license)

	token := "WyIxODM1MzIzNSIsIkNrV0toSEkyTVc2N2VjWktrTFhSZ0ZsMHNnMUhydWRFbHlGQXV2bU4iXQ=="
	publicKey := "<RSAKeyValue><Modulus>39PoIXYeHQijPn8v/4oMzJNQQyF8NdZyFDho7CZz0cczrGefc7/LB5/RSHA1akzwiG/xiNDE5XT+zaDr9qslkklKEWspYsC0oOAgvhBHKepYmOy5kf9JlIu8dTFAoJDGj4Mhb3hGV/DB8gvCbMnE+H7Q/jXPw289k0q9DXbqKicUot6EDkcB46PxMqMqFe0ykSzxc6vk0cyhbH+dX/ncL/eDkjwRQjPzl7SGkHK7pU+ndaDaRLmrb9BzcGLwepmQaQnJuJmnn0wSDRTAlJ3ynQOGiB5U+0B0zXpBz17qmcz6gVbY/soUMpO0OyDAynS98YU+spnshZRJG4yTlkg22Q==</Modulus><Exponent>AQAB</Exponent></RSAKeyValue>"

	licenseKey, err := cryptolens.KeyActivate(token, cryptolens.KeyActivateArguments{
		ProductId: 15153,
		Key:       license,
	})

	if err != nil || !licenseKey.HasValidSignature(publicKey) {
		log.Println("License key activation failed!")
		return false
	}

	if time.Now().After(licenseKey.Expires) && licenseKey.F1 {
		log.Println("Neptune license key has expired")
		return false
	}

	log.Println("Neptune license verified")
	return true
}

//func GetLocalIP() string {
//	addrs, err := net.InterfaceAddrs()
//	if err != nil {
//		return ""
//	}
//	for _, address := range addrs {
//		// check the address type and if it is not a loopback the display it
//		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
//			if ipnet.IP.To4() != nil {
//				return ipnet.IP.String()
//			}
//		}
//	}
//	return ""
//}

func getLocalIPAddr(version string) (ip string) {

	ifaces, err := net.Interfaces()

	if err == nil {
		for _, i := range ifaces {
			addrs, err := i.Addrs()

			if err == nil {
				for _, addr := range addrs {
					// in case you have an interface ip4 generated by another program for example docker
					if len(ip) > 0 {
						break
					}

					ipnet, status := addr.(*net.IPNet)

					if status == true {
						switch version {
						case "v4":
							// prevent 127 and check IP4.
							if !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
								ip = ipnet.IP.String()
							}
						case "v6":
							// prevent ::1 and check IP6.
							if !ipnet.IP.IsLoopback() && ipnet.IP.To4() == nil {
								ip = ipnet.IP.String()
							}
						}
					}
				}
			}
		}
	}
	return

}

func GetIPv4() (ip string) {
	return getLocalIPAddr("v4")
}

func GetIPv6() (ip string) {
	return getLocalIPAddr("v6")
}
