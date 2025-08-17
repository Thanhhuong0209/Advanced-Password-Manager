package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"password-manager/internal/crypto"

	_ "github.com/mattn/go-sqlite3"
)

// PasswordEntry represents a stored password entry
type PasswordEntry struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	URL         string    `json:"url"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []string  `json:"tags"`
}

// Database represents the encrypted password database
type Database struct {
	dbPath string
	db     *sql.DB
	masterPassword string
}

// NewDatabase creates a new database instance
func NewDatabase(dbPath, masterPassword string) (*Database, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Open SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{
		dbPath: dbPath,
		db:     db,
		masterPassword: masterPassword,
	}

	// Initialize database schema
	if err := database.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return database, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

// initSchema creates the database tables if they don't exist
func (db *Database) initSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS passwords (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			username TEXT,
			encrypted_password TEXT NOT NULL,
			url TEXT,
			notes TEXT,
			encrypted_tags TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS metadata (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_passwords_name ON passwords(name)`,
		`CREATE INDEX IF NOT EXISTS idx_passwords_username ON passwords(username)`,
	}

	for _, query := range queries {
		if _, err := db.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query '%s': %w", query, err)
		}
	}

	return nil
}

// SavePassword saves a password entry to the database
func (db *Database) SavePassword(entry *PasswordEntry) error {
	// Encrypt password
	encryptedPassword, err := crypto.Encrypt(entry.Password, db.masterPassword)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	// Encrypt tags
	encryptedTags, err := crypto.Encrypt(string(marshalTags(entry.Tags)), db.masterPassword)
	if err != nil {
		return fmt.Errorf("failed to encrypt tags: %w", err)
	}

	// Convert encrypted data to JSON
	passwordJSON, err := json.Marshal(encryptedPassword)
	if err != nil {
		return fmt.Errorf("failed to marshal encrypted password: %w", err)
	}

	tagsJSON, err := json.Marshal(encryptedTags)
	if err != nil {
		return fmt.Errorf("failed to marshal encrypted tags: %w", err)
	}

	// Insert or update password
	query := `INSERT OR REPLACE INTO passwords 
		(name, username, encrypted_password, url, notes, encrypted_tags, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	result, err := db.db.Exec(query, 
		entry.Name, 
		entry.Username, 
		string(passwordJSON), 
		entry.URL, 
		entry.Notes, 
		string(tagsJSON))
	
	if err != nil {
		return fmt.Errorf("failed to save password: %w", err)
	}

	// Get the ID if it's a new entry
	if entry.ID == 0 {
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID: %w", err)
		}
		entry.ID = id
	}

	return nil
}

// GetPassword retrieves a password entry by name
func (db *Database) GetPassword(name string) (*PasswordEntry, error) {
	query := `SELECT id, name, username, encrypted_password, url, notes, encrypted_tags, created_at, updated_at 
		FROM passwords WHERE name = ?`

	var entry PasswordEntry
	var passwordJSON, tagsJSON string
	var createdAt, updatedAt string

	err := db.db.QueryRow(query, name).Scan(
		&entry.ID,
		&entry.Name,
		&entry.Username,
		&passwordJSON,
		&entry.URL,
		&entry.Notes,
		&tagsJSON,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("password not found: %s", name)
		}
		return nil, fmt.Errorf("failed to query password: %w", err)
	}

	// Parse timestamps
	if entry.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
		entry.CreatedAt = time.Now()
	}
	if entry.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
		entry.UpdatedAt = time.Now()
	}

	// Decrypt password
	var encryptedPassword crypto.EncryptedData
	if err := json.Unmarshal([]byte(passwordJSON), &encryptedPassword); err != nil {
		return nil, fmt.Errorf("failed to unmarshal encrypted password: %w", err)
	}

	decryptedPassword, err := crypto.Decrypt(&encryptedPassword, db.masterPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}
	entry.Password = decryptedPassword

	// Decrypt tags
	var encryptedTags crypto.EncryptedData
	if err := json.Unmarshal([]byte(tagsJSON), &encryptedTags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal encrypted tags: %w", err)
	}

	decryptedTags, err := crypto.Decrypt(&encryptedTags, db.masterPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt tags: %w", err)
	}

	entry.Tags = unmarshalTags(decryptedTags)

	return &entry, nil
}

// ListPasswords returns all password entries
func (db *Database) ListPasswords() ([]*PasswordEntry, error) {
	query := `SELECT id, name, username, encrypted_password, url, notes, encrypted_tags, created_at, updated_at 
		FROM passwords ORDER BY name`

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query passwords: %w", err)
	}
	defer rows.Close()

	var entries []*PasswordEntry
	for rows.Next() {
		var entry PasswordEntry
		var passwordJSON, tagsJSON string
		var createdAt, updatedAt string

		err := rows.Scan(
			&entry.ID,
			&entry.Name,
			&entry.Username,
			&passwordJSON,
			&entry.URL,
			&entry.Notes,
			&tagsJSON,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Parse timestamps
		if entry.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
			entry.CreatedAt = time.Now()
		}
		if entry.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
			entry.UpdatedAt = time.Now()
		}

		// Decrypt password
		var encryptedPassword crypto.EncryptedData
		if err := json.Unmarshal([]byte(passwordJSON), &encryptedPassword); err != nil {
			continue // Skip invalid entries
		}

		decryptedPassword, err := crypto.Decrypt(&encryptedPassword, db.masterPassword)
		if err != nil {
			continue // Skip entries that can't be decrypted
		}
		entry.Password = decryptedPassword

		// Decrypt tags
		var encryptedTags crypto.EncryptedData
		if err := json.Unmarshal([]byte(tagsJSON), &encryptedTags); err != nil {
			entry.Tags = []string{}
		} else {
			decryptedTags, err := crypto.Decrypt(&encryptedTags, db.masterPassword)
			if err != nil {
				entry.Tags = []string{}
			} else {
				entry.Tags = unmarshalTags(decryptedTags)
			}
		}

		entries = append(entries, &entry)
	}

	return entries, nil
}

// DeletePassword deletes a password entry by name
func (db *Database) DeletePassword(name string) error {
	query := `DELETE FROM passwords WHERE name = ?`
	
	result, err := db.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to delete password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("password not found: %s", name)
	}

	return nil
}

// SearchPasswords searches for passwords by query
func (db *Database) SearchPasswords(query string) ([]*PasswordEntry, error) {
	searchQuery := `SELECT id, name, username, encrypted_password, url, notes, encrypted_tags, created_at, updated_at 
		FROM passwords WHERE name LIKE ? OR username LIKE ? OR url LIKE ? ORDER BY name`

	searchPattern := "%" + query + "%"
	rows, err := db.db.Query(searchQuery, searchPattern, searchPattern, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search passwords: %w", err)
	}
	defer rows.Close()

	var entries []*PasswordEntry
	for rows.Next() {
		var entry PasswordEntry
		var passwordJSON, tagsJSON string
		var createdAt, updatedAt string

		err := rows.Scan(
			&entry.ID,
			&entry.Name,
			&entry.Username,
			&passwordJSON,
			&entry.URL,
			&entry.Notes,
			&tagsJSON,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			continue
		}

		// Parse timestamps
		if entry.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt); err != nil {
			entry.CreatedAt = time.Now()
		}
		if entry.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt); err != nil {
			entry.UpdatedAt = time.Now()
		}

		// Decrypt password
		var encryptedPassword crypto.EncryptedData
		if err := json.Unmarshal([]byte(passwordJSON), &encryptedPassword); err != nil {
			continue
		}

		decryptedPassword, err := crypto.Decrypt(&encryptedPassword, db.masterPassword)
		if err != nil {
			continue
		}
		entry.Password = decryptedPassword

		// Decrypt tags
		var encryptedTags crypto.EncryptedData
		if err := json.Unmarshal([]byte(tagsJSON), &encryptedTags); err != nil {
			entry.Tags = []string{}
		} else {
			decryptedTags, err := crypto.Decrypt(&encryptedTags, db.masterPassword)
			if err != nil {
				entry.Tags = []string{}
			} else {
				entry.Tags = unmarshalTags(decryptedTags)
			}
		}

		entries = append(entries, &entry)
	}

	return entries, nil
}

// GetStats returns database statistics
func (db *Database) GetStats() (map[string]interface{}, error) {
	query := `SELECT COUNT(*) FROM passwords`
	
	var count int
	err := db.db.QueryRow(query).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to get password count: %w", err)
	}

	// Get file size
	fileInfo, err := os.Stat(db.dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	stats := map[string]interface{}{
		"total_passwords": count,
		"database_size":   fileInfo.Size(),
		"created_at":      fileInfo.ModTime(),
	}

	return stats, nil
}

// marshalTags converts tags slice to JSON string
func marshalTags(tags []string) []byte {
	if len(tags) == 0 {
		return []byte("[]")
	}
	
	data, err := json.Marshal(tags)
	if err != nil {
		return []byte("[]")
	}
	return data
}

// unmarshalTags converts JSON string to tags slice
func unmarshalTags(data string) []string {
	if data == "" || data == "[]" {
		return []string{}
	}
	
	var tags []string
	if err := json.Unmarshal([]byte(data), &tags); err != nil {
		return []string{}
	}
	return tags
}
