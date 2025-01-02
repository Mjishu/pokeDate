package database

import (
	"fmt"
	"os"
	"time"
)

type User struct {
	Id            string
	Username      string
	HashPassword  string
	Email         string
	Date_of_birth time.Time
}

type NewUser struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	// Email         string    `json:"Email"`
	// Date_of_birth time.Time `json:"Date_of_birth"`
}

// id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
// username VARCHAR(40) NOT NULL,
// email VARCHAR(100) ,
// password VARCHAR(50) NOT NULL,
// date_of_birth DATE,
// country_id INT REFERENCES locations(id) ON DELETE SET NULL,
// state_id INT REFERENCES locations(id) ON DELETE SET NULL,
// city_id INT REFERENCES locations(id) ON DELETE SET NULL,
// profile_pi

func GetUser(id any) (User, error) {
	ctx, pool := createConnection()
	var user User
	err := pool.QueryRow(ctx, "SELECT id,username,password,email,date_of_birth from users WHERE id = $1", id).Scan(
		&user.Id, &user.Username, &user.HashPassword, &user.Email, &user.Date_of_birth,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed!: %v\n", err)
		return User{}, err
	}

	return user, nil
}

func CreateUser(user NewUser, hashedPassword string) {
	sql := `INSERT INTO users(username,password) VALUES ($1,$2)`

	ctx, pool := createConnection()

	_, err := pool.Exec(ctx, sql, user.Username, hashedPassword) //add other options for new user like dob and email
	inserQueryFail(err, "creating user")
}
