package database

import (
	"context"
	"fmt"
	"log"
	"project-bioskop/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlserver"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RDB *redis.Client

func ConnectDB() {
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
		log.Fatal("Gagal konek ke SQL Server:", err)
	}
	log.Println("Sukses konek ke SQL Server!")

	runAutoMigration()
}

func runAutoMigration() {
	log.Println("Checking Migrations...")

	migrationDSN := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.DBName,
	)

	m, err := migrate.New(
		"file://database/migrations",
		migrationDSN,
	)
	if err != nil {
		log.Fatal("Gagal setup migrasi:", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Database sudah up-to-date")
		} else {
			log.Fatal("Gagal menjalankan migrasi:", err)
		}
	} else {
		log.Println("Migrasi Berhasil dijalankan otomatis!")
	}
}

func ConnectRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Redis tidak terdeteksi (Gak masalah, lanjut aja)", err)
	} else {
		log.Println("Sukses konek ke Redis!")
	}
}