package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	sqlDb  *sql.DB
	gormDb *gorm.DB
}

type Config struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"db_name"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

func (c Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
	)
}

func NewDatabase(config Config) *Database {
	newLogger := logger.New(
		&log.Logger,
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color

		},
	)

	db, err := gorm.Open(mysql.Open(config.GetDSN()), &gorm.Config{Logger: newLogger})

	if err != nil {
		log.Panic().Err(err).Msgf("Connect to Database failed")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Panic().Err(err).Msgf("failed to get DB")
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)

	log.Info().Msgf("Connect to Database [%v] Successful!", config.GetDSN())

	return &Database{
		gormDb: db,
		sqlDb:  sqlDB,
	}
}

func (db *Database) Shutdown(ctx context.Context) {
	if err := db.sqlDb.Close(); err != nil {
		log.Panic().Err(err).Msgf("failed to calse DB")
	}
}

func (db *Database) GetGorm() *gorm.DB {
	return db.gormDb
}
