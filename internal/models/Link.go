package models

import "time"

type Link struct {
	ShortKey string
	FullURL  string
	UserID   string
	Del      bool
	Created  time.Time
}
