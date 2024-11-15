package service

import "github.com/novando/go-cinema/internal/reservation/repository"

type Reservation struct {
	reservationRepo *repository.ReservationRepo
}

func NewReservation(rr *repository.ReservationRepo) *Reservation {
	return &Reservation{rr}
}

func (s *Reservation) GetNowPlaying() {

}
