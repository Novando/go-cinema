package repository

import "github.com/jackc/pgx/v5/pgtype"

type MovieInfo struct {
	Duration int16
	Title    string
	Synopsis string
	Poster   string
	ID       pgtype.UUID
}
