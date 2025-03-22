package main

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Storage struct {
	DB *sql.DB
}

func InitDB(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS links (
		handle TEXT PRIMARY KEY,
		target TEXT NOT NULL,
		owner TEXT NOT NULL,
		FOREIGN KEY(owner) REFERENCES users(username)
	);`
	if _, err := db.Exec(schema); err != nil {
		panic(err)
	}
	return db
}

func (s *Storage) CreateUser(username, password string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := s.DB.Exec("INSERT INTO users(username, password) VALUES (?, ?)", username, hash)
	return err
}

func (s *Storage) SaveLink(handle, target, owner string) error {
	_, err := s.DB.Exec(`
	INSERT INTO links(handle, target, owner)
	VALUES (?, ?, ?)
	ON CONFLICT(handle) DO UPDATE SET target = excluded.target WHERE owner = ?`,
		handle, target, owner, owner)
	return err
}

func (s *Storage) GetTarget(handle string) (string, error) {
	var target string
	err := s.DB.QueryRow("SELECT target FROM links WHERE handle = ?", handle).Scan(&target)
	return target, err
}

func (s *Storage) DeleteLink(handle, owner string) error {
	res, err := s.DB.Exec("DELETE FROM links WHERE handle = ? AND owner = ?", handle, owner)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("Not authorized or not found")
	}
	return nil
}
