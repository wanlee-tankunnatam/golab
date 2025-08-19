package config

import (
	"fmt"
	models "go-lab/model"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectPostgres() {
	host := "128.199.97.209"
	port := "5432"
	user := "testqa"
	password := "testqa"
	dbname := "testdb"

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		host, user, password, dbname, port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	// ตั้งค่า connection pool
	sqlDB.SetMaxOpenConns(50)                 // connection สูงสุด
	sqlDB.SetMaxIdleConns(10)                 // idle connection
	sqlDB.SetConnMaxLifetime(time.Minute * 5) // ปิด connection หลัง 5 นาที

	log.Println("Connected to PostgreSQL with connection pool")

	// Migrate tables
	tables := []interface{}{&models.Stock{}}
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
