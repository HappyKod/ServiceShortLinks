package models

import "time"

// Link сокращенная ссылка
type Link struct {
	ShortKey string
	FullURL  string
	UserID   string
	Created  time.Time
	Del      bool
}
