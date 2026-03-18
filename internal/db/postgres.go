package db

import (
	"fmt"

	"github.com/IUnlimit/ssh2a/configs"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *configs.Database) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&SSHAttempt{}, &HoneypotCredential{}, &AuthRecord{})
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}
	log.Info("Database connected and migrated")
}
