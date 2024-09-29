package models

import (
	"errors"
	"restapiproject/db"
	"restapiproject/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func GetAllUsers() ([]User, error) {
	query := "SELECT * FROM users"
	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (u *User) Save() error {
	query := `INSERT INTO users(email, password) VALUES (?,?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPwd, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPwd)

	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()

	u.ID = userId

	return err
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	retrievedPwd := ""
	err := row.Scan(&u.ID, &retrievedPwd)

	if err != nil {
		return errors.New("credential invalid")
	}

	pwdIsValid := utils.CheckPasswordHash(u.Password, retrievedPwd)

	if !pwdIsValid {
		return errors.New("credential invalid")
	}

	return nil
}
