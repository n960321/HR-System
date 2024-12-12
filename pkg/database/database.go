package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	db, err := gorm.Open(mysql.Open(config.GetDSN()), &gorm.Config{})

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
