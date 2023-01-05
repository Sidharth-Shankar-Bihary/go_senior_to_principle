package db

import (
	"fmt"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/models"
	environment "github.com/ramseyjiang/go_senior_to_principle/internal/env"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func GetDB() *gorm.DB {
	return db
}

func getModels() []interface{} {
	return []interface{}{&models.User{}}
}

func setupPostgres(gEnv *environment.Environment) (*gorm.DB, error) {
	host := gEnv.C().Postgres.Host
	port := gEnv.C().Postgres.Port
	user := gEnv.C().Postgres.User
	pwd := gEnv.C().Postgres.Pwd
	dbName := gEnv.C().Postgres.DB

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

func setupMySQL(gEnv *environment.Environment) (*gorm.DB, error) {
	host := gEnv.C().Mysql.Host
	port := gEnv.C().Mysql.Port
	user := gEnv.C().Mysql.User
	pwd := gEnv.C().Mysql.Pwd
	dbName := gEnv.C().Mysql.DB

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

// InitDB will return an instance of gorm.DB to an application.
func InitDB(gEnv *environment.Environment, dbType string) (err error) {
	switch dbType {
	case "mysql":
		db, err = setupMySQL(gEnv)
	case "postgres":
		db, err = setupPostgres(gEnv)
	case "sqlite":
		db, err = setupSQLite()
	default:
		return fmt.Errorf("no database found, set the DB baseenv")
	}

	if err != nil {
		return err
	}

	// after connect db, then do auto migrate.
	if err = db.AutoMigrate(getModels()...); err != nil {
		log.Fatal(err.Error())
	}

	return nil
}
