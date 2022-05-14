module brandonplank.org/signx/routes

go 1.17

require (
	brandonplank.org/neptune/database v0.0.0
	brandonplank.org/neptune/global v0.0.0
	brandonplank.org/neptune/embed v0.0.0
	github.com/gocarina/gocsv v0.0.0-20220310154401-d4df709ca055
	github.com/gofiber/fiber/v2 v2.31.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.3.0
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292
)

require (
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible // indirect
	gorm.io/driver/mysql v1.3.3 // indirect
	gorm.io/gorm v1.23.4 // indirect
)

require (
	brandonplank.org/neptune/models v0.0.0
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/klauspost/compress v1.15.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.34.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
)

replace (
	brandonplank.org/neptune/embed => ../embed
	brandonplank.org/neptune/database => ../database
	brandonplank.org/neptune/global => ../global
	brandonplank.org/neptune/models => ../models
)
