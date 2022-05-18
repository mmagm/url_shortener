package db

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

type Link struct {
	URL   string
	Token string
}

func generateToken() string {
	buff := make([]byte, 6)
	rand.Read(buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:6]
}

func (store *Store) RegisterURL(url string) (*Link, error) {
	token := generateToken()
	link := &Link{URL: url, Token: token}

	created := time.Now().UTC()

	// insert statement
	insertStatement := `INSERT INTO links (token, url, created_at, updated_at) VALUES ($1, $2, $3, $4)`

	_, err := store.db.Exec(insertStatement, link.Token, link.URL, created, created)

	if err != nil {
		return nil, err
	}

	return link, nil
}

func (store *Store) RetrieveURL(token string) (*Link, error) {
	row := store.db.QueryRow("SELECT token, url FROM links WHERE token = $1", token)

	link := &Link{Token: "", URL: ""}

	err := row.Scan(&link.Token, &link.URL)

	if err != nil {
		return nil, err
	}

	return link, err
}
