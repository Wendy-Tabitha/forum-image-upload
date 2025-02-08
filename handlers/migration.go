package handlers

import (
	"log"
)

func RunMigrations() error {
	// Add parent_id column to comments table if it doesn't exist
	_, err := db.Exec(`
		SELECT parent_id FROM comments LIMIT 1;
	`)
	if err != nil {
		log.Println("Adding parent_id column to comments table...")
		_, err = db.Exec(`
			ALTER TABLE comments ADD COLUMN parent_id INTEGER DEFAULT NULL REFERENCES comments(id) ON DELETE CASCADE;
		`)
		if err != nil {
			return err
		}
	}
	return nil
}
