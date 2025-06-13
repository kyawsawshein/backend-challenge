package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"backend-challenge/config"
	"backend-challenge/auth"
	"backend-challenge/service"

	"github.com/chasefleming/elem-go"
	"github.com/chasefleming/elem-go/attrs"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const realTimeOut = 5 * time.Second

var contextKey = "user"

var apiApp  *fiber.App
var apiOnce	sync.Once

type Api struct {
	apiApp  *fiber.App
	Ctx 	context.Context
}

var api *Api


func jwtAuth(config *config.ApiConfig) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.JwtSecret)},
		Filter: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/user/")
		},
		ContextKey: contextKey,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("Invalid, missing or expired JWT, %s", err))
		},
	})
}

func GetWebApi(Ctx context.Context) *Api {
	apiOnce.Do(func() {api = createFiberApp(Ctx)})
	return api
}

func createFiberApp(Ctx context.Context) *Api {
	api_config := config.Cfg.API
	app := fiber.New(fiber.Config{
		ReadTimeout: realTimeOut,
	})
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("config", api_config)
		c.Locals("context", Ctx)
		return c.Next()
	})
	
	app.Get("/", defaultHandler)
	app.Post("/login", login)
	app.Use(jwtAuth(&api_config))
	return &Api{
		apiApp: app,
	}
}

func defaultHandler(c *fiber.Ctx) error {
	head := elem.Head(attrs.Props{})
	body := elem.Body(
		attrs.Props{},
		elem.H1(attrs.Props{}, elem.Text("Hello World!")),
	)
	pageContent := elem.Html(attrs.Props{}, head, body)
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.Status(http.StatusOK).SendString(pageContent.Render())
}

func login(c *fiber.Ctx) error {
	var payload struct {
		UserEmail string `json:"email"`
		Password  string `json:"password"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid input", "detail": err})
	}
	user, err := service.GetUser(c, payload.UserEmail)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "User not found.", "detail": err})
	}
	if ok := auth.DoPasswordsMatch(user.Password, payload.Password); !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	api_config := config.Cfg.API
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(api_config.JwtSecret))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Token error!", "detail": err})
	}
	c.Locals("token", signedToken)
	return c.JSON(fiber.Map{"token": signedToken})
}

func (a Api) GetApp() *fiber.App {
	return a.apiApp
}

func (a *Api) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		go func() {
			<-ctx.Done()
			if err := a.apiApp.Shutdown(); err != nil {
				log.Fatal("Unable to shutdown API server", err)
			}
		}()
		if err := a.apiApp.Listen(":" + config.Cfg.API.Port); err != nil {
			log.Fatal("Unable to start Web Api", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}
