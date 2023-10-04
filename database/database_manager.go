package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
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
		dbType = "sqlite"
	}

	var dialect gorm.Dialector

	switch dbType {
	case "sqlite":
		dialect = sqlite.Open(fmt.Sprintf("%s.db", d.prefix))
		break
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
