package database

import (
	"fmt"
	"os"
	"time"
)

type User struct {
	Id            string
	Username      string
	Email         string
	Date_of_birth time.Time
}

type NewUser struct {
	Username      string    `json:"Username"`
	Email         string    `json:"Email"`
	Password      string    `json:"Password"`
	Date_of_birth time.Time `json:"Date_of_birth"`
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
	err := pool.QueryRow(ctx, "SELECT id,username,email,date_of_birth from users WHERE id = $1", id).Scan(
		&user.Id, &user.Username, &user.Email, &user.Date_of_birth,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed!: %v\n", err)
		return User{}, err
	}

	return user, nil
}

func CreateUser(user NewUser) {
	sql := `INSERT INTO users(username, email, password, date_of_birth) VALUES ($1,$2,$3,$4)`

	ctx, pool := createConnection()

	_, err := pool.Exec(ctx, sql, user.Username, user.Email, user.Password, user.Date_of_birth)
	inserQueryFail(err, "creating user")
}
