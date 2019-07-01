package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Database interface
type Database interface {
	Get(short string) (string, error)
	GetID(longURL string) (int64, string, error)
	Save(url string) (int64, error)
	saveShort(url string, shortURL string) (int64, error)
	shortExists(shortURL string) (bool, error)
	getLastID() (int64, error)
	setLastID(id int64) (int64, error)
}

type sqlite struct {
	Path string
}

func (s sqlite) getLastID() (int64, error) {
	db, err := sql.Open("sqlite3", s.Path)
	stmt, err := db.Prepare("SELECT value from variables WHERE name = 'short'")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var id int64
	stmt.QueryRow().Scan(&id)

	return id, nil
}

func (s sqlite) setLastID(id int64) (int64, error) {
	db, err := sql.Open("sqlite3", s.Path)
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare("REPLACE INTO variables (name, value) VALUES ('short', ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}

	newid, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	tx.Commit()

	return newid, nil
}

func (s sqlite) shortExists(shortURL string) (bool, error) {
	db, err := sql.Open("sqlite3", s.Path)
	stmt, err := db.Prepare("select 1 from urls where short = ?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var url int64
	stmt.QueryRow(shortURL).Scan(&url)

	if url == 1 {
		return true, nil
	}

	return false, nil
}

func (s sqlite) Save(url string) (int64, error) {
	db, err := sql.Open("sqlite3", s.Path)
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare("insert into urls(url) values(?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(url)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	tx.Commit()

	return id, nil
}

func (s sqlite) saveShort(url string, shortURL string) (int64, error) {
	db, err := sql.Open("sqlite3", s.Path)
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare("replace into urls(url, short) values(?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(url, shortURL)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	tx.Commit()

	return id, nil
}

func (s sqlite) GetID(longURL string) (int64, string, error) {
	db, err := sql.Open("sqlite3", s.Path)
	stmt, err := db.Prepare("select id, short from urls where url = ?")
	if err != nil {
		return 0, "0", err
	}
	defer stmt.Close()

	var id int64
	var short string
	stmt.QueryRow(longURL).Scan(&id, &short)

	return id, short, nil
}

func (s sqlite) Get(short string) (string, error) {
	db, err := sql.Open("sqlite3", s.Path)
	stmt, err := db.Prepare("select url from urls where short = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var url string
	err = stmt.QueryRow(short).Scan(&url)

	if err != nil {
		return "", err
	}

	return url, nil
}

func (s sqlite) Init() {
	c, err := sql.Open("sqlite3", s.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	sqlStmt := `CREATE TABLE if not exists urls (id integer NOT NULL PRIMARY KEY, short varchar (255) NOT NULL UNIQUE, url text NOT NULL);
		CREATE INDEX if not exists "short" ON urls (short);
		CREATE INDEX if not exists "url" ON urls (url);`
	_, err = c.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt = `CREATE TABLE if not exists variables (name varchar (255) NOT NULL primary key, value text NOT NULL);`
	_, err = c.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}
