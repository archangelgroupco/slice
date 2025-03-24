package vector

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	_ "github.com/lib/pq"
)

const (
	MAX_WORD_LENGTH = 1000      // Filter out words longer than 1000 chars
	MAX_TEXT_SIZE   = 1_000_000 // Truncate text over 1MB
	BATCH_SIZE      = 100       // Process 100 rows per batch
)

// filterText removes long words and truncates oversized text
func filterText(content string) string {
	if len(content) > MAX_TEXT_SIZE {
		content = content[:MAX_TEXT_SIZE]
	}
	words := strings.Fields(content)
	var filteredWords []string
	for _, word := range words {
		if len(word) <= MAX_WORD_LENGTH {
			filteredWords = append(filteredWords, word)
		}
	}
	return strings.Join(filteredWords, " ")
}

func ParseColToVec() {
	// Connect to the database
	db, err := sql.Open("postgres", os.Getenv("DB_CONN_STRING"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// TODO: create handling of file if not exists

	file, err := os.Open("documents_data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	header := true
	var wg sync.WaitGroup
	batch := make([][2]string, 0, BATCH_SIZE)

	// Read and process the CSV
	for scanner.Scan() {
		if header { // Skip header row
			header = false
			continue
		}
		line := scanner.Text()
		parts := strings.SplitN(line, ",", 2) // Split into id and content
		id := parts[0]
		content := parts[1]
		filteredContent := filterText(content)
		batch = append(batch, [2]string{id, filteredContent})

		// When batch is full, process it concurrently
		if len(batch) == BATCH_SIZE {
			wg.Add(1)
			go updateBatch(db, batch, &wg)
			batch = make([][2]string, 0, BATCH_SIZE)
		}
	}

	// Process any remaining rows
	if len(batch) > 0 {
		wg.Add(1)
		go updateBatch(db, batch, &wg)
	}

	wg.Wait() // Wait for all batches to complete
}

// updateBatch updates the database with a batch of rows
func updateBatch(db *sql.DB, batch [][2]string, wg *sync.WaitGroup) {
	defer wg.Done()
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Transaction error:", err)
		return
	}
	stmt, err := tx.Prepare("UPDATE documents SET content_tsv = to_tsvector('english', $1) WHERE id = $2")
	if err != nil {
		fmt.Println("Prepare error:", err)
		return
	}
	for _, item := range batch {
		_, err := stmt.Exec(item[1], item[0])
		if err != nil {
			fmt.Println("Exec error:", err)
		}
	}
	tx.Commit()
}
