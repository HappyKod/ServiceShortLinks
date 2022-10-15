package models

type User struct {
	Id      string   `json:"id"`
	KeysUrl []string `json:"keys_url"`
}
