package infrastructure

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMODE  string
}

func (dbConfig *DbConfig) Read() {
	dbConfig.Host = os.Getenv("DB_HOST")
	dbConfig.User = os.Getenv("DB_USER")
	dbConfig.Password = os.Getenv("DB_PASSWORD")
	dbConfig.DBName = os.Getenv("DB_NAME")
	dbConfig.Port = os.Getenv("DB_PORT")
	dbConfig.SSLMODE = os.Getenv("DB_SSLMODE")
}

type GormPostgres interface {
	GetConnection() *gorm.DB
}

type gormPostgresImpl struct {
	master *gorm.DB
}

func NewGormPostgres() GormPostgres {
	return &gormPostgresImpl{master: connect()}
}

func connect() *gorm.DB {
	var dbConfig = DbConfig{}
	dbConfig.Read()

	connectionStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMODE)

	db, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{})

	if err != nil {
		log.Fatalln("DB error when connecting: ", err)
	}

	return db
}

func (g *gormPostgresImpl) GetConnection() *gorm.DB {
	return g.master
}
