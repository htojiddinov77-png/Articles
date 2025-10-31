package store

import (
	"database/sql"
	"time"
)

type Article struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	AuthorId    int         `json:"author_id"`
	Paragraphs  []Paragraph `json:"paragraphs"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type Paragraph struct {
	ID         int    `json:"id"`
	Headline   string `json:"headline"`
	Body       string `json:"body"`
	OrderIndex int    `json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PostgresArticleStore struct {
	db *sql.DB
}

func NewPosgresArticleStore(db *sql.DB) *PostgresArticleStore {
	return &PostgresArticleStore{db: db}
}

type ArticleStore interface {
	CreateArticle(*Article) (*Article, error)
	GetArticleById(id int64) (*Article, error)
	UpdateArticle(*Article) error
	DeleteArticle(id int64) error
}

func (pg *PostgresArticleStore) CreateArticle(article *Article) (*Article, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query :=
		`INSERT INTO articles (title,description,image,author_id)
	VALUES($1, $2, $3, $4)
	RETURNING id`

	err = tx.QueryRow(query, article.Title, article.Description, article.Image, article.AuthorId).Scan(&article.ID)
	if err != nil {
		return nil, err
	}

	for _, paragraph := range article.Paragraphs {
		query :=
			`INSERT INTO paragraphs(article_id, headline, body,order_index)
		VALUES($1, $2, $3, $4)
		RETURNING id`

		err = tx.QueryRow(query, article.ID, paragraph.Headline, paragraph.Body, paragraph.OrderIndex).Scan(&paragraph.ID)
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

func (pg *PostgresArticleStore) GetArticleById(id int64) (*Article, error) {
	article := &Article{}
	query := `
	SELECT id, title, description,image, author_id,created_at, updated_at
	FROM articles WHERE id = $1`

	err := pg.db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Image, &article.AuthorId, &article.CreatedAt, &article.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	paragraphQuery := `
	SELECT id, headline,body,order_index,created_at,updated_at
	FROM paragraphs
	WHERE article_id = $1
	ORDER BY order_index`

	rows, err := pg.db.Query(paragraphQuery, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var entry Paragraph
		err = rows.Scan(
			&entry.ID,
			&entry.Headline,
			&entry.Body,
			&entry.OrderIndex,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		article.Paragraphs = append(article.Paragraphs, entry)
	}

	return article, nil
}

func (pg *PostgresArticleStore) UpdateArticle(article *Article) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	UPDATE articles
	SET title = $1, description = $2, image = $3, author_id = $4,updated_at = NOW()
	WHERE id = $5`

	result, err := tx.Exec(query, article.Title, article.Description, article.Image, article.AuthorId, article.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	_, err = tx.Exec(`DELETE FROM paragraphs WHERE article_id = $1`, article.ID)
	if err != nil {
		return err
	}

	for _, paragraph := range article.Paragraphs {
		query := `
		INSERT INTO paragraphs (article_id, headline, body, order_index)
		VALUES($1, $2, $3, $4);
		`
		_, err := tx.Exec(query,
			article.ID,
			paragraph.Headline,
			paragraph.Body,
			paragraph.OrderIndex,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (pg *PostgresArticleStore) DeleteArticle(id int64) error {
	query := `
	DELETE from articles
	WHERE id = $1`

	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
