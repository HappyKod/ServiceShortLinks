package pglinkssotorage

import (
	"HappyKod/ServiceShortLinks/utils"
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
	_, err := connect.Exec("CREATE TABLE if not exists public.urls (\n id text,\n long_url text primary key,\n created timestamp default now());")
	return err
}

func (PGS PGLinksStorage) Put(key string, url string) error {

	_, err := PGS.connect.Query("INSERT INTO public.urls (id, long_url) values ($1, $2);", key, url)
	return err
}
func (PGS PGLinksStorage) Get(key string) (string, error) {
	var longURL string
	rows, err := PGS.connect.Query("SELECT long_url from public.urls where id = $1", key)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		if err = rows.Scan(&longURL); err != nil {
			return "", err
		}
	}
	return longURL, nil
}

func (PGS PGLinksStorage) CreateUniqKey() (string, error) {
	var key string
	for {
		key = utils.GeneratorStringUUID()
		url, err := PGS.Get(key)
		if err != nil {
			return "", err
		}
		if url == "" {
			break
		}
	}
	return key, nil
}
