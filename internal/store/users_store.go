package store

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Bio          string    `json:"bio"`
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

func (pg *PostgresUserStore) CreateUser(user *User) (*User, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `
	INSERT INTO users(username,email, password_hash,bio, created_at, updated_at)
	VALUES ($1, $2, $3, $4, NOW(), NOW())
	RETURNING id;`

	err = tx.QueryRow(query, user.Username,user.Email, user.PasswordHash, user.Bio).Scan(&user.ID)
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
	SELECT id, username, email, password_hash, bio, created_at, updated_at
	FROM users 
	WHERE id = $1;
	`

	row := pg.db.QueryRow(query, id)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (pg *PostgresUserStore) UpdateUser(user *User) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	UPDATE users 
	SET username = $1, email = $2, password_hash = $3, bio = $4, updated_at = NOW()
	WHERE id = $5;
	`

	result, err := tx.Exec(query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Bio,
		user.ID,
	)
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

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}


func (pg *PostgresUserStore) DeleteUser(id int64) error {
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
