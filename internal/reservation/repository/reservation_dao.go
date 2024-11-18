package repository

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type MovieInfoDAO struct {
	Duration int16       `json:"duration"`
	Title    string      `json:"title"`
	Synopsis string      `json:"synopsis"`
	Poster   string      `json:"poster"`
	ID       pgtype.UUID `json:"id"`
}

type BookParams struct {
	Price    float64
	Name     string
	Seats    []string
	ScreenID pgtype.UUID
}

type ScreenSimpleDAO struct {
	Rows     int16    `json:"rows"`
	Config   []int16  `json:"config"`
	Disabled []string `json:"disabled"`
	Occupied []string `json:"occupied"`
}

type ScreenDAO struct {
	Rows      int16       `json:"rows"`
	Price     float64     `json:"price"`
	Config    []int16     `json:"config"`
	Disabled  []string    `json:"disabled"`
	Occupied  []string    `json:"occupied"`
	ID        pgtype.UUID `json:"id"`
	StartedAt time.Time   `json:"startedAt"`
}

type OrderDAO struct {
	Price float64
	Start time.Time
	ID    pgtype.UUID
}
