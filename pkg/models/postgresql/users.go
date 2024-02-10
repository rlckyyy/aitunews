package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"relucky.net/aitunews/pkg/models"
	"strings"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password, role string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `
		INSERT INTO users (name, email, hashed_password, role)
		VALUES($1, $2, $3, $4)
		RETURNING id
	`

	row := m.DB.QueryRow(stmt, name, email, string(hashedPassword), role)

	var id int
	err = row.Scan(&id)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "23505" && strings.Contains(err.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = $1"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	stmt := "SELECT id, name, email, hashed_password, role FROM users WHERE id = $1"
	row := m.DB.QueryRow(stmt, id)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *UserModel) GetUsers() ([]*models.User, error) {
	stmt := `SELECT id, name, email, hashed_password, role FROM users`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userList := []*models.User{}

	for rows.Next() {
		n := &models.User{}

		err := rows.Scan(&n.ID, &n.Name, &n.Email, &n.Role, &n.HashedPassword)
		if err != nil {
			return nil, err
		}

		userList = append(userList, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	fmt.Println(userList)
	return userList, nil
}

func (m *UserModel) UpdateUserRole(userID, newRole string) error {
	fmt.Println(userID + " " + newRole)
	stmt := `UPDATE users SET role = $1 WHERE id = $2`

	_, err := m.DB.Exec(stmt, newRole, userID)
	if err != nil {
		return fmt.Errorf("unable to update user role: %w", err)
	}

	return nil
}
