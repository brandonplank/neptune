module brandonplank.org/neptune

go 1.17

require (
	brandonplank.org/neptune/database v0.0.0
	brandonplank.org/neptune/routes v0.0.0
	github.com/aldy505/sentry-fiber v0.0.1
	github.com/gofiber/fiber/v2 v2.31.0
	github.com/gofiber/template v1.6.26
)

replace (
	brandonplank.org/neptune/controllers => ./controllers
	brandonplank.org/neptune/database => ./database
	brandonplank.org/neptune/global => ./global
	brandonplank.org/neptune/models => ./models
	brandonplank.org/neptune/routes => ./routes
)

require (
	brandonplank.org/neptune/global v0.0.0 // indirect
	brandonplank.org/neptune/models v0.0.0 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/getsentry/sentry-go v0.11.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gocarina/gocsv v0.0.0-20220310154401-d4df709ca055 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.15.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.35.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	gorm.io/driver/mysql v1.3.3 // indirect
	gorm.io/gorm v1.23.4 // indirect
)
