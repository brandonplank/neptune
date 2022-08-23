module brandonplank.org/neptune

go 1.19

require (
	brandonplank.org/neptune/database v0.0.0
	brandonplank.org/neptune/embed v0.0.0
	brandonplank.org/neptune/global v0.0.0
	brandonplank.org/neptune/routes v0.0.0
	github.com/crypticplank/israilway v0.0.0-20220506200846-efecf6a999e2
	github.com/getsentry/sentry-go v0.13.0
	github.com/gofiber/fiber/v2 v2.35.0
	github.com/gofiber/template v1.6.29
	github.com/joho/godotenv v1.4.0
	github.com/mileusna/crontab v1.2.0
)

replace (
	brandonplank.org/neptune/controllers => ./controllers
	brandonplank.org/neptune/database => ./database
	brandonplank.org/neptune/embed => ./embed
	brandonplank.org/neptune/global => ./global
	brandonplank.org/neptune/models => ./models
	brandonplank.org/neptune/routes => ./routes
)

require (
	brandonplank.org/neptune/models v0.0.0 // indirect
	github.com/Cryptolens/cryptolens-golang v0.0.0-20220630131701-c1b99a2da081 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gocarina/gocsv v0.0.0-20220729221910-a7386ae0b221 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/leaanthony/debme v1.2.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.38.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10 // indirect
	gorm.io/driver/mysql v1.3.5 // indirect
	gorm.io/gorm v1.23.8 // indirect
)
