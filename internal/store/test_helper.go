package store

import (
	"database/sql"
	"testing"
	_ "github.com/jackc/pgx/v4/stdlib"
)


func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}

	
	err = Migrate(db, "../../internal/migrations/")
	if err != nil {
		t.Fatalf("migrations test db error: %v", err)
	}
	
	_, err = db.Exec(`TRUNCATE users, articles, paragraphs RESTART IDENTITY CASCADE;`) // Restart identity har doim id ni 1 dan boshlashlikin ta'minlaydi
	
	// Truncate va Delete farqi, delete har bir rowni bittalab o'chiradi, truncate birdaniga

	if err != nil {
		t.Fatalf("truncating tables: %v", err)
	}

	return db
}
