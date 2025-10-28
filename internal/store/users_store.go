package store

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FirstName    string    `json:"firstname"`
	LastName     string    `json:"lastname"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(*User) (*User, error)
	GetUserById(id int64) (*User, error)
	UpdateUser(*User) error
	DeleteUser(id int64) error
}

func (pg *PostgresArticleStore) CreateUser(user *User) (*User, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `
	INSERT INTO users(email, password_hash, firstname, lastname, created_at, updated_at)
	VALUES ($1, $2, $3, $4, NOW(), NOW())
	RETURNING id;`

	err = tx.QueryRow(query, user.Email, user.PasswordHash, user.FirstName, user.LastName).Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (pg *PostgresUserStore) GetUserById(id int64) (*User, error) {
	user := &User{}
	query := `
	SELECT id, email, password_hash,firstname,lastname, created_at, updated_at
	FROM users 
	WHERE id = $1;`

	row := pg.db.QueryRow(query, id)
	err := row.Scan(&user.ID, &user, user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (pg *PostgresArticleStore) UpdateUser(user *User) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `
	UPDATE users SET email = $1, password_hash = $2, firstname = $3, lastname= $4,updated_at = NOW()
	WHERE id = $5`

	result, err := tx.Exec(query, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", user.ID)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresArticleStore) DeleteUser(id int64) error {
	query := `
	DELETE FROM users WHERE id = $1;`
	
	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("User with ID %d not found", id)
	}
	return nil
}
