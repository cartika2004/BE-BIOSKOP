package database

import (
	"context"
	"fmt"
	"log"
	"project-bioskop/config"
	"project-bioskop/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RDB *redis.Client

func ConnectDB() {
	// Connection String khusus SQL Server
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.DBName,
	)

	var err error
	DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(" Gagal konek ke SQL Server:", err)
	}
	log.Println(" Sukses konek ke SQL Server!")

	// AUTO MIGRATE
	log.Println(" Sedang membuat tabel otomatis...")
	err = DB.AutoMigrate(
		&models.User{},
		&models.Movie{},
		&models.Studio{},
		&models.Seat{},
		&models.Schedule{},
		&models.Transaction{},
		&models.Ticket{},
	)
	if err != nil {
		log.Fatal(" Gagal migrasi tabel:", err)
	}
	log.Println(" Semua tabel berhasil dibuat!")
}

func ConnectRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		log.Println(" Redis tidak terdeteksi (Gak masalah, lanjut aja)", err)
	} else {
		log.Println(" Sukses konek ke Redis!")
	}
}