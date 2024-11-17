package service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/novando/go-cinema/internal/reservation/repository"
	"github.com/novando/go-cinema/pkg/common/dto"
	"github.com/novando/go-cinema/pkg/helper"
	"github.com/novando/go-cinema/pkg/logger"
)

type Reservation struct {
	reservationRepo *repository.ReservationRepo
}

func NewReservation(rr *repository.ReservationRepo) *Reservation {
	return &Reservation{rr}
}

func (s *Reservation) GetNowPlaying() (res dto.StdService) {
	res.Code = fiber.StatusOK
	var rd []repository.MovieInfo
	rd, res.Err = s.reservationRepo.GetNowPlaying()
	if res.Err != nil {
		logger.Call().Errorf(res.Err.Error())
		res.Err = errors.New("ERR_GET_NOW_PLAYING")
		res.Code = fiber.StatusInternalServerError
		return
	}
	var c uint64
	c, res.Err = s.reservationRepo.CountNowPlaying()
	if res.Err != nil {
		logger.Call().Errorf(res.Err.Error())
		res.Err = errors.New("ERR_GET_NOW_PLAYING")
		res.Code = fiber.StatusInternalServerError
		return
	}
	res.Data = helper.CreateListResponse(c, rd)
	return
}

func (s *Reservation) GenerateScreen() (res dto.StdService) {
	res.Code = fiber.StatusOK
	res.Err = s.reservationRepo.GenerateScreen()
	return
}
