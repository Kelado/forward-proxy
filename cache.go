package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	path = "cache.db"
)

type Cache struct {
	db                *sql.DB
	expirationSeconds int
}

var cache *Cache

func initCache(expirationSeconds int) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Error opening SQLite database: %v", err)
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS cache (
        key TEXT PRIMARY KEY,
        value TEXT,
        timestamp INTEGER
    );
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Error creating cache table: %v", err)
	}
	cache = &Cache{db: db, expirationSeconds: expirationSeconds}
}

func (c *Cache) Get(key string) ([][]interface{}, bool) {
	var value string
	var timestamp int64
	query := `SELECT value, timestamp FROM cache WHERE key = ?`
	err := c.db.QueryRow(query, key).Scan(&value, &timestamp)
	if err == sql.ErrNoRows {
		return nil, false
	} else if err != nil {
		log.Printf("Error querying cache: %v", err)
		return nil, false
	}

	if time.Since(time.Unix(timestamp, 0)) > time.Duration(c.expirationSeconds)*time.Second {
		c.Delete(key)
		fmt.Println("Entry deleted after expiration")
		return nil, false
	}

	var result [][]interface{}
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		log.Printf("Error unmarshalling cache value: %v", err)
		return nil, false
	}

	return result, true
}

func (c *Cache) Set(key string, value [][]interface{}) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Printf("Error marshalling cache value: %v", err)
		return
	}

	timestamp := time.Now().Unix()
	insertOrUpdateQuery := `
    INSERT INTO cache (key, value, timestamp) VALUES (?, ?, ?)
    ON CONFLICT(key) DO UPDATE SET value = ?, timestamp = ?;
    `
	_, err = c.db.Exec(insertOrUpdateQuery, key, string(jsonValue), timestamp, string(jsonValue), timestamp)
	if err != nil {
		log.Printf("Error inserting/updating cache: %v", err)
	}
}

func (c *Cache) Delete(key string) {
	deleteQuery := `DELETE FROM cache WHERE key = ?`
	_, err := c.db.Exec(deleteQuery, key)
	if err != nil {
		log.Printf("Error deleting cache entry: %v", err)
	}
}
