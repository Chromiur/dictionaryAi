package restDb

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type User struct {
	ID        int
	Username  string
	Password  string
	LastLogin int64
	Admin     int
	Active    int
	Teacher   int
	Student   int
}

// FromJSON decodes a serialized JSON record - User{}
func (p *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

// ToJSON encodes a User JSON record
func (p *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// DeleteUser is for deleting a user defined by ID
func DeleteUser(ID int) bool {
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()

	// Check is the user ID exists
	t := FindUserByID(ID)
	if t.ID == 0 {
		log.Println("User", ID, "does not exist.")
		return false
	}

	stmt, err := db.Prepare("DELETE FROM users WHERE ID = $1")
	if err != nil {
		log.Println("DeleteUser:", err)
		return false
	}

	_, err = stmt.Exec(ID)
	if err != nil {
		log.Println("DeleteUser:", err)
		return false
	}

	return true
}

// InsertUser is for adding a new user to the database
func InsertUser(u User) bool {
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()

	if IsUserValid(u) {
		log.Println("User", u.Username, "already exists!")
		return false
	}

	stmt, err := db.Prepare("INSERT INTO users(Username, Password, LastLogin, Admin, Active, Teacher, Student) values($1,$2,$3,$4,$5,$6,$7)")
	if err != nil {
		log.Println("Adduser:", err)
		return false
	}

	stmt.Exec(u.Username, u.Password, u.LastLogin, u.Admin, u.Active, u.Teacher, u.Student)
	return true
}

// ListAllUsers is for returning all users from the database table
func ListAllUsers() []User {
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return []User{}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users \n")
	if err != nil {
		log.Println(err)
		return []User{}
	}

	all := []User{}
	var c1 int
	var c2, c3 string
	var c4 int64
	var c5, c6, c7, c8 int

	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6, &c7, &c8)
		temp := User{c1, c2, c3, c4, c5, c6, c7, c8}
		all = append(all, temp)
	}

	log.Println("All:", all)
	return all
}

// FindUserByID is for returning a user record defined by ID
func FindUserByID(ID int) User {
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return User{}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users where ID = $1\n", ID)
	if err != nil {
		log.Println("Query:", err)
		return User{}
	}
	defer rows.Close()

	u := User{}
	var c1 int
	var c2, c3 string
	var c4 int64
	var c5, c6, c7, c8 int

	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6, &c7, &c8)
		if err != nil {
			log.Println(err)
			return User{}
		}
		u = User{c1, c2, c3, c4, c5, c6, c7, c8}
		log.Println("Found user:", u)
	}
	return u
}

func IsUserValid(u User) bool {
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users WHERE Username = $1 \n", u.Username)
	if err != nil {
		log.Println(err)
		return false
	}

	temp := User{}
	var c1 int
	var c2, c3 string
	var c4 int64
	var c5, c6, c7, c8 int

	// If there exist multiple users with the same username,
	// we will get the FIRST ONE only.
	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6, &c7, &c8)
		if err != nil {
			log.Println(err)
			return false
		}
		temp = User{c1, c2, c3, c4, c5, c6, c7, c8}
	}

	if u.Username == temp.Username && u.Password == temp.Password {
		return true
	}
	return false
}

// UpdateUser allows you to update user name
func UpdateUser(u User) bool {
	log.Println("Updating user:", u)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE users SET Username=$1, Password=$2, Admin=$3, Active=$4 WHERE ID = $5")
	if err != nil {
		log.Println("Adduser:", err)
		return false
	}

	res, err := stmt.Exec(u.Username, u.Password, u.Admin, u.Active, u.ID)
	if err != nil {
		log.Println("UpdateUser failed:", err)
		return false
	}

	affect, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() failed:", err)
		return false
	}
	log.Println("Affected:", affect)
	return true
}

// IsUserAdmin determines whether a user is
// an administrator or not
func IsUserAdmin(u User) bool {
	log.Printf("Validating whether user %s admin or not...", u.Username)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()

	rows, err := db.Query("SELECT Username, Password, Admin FROM users WHERE Active = 1 AND Username = $1 \n", u.Username)
	if err != nil {
		log.Println(err)
		return false
	}

	temp := User{}
	var c1, c2 string
	var c3 int

	// If there exist multiple users with the same username,
	// we will get the FIRST ONE only.
	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3)
		if err != nil {
			log.Println(err)
			return false
		}
		temp = User{Username: c1, Password: c2, Admin: c3}
	}

	if u.Username == temp.Username && u.Password == temp.Password && temp.Admin == 1 {
		return true
	}
	return false
}

// FindUserByUsername is for returning a user record defined by a username
func FindUserByUsername(username string) User {
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return User{}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users where Username = $1 \n", username)
	if err != nil {
		log.Println("FindUserByUsername Query:", err)
		return User{}
	}
	defer rows.Close()

	u := User{}
	var c1 int
	var c2, c3 string
	var c4 int64
	var c5, c6, c7, c8 int

	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5, &c6, &c7, &c8)
		if err != nil {
			log.Println(err)
			return User{}
		}
		u = User{c1, c2, c3, c4, c5, c6, c7, c8}
		log.Println("Found user:", u)
	}
	return u
}
