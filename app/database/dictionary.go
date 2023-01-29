package database

import "time"

type Dictionary struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Updated     time.Time `json:"updated"`
}
