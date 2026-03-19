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
	ensureDatabase(cfg)

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

// ensureDatabase 连接默认 postgres 库，若目标数据库不存在则自动创建
func ensureDatabase(cfg *configs.Database) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password,
	)
	tmp, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to postgres for database creation: %v", err)
	}

	var count int64
	tmp.Raw("SELECT count(1) FROM pg_database WHERE datname = ?", cfg.DBName).Scan(&count)
	if count == 0 {
		// CREATE DATABASE 不支持参数化，这里拼接是安全的（值来自配置文件）
		if err := tmp.Exec(fmt.Sprintf("CREATE DATABASE %q", cfg.DBName)).Error; err != nil {
			log.Fatalf("Failed to create database %s: %v", cfg.DBName, err)
		}
		log.Infof("Database %s created", cfg.DBName)
	}

	sqlDB, _ := tmp.DB()
	sqlDB.Close()
}
