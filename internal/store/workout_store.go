package store

import (
	"database/sql"
)

type Article struct {
	ID          int          `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Image       string       `json:"image"`
	AuthorID    int          `json:"author_id"`
	Paraghraps  []Paraghraph `json:"paraghraps"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
}

type Paraghraph struct {
	ID         int    `json:"id"`
	Headline   string `json:"headline"`
	Body       string `json:"body"`
	OrderIndex int    `json:"order_index"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type PostgresArticleStore struct {
	db *sql.DB
}

func NewPostgresArticleStore(db *sql.DB) *PostgresArticleStore {
	return &PostgresArticleStore{db: db}
}

type ArticleStore interface {
	CreatedArticle(article *Article) (*Article, error)
	GetArticleByID(id int64) (*Article, error)
}

func (pg *PostgresArticleStore) CreateArticle(article *Article) (*Article, error) {
	tx, err := pg.db.Begin()
	if err != nil {

	}
	defer tx.Rollback()

	query :=
		`INSERT INTO articles (title, description, image, author_id)
	VALUES ($1, $2, $3, $4) RETURNING id;
	`
	err = tx.QueryRow(query, article.Title, article.Description, article.Image, article.AuthorID).Scan(&article.ID)
	if err != nil {
		return nil, err
	}
	for _, paraghraph := range article.Paraghraps {
		query :=
			`INSERT INTO paraghraps (article_id, headline, body, order_index)
		VALUES ($1, $2, $3, $4);
		`
		err = tx.QueryRow(query, article.ID, paraghraph.Headline, paraghraph.Body, paraghraph.OrderIndex).Scan(&paraghraph.ID)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (pg *PostgresArticleStore) GetArticleByID(id int64) (*Article, error) {
	Article := &Article{}
	return Article, nil
}
