package server

import "github.com/gofiber/fiber/v2"

type defaultRest struct{}

// defaultRestHandler will create an instace for default rest handler
func defaultRestHandler() *defaultRest {
	return &defaultRest{}
}

func (dr *defaultRest) Router(r fiber.Router) {}
