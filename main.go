package main

import (
	"fmt"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v2"
	"time"
)

const jwtSecret = "asecret"

//func authRequired() func (ctx *fiber.Ctx) error{
//	return jwtware.New(jwtware.Config{
//		SigningKey:     []byte(jwtSecret),
//	})
//}

func Init(ctx *fiber.Ctx) error {
	return ctx.Send([]byte("Hello world"))
}

func Hello(ctx *fiber.Ctx) error {
	return ctx.Send([]byte("Hello"))
}
func HelloId(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["sub"].(string)
	return ctx.Send([]byte(fmt.Sprintf("Hello user with id: %s", id)))
}

func Login(ctx *fiber.Ctx) error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body request
	err := ctx.BodyParser(&body)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})

	}
	if body.Email != "julio@mail.com" || body.Password != "password123" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Credentials",
		})
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "1"
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7) // a week

	s, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": s,
		"user": struct {
			Id    int    `json:"id"`
			Email string `json:"email"`
		}{
			Id:    1,
			Email: "julio@mail.com",
		},
	})

}

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Get("/", Init)
	app.Post("/login", Login)

	app.Use(jwtware.New(jwtware.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
		SigningKey: []byte(jwtSecret),
	}))

	app.Get("/hello", Hello)
	app.Get("/hello/id", HelloId)
	//app.Get("/hello", authRequired(), func(ctx *fiber.Ctx) error {
	//	return ctx.Send([]byte("Hello"))
	//
	//})

	app.Listen(":3000")

}
