package database

import (
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	"time"
)

type DB struct {
	*sqlx.DB
}

// New initializes a new MySQL database connection and performs migration if needed.
func New(dsn string, automigrate bool) (*DB, error) {
	var db *sqlx.DB
	var err error
	// Add retries in case of DB is not ready when app starts
	for i := 0; i < 5; i++ {
		db, err = sqlx.Connect("mysql", dsn)
		if nil == err {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)

	if automigrate {
		// Use the MySQL driver specific migration configuration
		migrate.SetTable("migrations")
		source := migrate.FileMigrationSource{
			Dir: "./internal/database/migrations",
		}

		// Run the migrations
		_, err := migrate.Exec(db.DB, "mysql", source, migrate.Up)
		if err != nil {
			return nil, err
		}
	}

	return &DB{db}, nil
}
