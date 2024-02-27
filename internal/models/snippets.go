package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

// Insert a new snippet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	sqlQuery := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(sqlQuery, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get a specific snippet by ID
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	sqlQuery := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(sqlQuery, id)
	// Initialize a pointer to an empty instance of a Snippet
	snippet := &Snippet{}
	// Copy the values from the returned row (if one) to the struct
	err := row.Scan(
		&snippet.ID,
		&snippet.Title,
		&snippet.Content,
		&snippet.Created,
		&snippet.Expires,
	)

	if err != nil {
		// If the DB driver returned no rows
		if errors.Is(err, sql.ErrNoRows) {
			// We'll return a custom error
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return snippet, nil
}

// Return the 10 most recent snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
