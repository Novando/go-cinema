package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/novando/go-cinema/pkg/db/pg"
	"math"
	"math/rand"
	"slices"
	"strings"
	"time"
)

type Reservation struct {
	db *pg.PG
}

func NewReservation(db *pg.PG) *Reservation {
	return &Reservation{db}
}

type movieDuration struct {
	duration uint16
	id       pgtype.UUID
}

func get7DaysMovies(ctx context.Context, tx pg.PGTX) (res []movieDuration, err error) {
	rows, err := tx.Query(ctx, `-- ReservationGet7DaysMovie
		SELECT id, duration FROM movies
		WHERE release_date < (NOW() + '1 day'::INTERVAL)
			AND deleted_at IS NULL
			AND (taken_off_date IS NULL OR taken_off_date > (NOW() + '1 week'::INTERVAL))
	`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var i movieDuration
		if err = rows.Scan(
			&i.id,
			&i.duration,
		); err != nil {
			return
		}
		res = append(res, i)
	}
	err = rows.Err()
	return
}

func getCinema(ctx context.Context, tx pg.PGTX) (res []pgtype.UUID, err error) {
	rows, err := tx.Query(ctx, `-- ReservationGetCinema
		SELECT id FROM cinemas WHERE deleted_at IS NULL
	`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var i pgtype.UUID
		if err = rows.Scan(&i); err != nil {
			return
		}
		res = append(res, i)
	}
	err = rows.Err()
	return
}

func checkLastTimetable(ctx context.Context, tx pg.PGTX) (res time.Time, err error) {
	row := tx.QueryRow(ctx, `-- ReservationCheckLastTimetable
		SELECT finished_at FROM screens WHERE deleted_at IS NULL AND finished_at > (NOW() - '1 day'::INTERVAL)
		ORDER BY finished_at DESC LIMIT 1
	`)
	if err = row.Scan(&res); err != nil && err.Error() == pgx.ErrNoRows.Error() {
		err = nil
	}
	return
}

type timetable struct {
	cinema pgtype.UUID
	movie  pgtype.UUID
	start  time.Time
	end    time.Time
}

func generateTimetable(m []movieDuration, cinemas []pgtype.UUID, lastTimetable time.Time) (res []timetable) {
	// Get filled days for the next 7 days
	durF := (time.Until(lastTimetable).Hours() + 23) / 24
	dur := int(math.Ceil(durF))
	if dur < 0 {
		dur = 0
	}

	// Set theater close and open time
	oh := 9
	ch := 21
	begin := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), oh+dur, 0, 0, 0, time.Local)
	end := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), ch+dur, 0, 0, 0, time.Local)

	movies := make([]movieDuration, len(m))
	copy(movies, m)
	for dur < 7 {
		begin = begin.Add(24 * time.Hour)
		end = end.Add(24 * time.Hour)
		for _, c := range cinemas {
			newBegin := begin
			for end.Sub(newBegin) >= time.Duration(0) {
				if len(movies) == 0 {
					movies = make([]movieDuration, len(m))
					copy(movies, m)
				}
				rs := rand.New(rand.NewSource(time.Now().UnixNano()))
				idx := rs.Intn(len(movies))
				movie := movies[idx]
				md := movie.duration / 15
				if movie.duration%15 != 0 {
					md++
				}
				md = (md * 15) + 30
				me := newBegin.Add(time.Minute * time.Duration(md))
				res = append(res, timetable{
					cinema: c,
					movie:  movie.id,
					start:  newBegin,
					end:    me,
				})
				newBegin = me
				movies = append(movies[:idx], movies[idx+1:]...)
			}
		}
		dur++
	}
	return
}

func insertTimetables(ctx context.Context, tx pg.PGTX, t []timetable) error {
	if len(t) < 1 {
		return nil
	}
	var values []string
	for _, v := range t {
		values = append(values, fmt.Sprintf(
			"\n('%x'::UUID, '%x'::UUID, '%v'::TIMESTAMPTZ, '%v'::TIMESTAMPTZ)",
			v.movie.Bytes,
			v.cinema.Bytes,
			v.end.Format(time.RFC3339),
			v.start.Format(time.RFC3339),
		))
	}
	sql := fmt.Sprintf(`-- ReservationInsertTimetables
		INSERT INTO screens (movie_id, cinema_id, started_at, finished_at) VALUES %v
	`, strings.Join(values, ","))
	_, err := tx.Exec(ctx, sql)
	return err
}

func (r *Reservation) GetNowPlaying() (res []MovieInfoDAO, err error) {
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
		var i MovieInfoDAO
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

func (r *Reservation) CountNowPlaying() (res uint64, err error) {
	row := r.db.QueryRow(`-- ReservationCountNowPlaying
		SELECT COUNT(id) FROM movies
		WHERE taken_off_date IS NULL AND release_date < NOW() AND deleted_at IS NULL
	`)
	if err = row.Scan(&res); err != nil && err.Error() == pgx.ErrNoRows.Error() {
		err = nil
	}
	return
}

func (r *Reservation) GenerateScreen() error {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	movies, err := get7DaysMovies(ctx, tx)
	if err != nil {
		return err
	}
	cinemas, err := getCinema(ctx, tx)
	if err != nil {
		return err
	}
	lastFinish, err := checkLastTimetable(ctx, tx)
	if err != nil {
		return err
	}
	newTimetable := generateTimetable(movies, cinemas, lastFinish)
	return insertTimetables(ctx, tx, newTimetable)
}

func (r *Reservation) Book(arg BookParams) (res OrderSimpleDAO, err error) {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	rowOccupancies := tx.QueryRow(ctx, `-- ReservationBookGetFilled
		SELECT occupancy, started_at FROM screens WHERE id = $1::UUID AND deleted_at IS NULL
	`, arg.ScreenID)
	var occupancies []string
	if err = rowOccupancies.Scan(&occupancies, &res.Start); err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	for _, seat := range arg.Seats {
		if slices.Contains(occupancies, seat) {
			err = errors.New("SEAT_OCCUPIED")
			return
		}
	}
	rowOrder := tx.QueryRow(ctx, `-- ReservationBookSeat
		INSERT INTO orders (screen_id, name, seats, price) VALUES ($1::UUID, $2::TEXT, $3, $4::FLOAT) RETURNING id
	`, arg.ScreenID, arg.Name, arg.Seats, arg.Price)
	if err = rowOrder.Scan(&res.ID); err != nil {
		return
	}
	_, err = tx.Exec(ctx, `-- ReservationBookUpdateOccupancies
		UPDATE screens SET occupancy = $2 WHERE id = $1::UUID
	`, arg.ScreenID, append(occupancies, arg.Seats...))
	res.Price = arg.Price
	return
}

func (r *Reservation) GetScreen(id pgtype.UUID) (res ScreenSimpleDAO, err error) {
	row := r.db.QueryRow(`-- ReservationGetScreen
		SELECT COALESCE(occupancy, '{}'), column_group, number_of_row, COALESCE(disabled_seat, '{}') FROM screens s
		LEFT JOIN cinemas c ON c.id = s.cinema_id
		WHERE s.id = $1::UUID AND s.deleted_at IS NULL
	`, id)
	if err = row.Scan(
		&res.Occupied,
		&res.Config,
		&res.Rows,
		&res.Disabled,
	); err != nil && err.Error() == pgx.ErrNoRows.Error() {
		err = nil
	}
	return
}

func (r *Reservation) GetScreenList(mid pgtype.UUID, start time.Time) (res []ScreenDAO, err error) {
	rows, err := r.db.Query(`-- ReservationGetScreenList
		SELECT
			s.id,
			started_at::TIME,
			COALESCE(occupancy, '{}'),
			column_group,
			number_of_row,
			COALESCE(disabled_seat, '{}')
		FROM screens s
		LEFT JOIN cinemas c ON c.id = s.cinema_id
		WHERE s.movie_id = $1::UUID AND started_at::DATE = $2::DATE AND s.deleted_at IS NULL
	`, mid, start)

	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var i ScreenDAO
		if err = rows.Scan(
			&i.ID,
			&i.StartedAt,
			&i.Occupied,
			&i.Config,
			&i.Rows,
			&i.Disabled,
		); err != nil {
			return
		}
		res = append(res, i)
	}
	return
}

func (r *Reservation) UpdateOrderPrice(id pgtype.UUID, price float64) error {
	_, err := r.db.Exec(`-- ReservationUpdateOrderPrice
		UPDATE orders SET price = $2::FLOAT WHERE id = $1::UUID
	`, id, price)
	return err
}

func (r *Reservation) GetOrdered() (res []OrderDAO, err error) {
	rows, err := r.db.Query(`-- ReservationGetOrdered
		SELECT o.id, o.price, o.name, m.title, started_at, seats FROM orders o
		LEFT JOIN screens s ON s.id = o.screen_id
		LEFT JOIN movies m ON m.id = s.movie_id
		WHERE o.deleted_at IS NULL
		ORDER BY created_at DESC
	`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var i OrderDAO
		if err = rows.Scan(
			&i.ID,
			&i.Price,
			&i.OrderBy,
			&i.Title,
			&i.Start,
			&i.Seats,
		); err != nil {
			return
		}
		res = append(res, i)
	}
	err = rows.Err()
	return
}
