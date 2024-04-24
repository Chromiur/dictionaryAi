package restDb

import (
	"database/sql"
	"fmt"
)

// Connection details
const (
	Hostname = "localhost"
	Port     = 5432
	Username = "ivankhromin"
	Password = "budva2024"
	Database = "dictionary_db"
)

func openConnection() (*sql.DB, error) {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Hostname, Port, Username, Password, Database)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	// check db
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
