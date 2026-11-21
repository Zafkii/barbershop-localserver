package db

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

// Open opens a SQLite database at path, enables foreign keys, and runs migrations.
func Open(path string) (*sql.DB, error) {
	database, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if _, err := database.Exec("PRAGMA foreign_keys = ON"); err != nil {
		database.Close()
		return nil, err
	}

	if err := Migrate(database); err != nil {
		database.Close()
		return nil, err
	}

	return database, nil
}

// Migrate creates the barber, client, service, sale, purchase, and admin tables
// if they do not exist yet. It is idempotent and never deletes data.
func Migrate(database *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS barber (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			phone TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS client (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			phone TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS service (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			service_type TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sale (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			client_id INTEGER NOT NULL,
			date TEXT NOT NULL,
			FOREIGN KEY (client_id) REFERENCES client(id)
		)`,
		`CREATE TABLE IF NOT EXISTS purchase (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sale_id INTEGER NOT NULL,
			service_id INTEGER NOT NULL,
			unit_price REAL NOT NULL,
			FOREIGN KEY (sale_id) REFERENCES sale(id),
			FOREIGN KEY (service_id) REFERENCES service(id)
		)`,
		`CREATE TABLE IF NOT EXISTS admin (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL
		)`,
	}

	for _, stmt := range statements {
		if _, err := database.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}

// VerifyAdmin checks username/password against the admin table.
// Returns the admin ID if valid, 0 if invalid, or error on DB issues.
func VerifyAdmin(database *sql.DB, username, password string) (int, error) {
	var admin Admin
	err := database.QueryRow(
		"SELECT id, password_hash FROM admin WHERE username = ?",
		username,
	).Scan(&admin.ID, &admin.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // not found
		}
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return 0, nil // password mismatch
	}

	return admin.ID, nil
}

// CreateAdmin inserts a new admin with hashed password.
// Returns the new admin ID.
func CreateAdmin(database *sql.DB, username, password string) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	result, err := database.Exec(
		"INSERT INTO admin (username, password_hash) VALUES (?, ?)",
		username, string(hash),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
