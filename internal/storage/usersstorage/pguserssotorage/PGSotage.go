package pguserssotorage

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type PGLinksStorage struct {
	connect *sql.DB
}

func New(url string) (*PGLinksStorage, error) {
	connect, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PGLinksStorage{connect: connect}, nil
}

func (PGS PGLinksStorage) Ping() error {
	err := PGS.connect.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (PGS PGLinksStorage) Put(key string, link string) error {
	return nil
}
func (PGS PGLinksStorage) Get(key string) ([]string, error) {
	return nil, nil
}
