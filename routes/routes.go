package routes

import (
	"brandonplank.org/neptune/database"
	"brandonplank.org/neptune/global"
	"brandonplank.org/neptune/models"
	"encoding/base64"
	"fmt"
	csv "github.com/gocarina/gocsv"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	guuid "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"sort"
	"strings"
	"time"
)

func Register(ctx *fiber.Ctx) error {
	type RegisterPayload struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		SchoolID string `json:"schoolID"`
	}

	var payload RegisterPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, "must set proper data")
	}

	if len(payload.Name) < 1 || len(payload.Email) < 1 || len(payload.Password) < 1 {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, "must set name, email, password")
	}

	var schoolID guuid.UUID

	if len(payload.SchoolID) > 0 {
		var err error
		schoolID, err = guuid.Parse(payload.SchoolID)
		if err != nil {
			return global.CraftReturnStatus(ctx, fiber.StatusBadRequest, "Malformed UUID")
		}
	} else {
		schoolID = guuid.Nil
	}

	var readUser models.User
	database.DB.Where("email = ?", payload.Email).First(&readUser)

	if readUser.Name == payload.Name {
		return global.CraftReturnStatus(ctx, fiber.StatusForbidden, "User already exists")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	user := models.User{
		Name:     payload.Name,
		SchoolId: schoolID,
		Email:    payload.Email,
		Password: string(password),
	}

	database.DB.Create(&user)

	return ctx.JSON(user)
}

func AdminRegister(ctx *fiber.Ctx) error {
	userSignedIn, err := global.GetUserFromToken(ctx)
	if err != nil {
		return err
	}

	if userSignedIn.PermissionLevel < 1 {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, "You must have permission level 1 or higher")
	}

	type RegisterPayload struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		SchoolID        string `json:"schoolID"`
		PermissionLevel uint   `json:"permissionLevel"`
	}

	var payload RegisterPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	if len(payload.Name) < 1 || len(payload.Email) < 1 || len(payload.Password) < 1 {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, "must set name, email, password")
	}

	var schoolID guuid.UUID

	if len(payload.SchoolID) > 0 {
		var err error
		schoolID, err = guuid.Parse(payload.SchoolID)
		if err != nil {
			return global.CraftReturnStatus(ctx, fiber.StatusBadRequest, "Malformed UUID")
		}
	} else {
		if userSignedIn.PermissionLevel < 5 {
			schoolID = userSignedIn.SchoolId
		} else {
			schoolID = guuid.Nil
		}
	}

	// Check to make sure use cannot elevate another user higher than oneself
	if userSignedIn.PermissionLevel <= payload.PermissionLevel {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, fmt.Sprintf("You can only set permissions lover than yourself, you are permission level %d", userSignedIn.PermissionLevel))
	}

	var readUser models.User
	database.DB.Where("email = ?", payload.Email).First(&readUser)

	if readUser.Name == payload.Name {
		return global.CraftReturnStatus(ctx, fiber.StatusForbidden, "User already exists")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	user := models.User{
		Name:            payload.Name,
		SchoolId:        schoolID,
		Email:           payload.Email,
		Password:        string(password),
		PermissionLevel: payload.PermissionLevel,
	}

	database.DB.Create(&user)

	return ctx.JSON(user)
}

func Login(ctx *fiber.Ctx) error {
	type LoginPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var payload LoginPayload

	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	if len(payload.Email) < 1 || len(payload.Password) < 1 {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, "must set proper data")
	}

	var user models.User
	database.DB.Where("email = ?", strings.ToLower(payload.Email)).First(&user)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		log.Println(err)
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, "Passwords do not match")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Id.String(),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 12).Unix(), // half a day
	})

	token, err := claims.SignedString([]byte(global.SignedJWTKey))
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusInternalServerError, err.Error())
	}

	cookie := fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 12),
		HTTPOnly: true,
	}
	ctx.Cookie(&cookie)
	log.Println(fmt.Sprintf("%s is logging in.", user.Name))
	return global.CraftReturnStatus(ctx, fiber.StatusOK, "Logged in")
}

func Logout(ctx *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	ctx.Cookie(&cookie)
	return global.CraftReturnStatus(ctx, fiber.StatusOK, "success")
}

func User(ctx *fiber.Ctx) error {
	user, err := global.GetUserFromToken(ctx)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}
	return ctx.JSON(user)
}

func AddSchool(ctx *fiber.Ctx) error {
	user, err := global.GetUserFromToken(ctx)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	if user.PermissionLevel < 9 {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, "You must have permission level 9 or higher")
	}

	var data map[string]interface{}
	if err := ctx.BodyParser(&data); err != nil {
		return err
	}

	if data["name"] == nil {
		return global.CraftReturnStatus(ctx, fiber.StatusBadRequest, "Must add a school name")
	}

	school := models.School{
		Name: data["name"].(string),
	}

	database.DB.Create(&school)

	return global.CraftReturnStatus(ctx, fiber.StatusOK, "Created school")
}

func GetSchool(ctx *fiber.Ctx) error {
	user, err := global.GetUserFromToken(ctx)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}
	var school models.School
	database.DB.Where("id = ?", user.SchoolId).First(&school)
	if school.Id == guuid.Nil {
		return global.CraftReturnStatus(ctx, fiber.StatusBadRequest, "You are not apart of a school")
	}
	return ctx.JSON(school)
}

func GetSchools(ctx *fiber.Ctx) error {
	var schools []models.School
	database.DB.Find(&schools)
	return ctx.JSON(schools)
}

func IsOut(ctx *fiber.Ctx) error {
	nameBase64 := ctx.Params("name")
	nameData, err := base64.URLEncoding.DecodeString(nameBase64)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	studentName := string(nameData)

	// Get the teacher
	user, err := global.GetUserFromToken(ctx)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	students := global.GetStudentsFromUserID(user.Id.String())

	type out struct {
		IsOut bool   `json:"isOut"`
		Name  string `json:"name"`
	}
	return ctx.JSON(out{IsOut: global.IsStudentOut(studentName, students), Name: studentName})
}

func Id(ctx *fiber.Ctx) error {
	nameBase64 := ctx.Params("name")
	nameData, err := base64.URLEncoding.DecodeString(nameBase64)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	studentName := string(nameData)

	// Get the teacher
	user, err := global.GetUserFromToken(ctx)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	students := global.GetStudentsFromUserID(user.Id.String())

	if global.IsStudentOut(studentName, students) {
		log.Println(fmt.Sprintf("%s has retured to %s's classroom", studentName, user.Name))
		var student models.Student
		database.DB.Where("name = ?", studentName).Where("sign_in = ?", "Signed Out").First(&student).Update("sign_in", time.Now().Format("3:04 pm"))
	} else {
		log.Println(fmt.Sprintf("%s has left from %s's classroom", studentName, user.Name))
		student := models.Student{
			Name:      studentName,
			SignOut:   time.Now().Format("3:04 pm"),
			SignIn:    "Signed Out",
			Date:      time.Now().Format("01/02/2006"),
			TeacherId: user.Id,
		}
		database.DB.Create(&student)
	}
	return global.CraftReturnStatus(ctx, fiber.StatusOK, "success")
}

func GetCSV(ctx *fiber.Ctx) error {
	user, err := global.GetUserFromToken(ctx)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	students := global.GetStudentsFromUserID(user.Id.String())

	publicStudents := models.StudentsToPublicStudents(students)
	sort.Sort(publicStudents)
	global.ReverseSlice(publicStudents)
	studentsBytes, err := csv.MarshalBytes(publicStudents)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusInternalServerError, "Could not marshal student list")
	}
	return ctx.Send(studentsBytes)
}

func GetAdminCSV(ctx *fiber.Ctx) error {
	user, err := global.GetUserFromToken(ctx)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	students := global.GetSchoolSignoutsFromSchoolID(user.SchoolId.String())

	publicStudents := models.StudentsToPublicStudents(students)
	sort.Sort(publicStudents)
	global.ReverseSlice(publicStudents)
	studentsBytes, err := csv.MarshalBytes(publicStudents)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusInternalServerError, "Could not marshal student list")
	}
	return ctx.Send(studentsBytes)
}

func CSVFile(ctx *fiber.Ctx) error {
	user, err := global.GetUserFromToken(ctx)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	students := global.GetStudentsFromUserID(user.Id.String())

	publicStudents := models.StudentsToPublicStudents(students)
	sort.Sort(publicStudents)
	studentsBytes, err := csv.MarshalBytes(publicStudents)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusInternalServerError, "Could not marshal student list")
	}
	ctx.Append("Content-Disposition", "attachment; filename=\"classroom.csv\"")
	ctx.Append("Content-Type", "text/csv")
	return ctx.Send(studentsBytes)
}

func GetAdminCSVFile(ctx *fiber.Ctx) error {
	user, err := global.GetUserFromToken(ctx)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	allStudents := global.GetSchoolSignoutsFromSchoolID(user.SchoolId.String())

	publicStudents := models.StudentsToPublicStudents(allStudents)
	sort.Sort(publicStudents)
	studentsBytes, err := csv.MarshalBytes(publicStudents)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusInternalServerError, "Could not marshal student list")
	}
	ctx.Append("Content-Disposition", "attachment; filename=\"admin.csv\"")
	ctx.Append("Content-Type", "text/csv")
	return ctx.Send(studentsBytes)
}

func GetUserPermissionLevel(ctx *fiber.Ctx) error {
	user, err := global.GetUserFromToken(ctx)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	type Level struct {
		Level uint `json:"level"`
	}
	return ctx.JSON(Level{Level: user.PermissionLevel})
}

func AdminSearchStudent(ctx *fiber.Ctx) error {
	nameBase64 := ctx.Params("name")
	nameData, err := base64.URLEncoding.DecodeString(nameBase64)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	studentName := string(nameData)

	user, err := global.GetUserFromToken(ctx)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	allStudents := global.GetSchoolSignoutsFromSchoolID(user.SchoolId.String())

	var retStudents models.Students

	for _, student := range allStudents {
		if strings.Contains(strings.ToLower(student.Name), strings.ToLower(studentName)) {
			retStudents = append(retStudents, student)
		}
	}

	publicStudents := models.StudentsToPublicStudents(retStudents)
	sort.Sort(publicStudents)
	global.ReverseSlice(publicStudents)

	content, _ := csv.MarshalBytes(publicStudents)
	return ctx.Send(content)
}

func PasswordChange(ctx *fiber.Ctx) error {
	user, err := global.GetUserFromToken(ctx)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	type PasswordChnagePayload struct {
		Password    string `json:"password"`
		NewPassword string `json:"newPassword"`
	}

	var payload PasswordChnagePayload

	if err := ctx.BodyParser(&payload); err != nil {
		return err
	}

	if len(payload.Password) < 1 || len(payload.NewPassword) < 1 {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, "must set proper data")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		log.Println(err)
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, "Passwords do not match")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	database.DB.Model(&user).Update("password", string(password))

	return global.CraftReturnStatus(ctx, fiber.StatusOK, "changed password")
}

func Admin(ctx *fiber.Ctx) error {
	user, err := global.GetUserFromToken(ctx)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	if user.PermissionLevel < 1 {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, "you lack admin level 1 or above")
	}

	return ctx.Render("superadmin", fiber.Map{
		"year":                  time.Now().Format("2006"),
		"logo":                  "assets/img/viking_logo.png",
		"canAddSchool":          user.PermissionLevel >= 9,
		"canAddDifferentSchool": user.PermissionLevel >= 5,
	})
}

func PasswordChangeSite(ctx *fiber.Ctx) error {
	_, err := global.GetUserFromToken(ctx)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	return ctx.Render("passwordreset", fiber.Map{
		"year": time.Now().Format("2006"),
		"logo": "assets/img/viking_logo.png",
	})
}

func Home(ctx *fiber.Ctx) error {
	logoURL := "assets/img/viking_logo.png"

	cookie := ctx.Cookies("token")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.SignedJWTKey), nil
	})
	if err != nil || !token.Valid {
		return ctx.Render("login", fiber.Map{
			"year": time.Now().Format("2006"),
			"logo": logoURL,
		})
	}

	user, err := global.GetUserFromToken(ctx)
	if err != nil {
		return global.CraftReturnStatus(ctx, fiber.StatusUnauthorized, err.Error())
	}

	if user.PermissionLevel >= 3 {
		return Admin(ctx)
	}

	if user.PermissionLevel > 1 {
		return ctx.Render("admin", fiber.Map{
			"year":             time.Now().Format("2006"),
			"logo":             logoURL,
			"hasAdminFeatures": true,
		})
	}

	return ctx.Render("main", fiber.Map{
		"year":             time.Now().Format("2006"),
		"logo":             logoURL,
		"hasAdminFeatures": user.PermissionLevel > 0,
	})
}
