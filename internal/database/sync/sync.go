package sync

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	SQLiteDBPath string
	PostgresDSN  string
}

// sanitizeString removes null bytes from a string
func sanitizeString(s string) string {
	return strings.Map(func(r rune) rune {
		if r == 0 {
			return -1 // Will be removed
		}
		return r
	}, s)
}

// sanitizeNullString removes null bytes from a nullable string
func sanitizeNullString(s *string) *string {
	if s == nil {
		return nil
	}
	cleaned := sanitizeString(*s)
	return &cleaned
}

type Record struct {
	ID                int64   `json:"id"`
	FilePath          string  `json:"file_path"`
	FileType          *string `json:"file_type,omitempty"`
	Content           string  `json:"content"`
	ContentHash       *string `json:"content_hash,omitempty"`
	ExtractionVersion *int64  `json:"extraction_version,omitempty"`
	Urls              *string `json:"urls,omitempty"`
	Names             *string `json:"names,omitempty"`
	Tokens            *string `json:"tokens,omitempty"`
	Places            *string `json:"places,omitempty"`
	Metadata          *string `json:"metadata,omitempty"`
}

func SyncWithRemote(sqlitePath string) {
	config := Config{
		SQLiteDBPath: sqlitePath,
		PostgresDSN:  os.Getenv("DB_CONN_STRING"),
	}

	// Connect to SQLite
	sqliteDB, err := sql.Open("sqlite3", config.SQLiteDBPath)
	if err != nil {
		log.Fatal("Failed to connect to SQLite:", err)
	}
	defer sqliteDB.Close()

	// Connect to PostgreSQL
	pgDB, err := sql.Open("postgres", config.PostgresDSN)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	defer pgDB.Close()

	// Verify connections
	if err = sqliteDB.Ping(); err != nil {
		log.Fatal("SQLite ping failed:", err)
	}
	if err = pgDB.Ping(); err != nil {
		log.Fatal("PostgreSQL ping failed:", err)
	}

	// Migrate data
	err = migrateData(sqliteDB, pgDB)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	fmt.Println("Database migration completed successfully!")
}

func migrateData(sqliteDB, pgDB *sql.DB) error {
	rows, err := sqliteDB.Query("SELECT * FROM documents")
	if err != nil {
		return fmt.Errorf("failed to query SQLite: %v", err)
	}
	defer rows.Close()

	tx, err := pgDB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	stmt, err := tx.Prepare(`
        INSERT INTO documents (id, file_path, file_type, content, content_hash, extraction_version, urls, names, tokens, places, metadata)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        ON CONFLICT (id) DO UPDATE
        SET file_path = EXCLUDED.file_path,
            file_type = EXCLUDED.file_type,
            content = EXCLUDED.content,
            content_hash = EXCLUDED.content_hash,
            extraction_version = EXCLUDED.extraction_version,
	    urls = EXCLUDED.urls,
	    names = EXCLUDED.names,
	    tokens = EXCLUDED.tokens,
	    places = EXCLUDED.places,
            metadata = EXCLUDED.metadata
    `)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	batchSize := 1000
	count := 0
	var record Record

	for rows.Next() {
		err = rows.Scan(&record.ID,
			&record.FilePath,
			&record.FileType,
			&record.Content,
			&record.ContentHash,
			&record.ExtractionVersion,
			&record.Urls,
			&record.Names,
			&record.Tokens,
			&record.Places,
			&record.Metadata)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to scan row: %v", err)
		}

		// Sanitize the data before insertion
		_, err = stmt.Exec(record.ID,
			sanitizeString(record.FilePath),
			sanitizeNullString(record.FileType),
			sanitizeString(record.Content),
			sanitizeNullString(record.ContentHash),
			record.ExtractionVersion,
			sanitizeNullString(record.Urls),
			sanitizeNullString(record.Names),
			sanitizeNullString(record.Tokens),
			sanitizeNullString(record.Places),
			sanitizeNullString(record.Metadata))
		if err != nil {
			tx.Rollback()
			log.Printf("Error processing record ID %d: %v\nFilePath: %s\n", record.ID, err, record.FilePath)
			return fmt.Errorf("failed to insert/update record %d: %v", record.ID, err)
		}

		count++
		if count%batchSize == 0 {
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("failed to commit batch: %v", err)
			}
			tx, err = pgDB.Begin()
			if err != nil {
				return fmt.Errorf("failed to begin new transaction: %v", err)
			}
			stmt, err = tx.Prepare(`
                INSERT INTO documents (id, file_path, file_type, content, content_hash, extraction_version, urls, names, tokens, places, metadata)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
                ON CONFLICT (id) DO UPDATE
                SET file_path = EXCLUDED.file_path,
                    file_type = EXCLUDED.file_type,
                    content = EXCLUDED.content,
                    content_hash = EXCLUDED.content_hash,
                    extraction_version = EXCLUDED.extraction_version,
		    urls = EXCLUDED.urls,
		    names = EXCLUDED.names,
		    tokens = EXCLUDED.tokens,
		    places = EXCLUDED.places,
                    metadata = EXCLUDED.metadata
            `)
			if err != nil {
				return fmt.Errorf("failed to prepare statement: %v", err)
			}
			fmt.Printf("Processed %d records\n", count)
		}
	}

	if err = rows.Err(); err != nil {
		tx.Rollback()
		return fmt.Errorf("row iteration error: %v", err)
	}

	return tx.Commit()
}
