package global

import (
	"brandonplank.org/neptune/database"
	"brandonplank.org/neptune/models"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"reflect"
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

func IsRailway() bool {
	vars := []string{"RAILWAY_STATIC_URL", "RAILWAY_ENVIRONMENT", "RAILWAY_HEALTHCHECK_TIMEOUT_SEC", "RAILWAY_GIT_COMMIT_SHA", "RAILWAY_GIT_AUTHOR", "RAILWAY_GIT_BRANCH", "RAILWAY_GIT_REPO_NAME", "RAILWAY_GIT_REPO_OWNER", "RAILWAY_GIT_COMMIT_MESSAGE"}
	for _, env := range vars {
		_, hasVar := os.LookupEnv(env)
		if hasVar {
			return true
		}
	}
	return false
}
