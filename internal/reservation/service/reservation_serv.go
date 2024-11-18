package service

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/novando/go-cinema/internal/reservation/repository"
	"github.com/novando/go-cinema/pkg/common/dto"
	"github.com/novando/go-cinema/pkg/helper"
	"github.com/novando/go-cinema/pkg/logger"
	"github.com/novando/go-cinema/pkg/uuid"
	"slices"
	"time"
)

type Reservation struct {
	reservationRepo *repository.Reservation
}

func NewReservation(rr *repository.Reservation) *Reservation {
	return &Reservation{rr}
}

func (s *Reservation) GetNowPlaying() (res dto.StdService) {
	res.Code = fiber.StatusOK
	var rd []repository.MovieInfoDAO
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

func (s *Reservation) GetScreen(idStr string) (res dto.StdService) {
	res.Code = fiber.StatusOK
	id, err := uuid.ParseUUID(idStr)
	if err != nil {
		res.Code = fiber.StatusBadRequest
		return
	}
	res.Data, res.Err = s.reservationRepo.GetScreen(id)
	return
}

func (s *Reservation) GetScreenList(arg ScreenRequestDTO) (res dto.StdService) {
	res.Code = fiber.StatusOK
	var id pgtype.UUID
	id, res.Err = uuid.ParseUUID(arg.MovieID)
	if res.Err != nil {
		res.Code = fiber.StatusBadRequest
		return
	}
	var date time.Time
	date, res.Err = time.Parse(time.DateOnly, arg.Date)
	if res.Err != nil {
		res.Code = fiber.StatusBadRequest
		return
	}
	var rd []repository.ScreenDAO
	rd, res.Err = s.reservationRepo.GetScreenList(id, date)
	if res.Err != nil {
		res.Code = fiber.StatusInternalServerError
		return
	}
	for i := range rd {
		rd[i].Price = 40000
		if slices.Contains([]time.Weekday{time.Sunday, time.Saturday}, rd[i].StartedAt.Weekday()) {
			rd[i].Price = 60000
		}
	}
	res.Data = helper.CreateListResponse(uint64(len(rd)), rd)
	return
}

func (s *Reservation) Book(arg BookRequestDTO) (res dto.StdService) {
	res.Code = fiber.StatusOK
	var sid pgtype.UUID
	sid, res.Err = uuid.ParseUUID(arg.ScreenID)
	if res.Err != nil {
		res.Code = fiber.StatusBadRequest
		return
	}

	var rd repository.OrderSimpleDAO
	rd, res.Err = s.reservationRepo.Book(repository.BookParams{
		ScreenID: sid,
		Seats:    arg.Seats,
		Name:     arg.Name,
		Price:    float64(40000 * len(arg.Seats)),
	})
	if res.Err != nil {
		logger.Call().Errorf(res.Err.Error())
		res.Code = fiber.StatusInternalServerError
		return
	}
	if !slices.Contains([]time.Weekday{time.Sunday, time.Saturday}, rd.Start.Weekday()) {
		return
	}
	res.Err = s.reservationRepo.UpdateOrderPrice(rd.ID, float64(60000*len(arg.Seats)))
	if res.Err != nil {
		logger.Call().Errorf(res.Err.Error())
		res.Code = fiber.StatusInternalServerError
		return
	}
	res.Data = fmt.Sprintf("%x", rd.ID.Bytes)
	return
}

func (s *Reservation) GetOrders() (res dto.StdService) {
	res.Code = fiber.StatusOK
	var rd []repository.OrderDAO
	rd, res.Err = s.reservationRepo.GetOrdered()
	if res.Err != nil {
		logger.Call().Errorf(res.Err.Error())
		res.Code = fiber.StatusInternalServerError
		return
	}
	res.Data = helper.CreateListResponse(uint64(len(rd)), rd)
	return
}
