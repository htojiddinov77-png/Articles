package store

import (
	"database/sql"
	"time"
	"github.com/htojiddinov77-png/Articles/internal/tokens"
)


type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore{
	return &PostgresTokenStore{db: db,}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	GetTokenByHash(hash []byte) (*tokens.Token, error)
	CreateNewToken(userId int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, scope string) error
}

func (t *PostgresTokenStore) CreateNewToken(userID int, ttl time.Duration, scope string)(*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil,err
	}
	err = t.Insert(token)
	return token,err
}

func (t *PostgresTokenStore) Insert(token *tokens.Token) error{
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES($1, $2, $3, $4)`

	_, err := t.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	return err
}

func (t *PostgresTokenStore) GetTokenByHash(hash []byte) (*tokens.Token, error) {
	token := &tokens.Token{}
	query := `
	SELECT user_id, expiry, scope
	FROM tokens
	WHERE hash = $1`

	row := t.db.QueryRow(query, hash)

	err := row.Scan(
		&token.UserID,
		&token.Expiry,
		&token.Scope,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	token.Hash = hash
	return token, nil
}



func (t *PostgresTokenStore) DeleteAllTokensForUser(userID int, scope string) error{ 
	query := `
	DELETE FROM tokens 
	WHERE scope = $1 AND user_id = $2`

	_,err := t.db.Exec(query, scope, userID)
	return err
}