package domain

import "time"

type Budget struct {
	ID        int
	Name      string
	CreatedAt time.Time
}
