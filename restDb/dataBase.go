package restDb

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

// Connection details
const (
	Hostname = "localhost"
	Port     = 5432
	Username = "ivankhromin"
	Password = "budva2024"
	Database = "dictionary_db"
)

// Userdata is for holding full word data
// Userdata table + Username
type WordData struct {
	ID          int `json:"id"`
	UserId      int
	Word        string `json:"word"`
	Description string `json:"description"`
}

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

// The function returns the Word ID of the word
// -1 if the word does not exist
func exists(word string) int {
	word = strings.ToLower(word)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	ID := -1

	statement := fmt.Sprintf(`SELECT id FROM words where word = '%s'`, word)
	fmt.Println(statement)
	rows, err := db.Query(statement)
	if err != nil {
		fmt.Println("Error: ", err)
		return -1
	}
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Scan", err)
			return -1
		}
		ID = id
	}
	defer rows.Close()
	return ID
}

// AddWord adds a new word to the database
// Returns new Word ID
// -1 if there was an error
func AddWord(d WordData) int {
	d.Word = strings.ToLower(d.Word)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	//fmt.Println(d)
	wordID := exists(d.Word)
	if wordID != -1 {
		fmt.Println("Word already exists:", Username)
		return -1
	}

	insertStatement := `insert into "words" (word) values ($1)`
	fmt.Println(insertStatement)
	_, err = db.Exec(insertStatement, d.Word)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	wordID = exists(d.Word)
	if wordID == -1 {
		return wordID
	}

	fmt.Println("TEST")

	insertStatement = `insert into "wordslist" ("id", "userid", "word", "description")
	values ($1, $2, $3, $4)`
	_, err = db.Exec(insertStatement, wordID, d.UserId, d.Word, d.Description)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1
	}

	return wordID
}

// DeleteWords deletes an existing word
func DeleteWords(ids []int) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Does the ID exist?
	for _, id := range ids {
		statement := fmt.Sprintf(`SELECT word FROM words where id = '%d'`, id)
		rows, err := db.Query(statement)
		if err != nil {
			fmt.Println("Error: ", err)
			return err
		}

		fmt.Println(statement)
		fmt.Println(rows)

		var word string
		for rows.Next() {
			err = rows.Scan(&word)
			if err != nil {
				return err
			}
			if exists(word) != id {
				return fmt.Errorf("Word with ID %d does not exist", id)
			}
			// Delete from Userdata
			deleteStatement := `delete from "wordslist" where id=$1`
			_, err = db.Exec(deleteStatement, id)
			if err != nil {
				return err
			}

			// Delete from Users
			deleteStatement = `delete from "words" where id=$1`
			_, err = db.Exec(deleteStatement, id)
			if err != nil {
				return err
			}
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {

			}
		}(rows)
	}

	return nil
}

// ListUsers lists all users in the database
func ListWords() ([]WordData, error) {
	Data := []WordData{}
	db, err := openConnection()

	if err != nil {
		return Data, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, word, description FROM wordslist \n")
	if err != nil {
		fmt.Println("Error DB:", err)
		return Data, err
	}

	for rows.Next() {
		var id int
		//var userId int
		var word string
		var description string
		err = rows.Scan(&id, &word, &description)
		temp := WordData{ID: id, Word: word, Description: description}
		Data = append(Data, temp)
		if err != nil {
			return Data, err
		}
	}
	defer rows.Close()
	return Data, nil
}
