package reservation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/novando/go-cinema/internal/reservation/controller"
	"github.com/novando/go-cinema/internal/reservation/repository"
	"github.com/novando/go-cinema/internal/reservation/service"
	"github.com/novando/go-cinema/pkg/db/pg"
	"github.com/zishang520/socket.io/v2/socket"
)

func Init(r fiber.Router, soc *socket.Server, db *pg.PG) {
	rr := repository.NewReservation(db)

	rs := service.NewReservation(rr, soc)

	rc := controller.NewReservation(rs)
	r.Get("/now-playing", rc.GetNowPlaying)
	r.Get("/generate", rc.GenerateScreen)
	r.Get("/screen", rc.GetScreens)
	r.Get("/order", rc.GetOrders)
	r.Post("/order", rc.Book)
	r.Post("/login", rc.Login)
}
