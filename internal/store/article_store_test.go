package store

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}


	err = Migrate(db, "../..migrations/")
	if err != nil {
		t.Fatalf("migrations test db error: %v", err)
	}

	_, err = db.Exec("TRUNCATE articles, paragraphs CASCADE")
	if err != nil{
		t.Fatalf("truncation tables %v", err)
	}

	return db
}

func TestCreate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store :=NewPosgresArticleStore(db)

	tests := []struct {
		name string
		article *Article
		wantErr bool
	}{
		{
			name: "valid article",
			article: &Article{
				Title: "How to test your code",
				Description: "A simple tutorial for beginners",
				Image: "https://example.com/go",
				AuthorId: 1,
				Paragraphs: []Paragraph{
					{
						Headline: "Introduction",
						Body: "This section explain the basics of testing your apis",
						OrderIndex: 1,
					},
					{
						Headline: "Implementation",
						Body: "Here we walk through an example",
						OrderIndex: 2,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "article with invalid author id",
			article: &Article{
				Title: "Invalid Author",
				Description: "Should fail because author_id doesn't exist",
				Image: "https://example.com/invalid",
				AuthorId: 9999,
				Paragraphs: []Paragraph{
					{
						Headline: "Test Headline",
						Body: "Body content",
						OrderIndex: 1,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdArticle, err := store.CreateArticle(tt.article)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.article.Title, createdArticle.Title)
			assert.Equal(t, tt.article.Description, createdArticle.Description)
			assert.Equal(t, tt.article.Paragraphs, createdArticle.Paragraphs)

			retrieved, err := store.GetArticleById(int64(createdArticle.ID))
			require.NoError(t, err)

			assert.Equal(t, createdArticle.ID, retrieved.ID)
			assert.Equal(t, len(tt.article.Paragraphs), len(retrieved.Paragraphs))

			for i := range retrieved.Paragraphs{
				assert.Equal(t, tt.article.Paragraphs[i].Headline, retrieved.Paragraphs[i].Headline)
				assert.Equal(t, tt.article.Paragraphs[i].OrderIndex, retrieved.Paragraphs[i].OrderIndex)
			}

		})
	}
}