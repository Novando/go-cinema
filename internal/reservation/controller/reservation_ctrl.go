package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/novando/go-cinema/internal/reservation/service"
	"github.com/novando/go-cinema/pkg/common/dto"
)

type Reservation struct {
	reservationServ *service.Reservation
}

func NewReservation(sr *service.Reservation) *Reservation {
	return &Reservation{sr}
}

func (c *Reservation) GetNowPlaying(ctx *fiber.Ctx) error {
	serv := c.reservationServ.GetNowPlaying()
	res := dto.StdResponse{Message: "NOW_PLAYING_FETCHED", Data: serv.Data}
	if serv.Err != nil {
		res.Message = serv.Err.Error()
	}
	return ctx.Status(serv.Code).JSON(res)
}

func (c *Reservation) GenerateScreen(ctx *fiber.Ctx) error {
	serv := c.reservationServ.GenerateScreen()
	res := dto.StdResponse{Message: "SCREEN_GENERATED", Data: serv.Data}
	if serv.Err != nil {
		res.Message = serv.Err.Error()
	}
	return ctx.Status(serv.Code).JSON(res)
}
