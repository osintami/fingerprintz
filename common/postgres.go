// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gorm_log "gorm.io/gorm/logger"
)

type PostgresConfig struct {
	PgHost     string
	PgPort     string
	PgUser     string
	PgPassword string
	PgDB       string
}

func OpenDB(cfg *PostgresConfig, logPath string) (*gorm.DB, error) {
	fh, err := os.OpenFile(logPath+"/postgres.log", os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return nil, err
	}
	myLog := gorm_log.New(log.New(fh, "\r\n", log.LstdFlags),
		gorm_log.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gorm_log.Error,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable timezone=utc",
		cfg.PgHost,
		cfg.PgPort,
		cfg.PgUser,
		cfg.PgDB,
		cfg.PgPassword)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: myLog})
}

func CreateMockDatabase() (*gorm.DB, sqlmock.Sqlmock) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	return db, mock
}
