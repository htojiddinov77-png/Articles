package store

import (
	"database/sql"
	"fmt"
	"time"
)

type Review struct {
	ID         int       `json:"id"`
	UserId     int       `json:"user_id"`
	ArticleId  int       `json:"article_id"`
	ReviewText string    `json:"review_text"`
	Rating     int       `json:"rating"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PostgresReviewStore struct {
	db *sql.DB
}

func NewPostgresReviewStore(db *sql.DB) *PostgresReviewStore {
	return &PostgresReviewStore{db: db}
}

type ReviewStore interface {
	CreateReview(*Review) (*Review, error)
	GetReviewById(id int64) (*Review, error)
	UpdateReview(*Review) error
	DeleteReview(id int64) error
}

func (pg *PostgresReviewStore) CreateReview(review *Review) (*Review, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction %v", err)
	}
	tx.Rollback()

	query := `
	INSERT INTO reviews(user_id, article_id, review_text, rating, created_at, updated_at)
	VALUES($1, $2, $3, $4 NOW(), NOW())
	RETURNING id;`

	err = tx.QueryRow(query, review.UserId, review.ArticleId, review.ReviewText, review.Rating).Scan(&review.ID)
	if err != nil {
		return nil, err
	}
	
	return review, tx.Commit()
}

func (pg *PostgresReviewStore) GetReviewById(id int64) (*Review, error) {
	review := &Review{}
	query := `
	SELECT * from reviews
	WHERE id = $;`

	row := pg.db.QueryRow(query, id)
	err := row.Scan(
		&review.ID,
		&review.ReviewText,
		&review.Rating,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (pg *PostgresReviewStore) UpdateReview(review *Review) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil
	}
	tx.Rollback()

	query := `UPDATE reviews
	SET review_text = $1, rating = $2, updated_at = NOW()
	WHERE id = $3`

	result, err := tx.Exec(query, review.ReviewText, review.Rating, review.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("review with ID %d not found", review.ID)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (pg *PostgresReviewStore) DeleteReview(id int64) error {
	query := `
	DELETE FROM reviews WHERE id = $1;`

	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Review with ID %d not found", id)
	}
	return nil
}
