package oui

import (
	"bufio"
	"database/sql"
	"net/http"
	"strings"
)

type (
	metadata struct {
		lastModified  string
		contentLength int64
		etag          string
	}

	Vendor struct {
		Oui string
		Org string
	}
)

// ensureVendors ensures that the vendors database populated and up to date.
func (v *oui) ensureVendors(update bool) error {
	isLatest, err := v.check()
	if err != nil {
		return err
	}

	if isLatest && !update {
		return nil
	}

	// Download the new vendor data
	if err := v.download(); err != nil {
		return err
	}

	return nil
}

// check checks if the vendors database is up to date.
func (v *oui) check() (bool, error) {
	// Get the current headers from the server
	resp, err := http.Head(vendorSource)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Get the latest headers from our database
	var stored metadata
	err = v.db.QueryRow(`
		SELECT last_modified, content_length, etag 
		FROM metadata 
		ORDER BY updated_at DESC
		LIMIT 1
	`).Scan(&stored.lastModified, &stored.contentLength, &stored.etag)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	// Compare headers to check if update is needed
	needsUpdate := err == sql.ErrNoRows || // No previous headers stored
		stored.lastModified != resp.Header.Get("Last-Modified") ||
		stored.contentLength != resp.ContentLength ||
		stored.etag != resp.Header.Get("ETag")

	if !needsUpdate {
		return true, nil // Database is up to date
	}

	// Store new metadata
	_, err = v.db.Exec(`
		INSERT INTO metadata (last_modified, content_length, etag)
		VALUES (?, ?, ?)
	`, resp.Header.Get("Last-Modified"), resp.ContentLength, resp.Header.Get("ETag"))
	if err != nil {
		return false, err
	}

	return false, nil
}

// download downloads the latest vendor data from the source and parses it.
func (v *oui) download() error {
	// Get the data from vendorSource
	resp, err := http.Get(vendorSource)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var vendors []Vendor

	for scanner.Scan() {
		line := scanner.Text()
		// Look for lines containing "(base 16)" as they contain the format we want
		if strings.Contains(line, "(base 16)") {
			parts := strings.Split(line, "(base 16)")
			if len(parts) != 2 {
				continue
			}

			oui := strings.TrimSpace(parts[0])
			org := strings.TrimSpace(parts[1])

			vendors = append(vendors, Vendor{
				Oui: oui,
				Org: org,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Begin transaction to update vendors
	tx, err := v.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear existing vendors
	_, err = tx.Exec("DELETE FROM vendors")
	if err != nil {
		return err
	}

	// Insert new vendors
	stmt, err := tx.Prepare("INSERT INTO vendors (oui, org) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, vendor := range vendors {
		_, err = stmt.Exec(vendor.Oui, vendor.Org)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
