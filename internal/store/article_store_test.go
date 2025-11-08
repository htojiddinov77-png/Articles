package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestCreate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresArticleStore(db)
	userStore := NewPostgresUserStore(db)
	createdUser, err := userStore.CreateUser(&User{
		Email:        "c2KQw@example.com",
		PasswordHash: "hashed_password",
		Username:    "JohnDoe",
	  })
	  require.NoError(t, err) // t orqali testing packagega kiradi // err chiqmasligini tekshiradi
	  require.NotNil(t, createdUser)	// user create bo'lganini tekshiradi
	

	  
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
				AuthorId: createdUser.ID,
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

	for _, tt := range tests { //
		t.Run(tt.name, func(t *testing.T) { // har bir test xolatini alohidadan check qilish uchun
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
			require.NoError(t, err) // agar db dan articleni olish muammo bo'lsa, test to'xtaydi

			assert.Equal(t, createdArticle.ID, retrieved.ID)
			assert.Equal(t, len(tt.article.Paragraphs), len(retrieved.Paragraphs))

			for i := range retrieved.Paragraphs{
				assert.Equal(t, tt.article.Paragraphs[i].Headline, retrieved.Paragraphs[i].Headline)
				assert.Equal(t, tt.article.Paragraphs[i].OrderIndex, retrieved.Paragraphs[i].OrderIndex)
			}

		})
	}
}