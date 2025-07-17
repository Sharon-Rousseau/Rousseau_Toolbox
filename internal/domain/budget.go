package domain

import "time"

type Budget struct {
	ID        string
	Name      string
	CreatedAt time.Time
}
