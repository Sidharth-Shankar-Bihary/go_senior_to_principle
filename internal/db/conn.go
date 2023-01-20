package db

import (
	"fmt"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/models"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/utils"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func GetDB() *gorm.DB {
	return db
}

// GetModels includes all database structs you want to migrate.
func GetModels() []interface{} {
	return []interface{}{
		&models.User{},
	}
}

// InitDB will return an instance of gorm.DB to an application.
func InitDB() (err error) {
	switch os.Getenv("DB_DRIVER") {
	case utils.Mysql:
		db, err = setupMySQL()
	case utils.Postgres:
		db, err = setupPostgres()
	case utils.Sqlite:
		db, err = setupSQLite()
	default:
		return fmt.Errorf("no database found, set the DB baseenv")
	}

	if err != nil {
		return err
	}

	// after connect db, then do auto migrate.
	if err = db.AutoMigrate(GetModels()...); err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

func setupPostgres() (*gorm.DB, error) {
	host := os.Getenv("POSTGRESQL_HOST")
	port := os.Getenv("POSTGRESQL_PORT")
	user := os.Getenv("POSTGRESQL_USER")
	pwd := os.Getenv("POSTGRESQL_PWD")
	dbName := os.Getenv("POSTGRESQL_DB")

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		user,
		dbName,
		pwd,
	)

	db, err = gorm.Open(postgres.Open(connStr))
	if err != nil {
		log.Fatal(err.Error())
	}

	return db, nil
}

func setupMySQL() (*gorm.DB, error) {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	pwd := os.Getenv("MYSQL_PWD")
	dbName := os.Getenv("MYSQL_DB")

	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		pwd,
		host,
		port,
		dbName,
	)

	db, err = gorm.Open(mysql.Open(connStr))
	if err != nil {
		log.Fatal(err.Error())
	}
	return db, nil
}

func setupSQLite() (*gorm.DB, error) {
	dbLocation := os.Getenv("DATABASE_PATH")
	if dbLocation == "" {
		dbLocation = "/opt/auth-service/gorm.db"
	}

	// Create the sqlite file if it's not available
	if _, err = os.Stat(dbLocation); err != nil {
		if _, err = os.Create(dbLocation); err != nil {
			return nil, err
		}
	}

	db, err = gorm.Open(sqlite.Open(dbLocation), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	return db, err
}
