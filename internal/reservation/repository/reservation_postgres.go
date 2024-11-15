package repository

import "github.com/novando/go-cinema/pkg/db/pg"

type ReservationRepo struct {
	db *pg.PG
}

func NewReservation(db *pg.PG) *ReservationRepo {
	return &ReservationRepo{db}
}

func (r *ReservationRepo) GetNowPlaying() (res []MovieInfo, err error) {
	rows, err := r.db.Query(`-- ReservationGetNowPlaying
		SELECT id, title, synopsis, poster, duration FROM movies
		WHERE taken_off_date IS NULL AND release_date < NOW() AND deleted_at IS NULL
		ORDER BY title
	`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var i MovieInfo
		if err = rows.Scan(
			&i.ID,
			&i.Title,
			&i.Synopsis,
			&i.Poster,
			&i.Duration,
		); err != nil {
			return
		}
		res = append(res, i)
	}
	err = rows.Err()
	return
}
