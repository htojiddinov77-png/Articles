package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUserDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}

	// run migrations for test db
	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("migrations test db error: %v", err)
	}

	// clean up old data
	_, err = db.Exec("TRUNCATE users CASCADE")
	if err != nil {
		t.Fatalf("truncation tables %v", err)
	}

	return db
}

func TestCreateUser(t *testing.T) {
	db := setupUserDB(t)
	defer db.Close()

	store := NewPostgresUserStore(db)

	tests := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &User{
				Username:     "Noah",
				Email:        "Noah_Faris@example.com",
				PasswordHash: "hash_12345",
				Bio:          "Backend developer from Tashkent who loves Go and coffee ☕️",
			},
			wantErr: false,
		},
		{
			name: "valid user 2",
			user: &User{
				Username:     "muhammad_rasul",
				Email:        "muhammad.rasul@example.com",
				PasswordHash: "secure_as!2kge",
				Bio:          "Tech enthusiast and aspiring software engineer 🚀",
			},
			wantErr: false,
		},
		{
			name: "missing username",
			user: &User{
				Email:        "invalid@example.com",
				PasswordHash: "hash_no_username",
				Bio:          "This should fail because username is missing",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdUser, err := store.CreateUser(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			
			assert.Equal(t, tt.user.Username, createdUser.Username)
			assert.Equal(t, tt.user.Email, createdUser.Email)
			assert.Equal(t, tt.user.Bio, createdUser.Bio)

			
			retrieved, err := store.GetUserById(int64(createdUser.ID))
			require.NoError(t, err)

			assert.Equal(t, createdUser.ID, retrieved.ID)
			assert.Equal(t, tt.user.Username, retrieved.Username)
			assert.Equal(t, tt.user.Email, retrieved.Email)
			assert.Equal(t, tt.user.Bio, retrieved.Bio)
		})
	}
}
