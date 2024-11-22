package oui

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Oui interface {
	Resolve(mac string) ([]Vendor, string, error)
}

type oui struct {
	db *sql.DB
}

const (
	vendorSource = "https://standards-oui.ieee.org/oui/oui.txt"
)

func New(update bool) (Oui, func(), error) {
	statePath, err := statePath()
	if err != nil {
		return nil, nil, err
	}

	// Create the state path if it doesn't exist
	dbPath := filepath.Join(statePath, "vendors.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, nil, err
	}

	v := &oui{db: db}

	if err := v.ensureTables(); err != nil {
		return nil, nil, err
	}

	if err := v.ensureVendors(update); err != nil {
		return nil, nil, err
	}

	return v, func() { db.Close() }, nil
}

// Resolve resolves the vendor for a given MAC address.
func (v *oui) Resolve(mac string) ([]Vendor, string, error) {
	rows, err := v.db.Query("SELECT oui, org FROM vendors WHERE ? LIKE oui || '%'", normalize(mac))
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var vendors []Vendor
	for rows.Next() {
		var vendor Vendor
		if err := rows.Scan(&vendor.Oui, &vendor.Org); err != nil {
			return nil, "", err
		}
		vendors = append(vendors, vendor)
	}

	lastModified, err := v.lastModified()
	if err != nil {
		return nil, "", err
	}

	return vendors, lastModified, nil
}

// lastModified returns the last modified date of the vendors database.
func (v *oui) lastModified() (string, error) {
	var lastModified string
	if err := v.db.QueryRow("SELECT last_modified FROM metadata").Scan(&lastModified); err != nil {
		return "", err
	}
	return lastModified, nil
}

// normalize normalizes a MAC address. replacing all possible separators with nothing.
func normalize(mac string) string {
	separators := strings.NewReplacer("-", "", ":", "", ".", "")
	return separators.Replace(strings.ToUpper(mac))
}

// ensureTables ensures the required tables exist in the database.
func (v *oui) ensureTables() error {
	// Metadata table
	_, err := v.db.Exec(`
		CREATE TABLE IF NOT EXISTS metadata (
			id INTEGER PRIMARY KEY,
			last_modified TEXT,
			content_length INTEGER,
			etag TEXT,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Vendors table
	_, err = v.db.Exec(`
		CREATE TABLE IF NOT EXISTS vendors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			oui TEXT,
			org TEXT
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// statePath returns the path to the state directory.
// See: https://specifications.freedesktop.org/basedir-spec/latest/index.html#variables
func statePath() (string, error) {
	path := os.Getenv("XDG_STATE_HOME")
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(homeDir, ".local", "state")
	}
	return filepath.Join(path, "yeah"), nil
}
