package store

import (
	"database/sql"
	"testing"
	_ "github.com/jackc/pgx/v4/stdlib"
)


func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}

	
	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("migrations test db error: %v", err)
	}
	
	_, err = db.Exec(`
		TRUNCATE users, articles, paragraphs, workouts, workout_entries
		CASCADE;
	`)
	if err != nil {
		t.Fatalf("truncating tables: %v", err)
	}

	return db
}
