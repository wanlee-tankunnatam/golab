package db

import (
	"fmt"
	"log"

	models "go-lab/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectPostgres() {
	// โหลด config จาก ENV (หรือจะฮาร์ดโค้ดไว้ก็ได้ใน dev)
	host := "128.199.97.209" // eg. "localhost"
	port := "5432"           // eg. "5432"
	user := "testqa"         // eg. "postgres"
	password := "testqa"     // eg. "yourpassword"
	dbname := "testdb"       // eg. "testdb"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		host, user, password, dbname, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	tables := []interface{}{&models.Users{}, &models.Stock{}}

	for _, table := range tables {
		if !DB.Migrator().HasTable(table) {
			log.Printf("Table %T not found. Migrating...", table)
			if err := DB.AutoMigrate(table); err != nil {
				log.Fatalf("failed to migrate %T: %v", table, err)
			}
			log.Printf("Table %T created.", table)
		} else {
			log.Printf("Table %T already exists. Skipping migration.", table)
		}

	}
}
