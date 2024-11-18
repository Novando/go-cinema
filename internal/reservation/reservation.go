package reservation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/novando/go-cinema/internal/reservation/controller"
	"github.com/novando/go-cinema/internal/reservation/repository"
	"github.com/novando/go-cinema/internal/reservation/service"
	"github.com/novando/go-cinema/pkg/db/pg"
)

func Init(r fiber.Router, db *pg.PG) {
	rr := repository.NewReservation(db)

	rs := service.NewReservation(rr)

	rc := controller.NewReservation(rs)
	r.Get("/now-playing", rc.GetNowPlaying)
	r.Get("/generate", rc.GenerateScreen)
	r.Get("/screen", rc.GetScreens)
}
