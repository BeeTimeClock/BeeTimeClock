package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DatabaseManager struct {
	prefix string
}

func NewDatabaseManager(prefix string) *DatabaseManager {
	return &DatabaseManager{
		prefix: prefix,
	}
}

func (d *DatabaseManager) newConnection() (*gorm.DB, error) {
	config := &gorm.Config{}

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "psql"
	}

	var dialect gorm.Dialector

	hasMissing := false

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DATABASE")
	port := os.Getenv("DB_PORT")
	
	if host == "" {
		fmt.Println("Missing DB_HOST")
		hasMissing = true
	}
	
	if user == "" {
		fmt.Println("Missing DB_USER")
		hasMissing = true
	}
	
	if password  == "" {
		fmt.Println("Missing DB_PASSWORD")
		hasMissing = true
	}

	if database == "" {
		fmt.Println("Missing DB_DATABASE")
		hasMissing = true
	}

	if port == "" {
		fmt.Println("Missing DB_PORT")
		hasMissing = true
	}

	if hasMissing {
		return nil, fmt.Errorf("missing database env vars")
	}

	switch dbType {
	case "psql":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable application_name=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DATABASE"),
			os.Getenv("DB_PORT"),
			d.prefix,
		)

		dialect = postgres.Open(dsn)
		break
	default:
		return nil, fmt.Errorf("database type %s not supported", dbType)
	}

	config.NamingStrategy = schema.NamingStrategy{
		SingularTable: true,
		TablePrefix:   fmt.Sprintf("%s_", d.prefix),
	}

	conn, err := gorm.Open(dialect, config)
	return conn, err
}

func (d *DatabaseManager) GetConnection() (*gorm.DB, error) {
	conn, err := d.newConnection()
	if err != nil {
		return nil, err
	}

	return conn, err
}

func (d *DatabaseManager) CloseConnection(conn *gorm.DB) {
	db, _ := conn.DB()
	db.Close()
}
