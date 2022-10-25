package pglinkssotorage

import (
	"HappyKod/ServiceShortLinks/internal/models"
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
	return createTable(PGS.connect)
}

func createTable(connect *sql.DB) error {
	_, err := connect.Exec("CREATE TABLE if not exists public.urls (id text,\nlong_url text primary key,\nuser_id text,\ncreated timestamp default now());")
	return err
}

func (PGS PGLinksStorage) PutShortLink(key string, link models.Link) error {
	_, err := PGS.connect.Exec("INSERT INTO public.urls (id, long_url, user_id) values ($1, $2, $3);", key, link.FullURL, link.UserID)
	return err
}
func (PGS PGLinksStorage) GetShortLink(key string) (string, error) {
	var longURL string
	row := PGS.connect.QueryRow("SELECT long_url from public.urls where id = $1", key)
	err := row.Scan(&longURL)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return longURL, row.Err()
}

// ManyPutShortLink добавляем множества значений
func (PGS PGLinksStorage) ManyPutShortLink(links []models.Link) error {
	scope, err := PGS.connect.Begin()
	if err != nil {
		return err
	}
	batch, err := scope.Prepare("INSERT INTO public.urls (id, long_url, user_id) values ($1, $2, $3)")
	if err != nil {
		return err
	}
	for _, link := range links {
		_, err = batch.Exec(link.ShortKey, link.FullURL, link.UserID)
		if err != nil {
			return err
		}
	}
	return scope.Commit()
}

func (PGS PGLinksStorage) GetKey(fullURL string) (string, error) {
	var key string
	err := PGS.connect.QueryRow("SELECT id FROM public.urls where long_url = $1", fullURL).Scan(&key)
	if err != nil {
		return "", err
	}
	return key, nil
}

func (PGS PGLinksStorage) GetShortLinkUser(UserID string) ([]models.Link, error) {
	var links []models.Link
	rows, err := PGS.connect.Query("SELECT id, long_url from public.urls where user_id = $1", UserID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var link models.Link
		if err = rows.Scan(&link.ShortKey, &link.FullURL); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, rows.Err()
}
