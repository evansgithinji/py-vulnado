package persistence

import (
	"database/sql"
	"fmt"

	"goapp/internal/domain/entity"
)

// SQLiteUserRepository implements UserRepository with SQLite
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository creates a new SQLiteUserRepository
func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

// FindByCredentials authenticates a user
// VULNERABLE: SQL Injection via string formatting
func (r *SQLiteUserRepository) FindByCredentials(username, password string) (*entity.User, error) {
	// VULNERABLE: Direct string concatenation in SQL query
	query := fmt.Sprintf(
		"SELECT id, username, email, role, is_admin FROM users WHERE username='%s' AND password='%s'",
		username, password,
	)

	var user entity.User
	err := r.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Search searches users by term
// VULNERABLE: SQL Injection via LIKE clause
func (r *SQLiteUserRepository) Search(term string) ([]entity.User, error) {
	// VULNERABLE: Direct string concatenation
	query := fmt.Sprintf(
		"SELECT id, username, email, role, is_admin FROM users WHERE username LIKE '%%%s%%' OR email LIKE '%%%s%%'",
		term, term,
	)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.IsAdmin); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// FindByID finds user by ID
func (r *SQLiteUserRepository) FindByID(id int64) (*entity.User, error) {
	var user entity.User
	err := r.db.QueryRow(
		"SELECT id, username, email, role, is_admin FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateProfile updates user profile
// VULNERABLE: SQL Injection + Privilege Escalation via is_admin
func (r *SQLiteUserRepository) UpdateProfile(userID int64, email string, isAdmin bool) error {
	// VULNERABLE: Direct string concatenation allows SQL injection
	// VULNERABLE: is_admin parameter allows privilege escalation
	isAdminInt := 0
	if isAdmin {
		isAdminInt = 1
	}

	query := fmt.Sprintf(
		"UPDATE users SET email = '%s', is_admin = %d WHERE id = %d",
		email, isAdminInt, userID,
	)

	_, err := r.db.Exec(query)
	return err
}

// Create creates a new user
func (r *SQLiteUserRepository) Create(user *entity.User) error {
	result, err := r.db.Exec(
		"INSERT INTO users (username, password, email, role, is_admin) VALUES (?, ?, ?, ?, ?)",
		user.Username, user.Password, user.Email, user.Role, user.IsAdmin,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id

	return nil
}
