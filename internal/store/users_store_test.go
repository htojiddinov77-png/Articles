// package store

// import (
// 	"testing"

// 	_ "github.com/jackc/pgx/v4/stdlib"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )



// func TestCreateUser(t *testing.T) {
// 	db := setupTestDB(t)
// 	defer db.Close()

// 	store := NewPostgresUserStore(db)

// 	tests := []struct {
// 		name    string
// 		user    *User
// 		wantErr bool
// 	}{
// 		{
// 			name: "valid user",
// 			user: &User{
// 				Username:     "Noah",
// 				Email:        "Noah_Faris@example.com",
// 				Bio:          "Backend developer from Tashkent who loves Go and coffee",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "valid user 2",
// 			user: &User{
// 				Username:     "muhammad_rasul",
// 				Email:        "muhammad.rasul@example.com",
// 				Bio:          "Tech enthusiast software engineer ",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "missing username",
// 			user: &User{
// 				Email:        "invalid@example.com",
// 				Bio:          "This should fail because username is missing",
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			createdUser, err := store.CreateUser(tt.user)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 				return
// 			}

// 			require.NoError(t, err)
			
// 			assert.Equal(t, tt.user.Username, createdUser.Username)
// 			assert.Equal(t, tt.user.Email, createdUser.Email)
// 			assert.Equal(t, tt.user.Bio, createdUser.Bio)

			
		
// 			require.NoError(t, err)

// 		})
// 	}
// }
package store