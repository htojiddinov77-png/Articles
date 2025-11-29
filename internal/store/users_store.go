package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintText *string
	hash       []byte
}

func (p *password) Set(plaintTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintTextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintText = &plaintTextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err // internal server error
		}
	}

	return true, nil
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
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
	CreateUser(*User) error
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int64) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
	DeleteUser(id int64) error
}

func (pg *PostgresUserStore) CreateUser(user *User) error {
	query := `
    INSERT INTO users (username, email, password_hash, bio, created_at, updated_at)
    VALUES ($1, $2, $3, $4, NOW(), NOW())
    RETURNING id;
    `
	err := pg.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash, user.Bio).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresUserStore) GetUserByEmail(email string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}
	query := `SELECT id, username, email, bio, created_at, updated_at
	FROM users
	WHERE email = $1`

	err := pg.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (pg *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}
	query := `SELECT id, username, password_hash, email, bio, created_at, updated_at
	FROM users
	WHERE username = $1`

		err := pg.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash.hash,
		&user.Email,
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





func (pg *PostgresUserStore) GetUserById(id int64) (*User, error) {
	user := &User{}
	query := `
	SELECT id, username, password_hash, email, bio, created_at, updated_at
	FROM users 
	WHERE id = $1;
	`

	row := pg.db.QueryRow(query, id)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash.hash,
		&user.Email,
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
		user.PasswordHash.hash,
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




