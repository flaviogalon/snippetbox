package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (*User, error)
	PasswordUpdate(id int, currentPassword, newPassword string) error
}

type UserModel struct {
	DB *sql.DB
}

// Create a user
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	sql := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(sql, name, email, string(hashedPassword))
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) {
			// Check for duplicate email
			// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html#error_er_dup_entry
			if mysqlError.Number == 1062 &&
				strings.Contains(mysqlError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// Authenticate an user
func (m *UserModel) Authenticate(email, password string) (int, error) {
	user := User{}
	sqlQuery := "SELECT id, hashed_password FROM users where email = ?"

	err := m.DB.QueryRow(sqlQuery, email).Scan(&user.ID, &user.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return user.ID, nil
}

// Return true if there's a user with the given id
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	sqlQuery := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := m.DB.QueryRow(sqlQuery, id).Scan(&exists)
	return exists, err
}

// Get a new user
func (m *UserModel) Get(id int) (*User, error) {
	var user User
	sqlQuery := "SELECT id, name, email, created FROM users where id = ?"

	err := m.DB.QueryRow(sqlQuery, id).Scan(&user.ID, &user.Name, &user.Email, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return &user, nil
}

// Update an user's password
func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	var currentHashedPassword []byte

	sqlQuery := "SELECT hashed_password FROM users WHERE id = ?"
	err := m.DB.QueryRow(sqlQuery, id).Scan(&currentHashedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(currentHashedPassword, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		} else {
			return err
		}
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(currentPassword), 12)
	if err != nil {
		return err
	}

	sqlQuery = "UPDATE users SET hashed_password = ? WHERE id = ?"
	_, err = m.DB.Exec(sqlQuery, hashedNewPassword, id)
	return err
}
