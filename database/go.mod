module brandonplank.org/neptune/database

go 1.19

require (
	brandonplank.org/neptune/models v0.0.0
	gorm.io/driver/mysql v1.3.3
	gorm.io/gorm v1.23.4
)

require (
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
)

replace brandonplank.org/neptune/models => ../models
