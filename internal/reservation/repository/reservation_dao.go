package repository

import "github.com/jackc/pgx/v5/pgtype"

type MovieInfo struct {
	Duration int16       `json:"duration"`
	Title    string      `json:"title"`
	Synopsis string      `json:"synopsis"`
	Poster   string      `json:"poster"`
	ID       pgtype.UUID `json:"id"`
}
