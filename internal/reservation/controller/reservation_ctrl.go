package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/novando/go-cinema/internal/reservation/service"
	"github.com/novando/go-cinema/pkg/common/dto"
	"github.com/novando/go-cinema/pkg/common/value"
	"github.com/novando/go-cinema/pkg/logger"
	"golang.org/x/crypto/bcrypt"
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

func (c *Reservation) GetScreens(ctx *fiber.Ctx) error {
	var q service.ScreenRequestDTO
	if err := ctx.QueryParser(&q); err != nil {
		logger.Call().Errorf(err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.StdResponse{Message: "PARAM_INVALID"})
	}
	serv := c.reservationServ.GetScreenList(q)
	res := dto.StdResponse{Message: "SCREEN_FETCHED", Data: serv.Data}
	if serv.Err != nil {
		res.Message = serv.Err.Error()
	}
	return ctx.Status(serv.Code).JSON(res)
}

func (c *Reservation) Book(ctx *fiber.Ctx) error {
	var p service.BookRequestDTO
	if err := ctx.BodyParser(&p); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.StdResponse{Message: "PARAM_INVALID"})
	}
	serv := c.reservationServ.Book(p)
	res := dto.StdResponse{Message: "BOOKED", Data: serv.Data}
	if serv.Err != nil {
		res.Message = serv.Err.Error()
	}
	return ctx.Status(serv.Code).JSON(res)
}

func (c *Reservation) GetOrders(ctx *fiber.Ctx) error {
	apiKey := ctx.Get("x-api-key")
	if apiKey == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.StdResponse{Message: "ACCESS_DENIED"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(apiKey), []byte(value.ADMIN_PASS)); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.StdResponse{Message: "ACCESS_DENIED"})
	}
	serv := c.reservationServ.GetOrders()
	res := dto.StdResponse{Message: "ORDER_FETCH", Data: serv.Data}
	if serv.Err != nil {
		res.Message = serv.Err.Error()
	}
	return ctx.Status(serv.Code).JSON(res)
}

// Temporary in reservation, supposed to be in auth/user/iam

type loginDTO struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

func (c *Reservation) Login(ctx *fiber.Ctx) error {
	var p loginDTO
	if err := ctx.BodyParser(&p); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.StdResponse{Message: "PARAM_INVALID"})
	}
	if p.User != value.ADMIN_USER || p.Pass != value.ADMIN_PASS {
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.StdResponse{Message: "ACCESS_DENIED"})
	}
	res, err := bcrypt.GenerateFromPassword([]byte(p.Pass), 10)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.StdResponse{Message: "PARAM_INVALID"})

	}
	return ctx.Status(fiber.StatusOK).JSON(dto.StdResponse{Message: "LOGGED_IN", Data: string(res)})
}
