package main

import (
	"log"

	"github.com/bytedance/sonic"
	"github.com/caarlos0/env/v11"
	"github.com/gofiber/fiber/v2"
	"github.com/nozzlium/halosuster/internal/client"
	"github.com/nozzlium/halosuster/internal/config"
	"github.com/nozzlium/halosuster/internal/handler"
	"github.com/nozzlium/halosuster/internal/middleware"
	"github.com/nozzlium/halosuster/internal/repository"
	"github.com/nozzlium/halosuster/internal/service"
)

func main() {
	fiberApp := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		Prefork:     true,
	})

	err := setupApp(fiberApp)
	if err != nil {
		log.Fatal(err)
	}

	err = fiberApp.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func setupApp(app *fiber.App) error {
	var cfg config.Config
	opts := env.Options{
		TagName: "json",
	}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		log.Fatalf("%+v\n", err)
		return err
	}

	db, err := client.InitDB(cfg.DB)
	if err != nil {
		log.Fatal(err)
		return err
	}

	authRepo := repository.NewAuthRepository(
		db,
	)

	userService := service.NewUserService(
		authRepo,
		int(cfg.BCryptSalt),
		cfg.JWTSecret,
	)
	nurseService := service.NewNurseService(
		authRepo,
	)

	userHandler := handler.NewUserHandler(
		userService,
	)
	nurseHandler := handler.NewNurseHandler(
		nurseService,
	)

	v1 := app.Group("/v1")

	userIt := v1.Group("/user/it")
	userIt.Post(
		"/register",
		userHandler.Register,
	)
	userIt.Post(
		"/login",
		userHandler.Login,
	)

	userNurse := v1.Group("/user/nurse")
	userNurse.Post(
		"/login",
		nurseHandler.Login,
	)
	userNurseProtected := userNurse.
		Use(middleware.Protected()).
		Use(middleware.SetClaimsData())
	userNurseProtected.Post(
		"/register",
		nurseHandler.Register,
	)
	userNurseProtected.Post(
		"/:userId/access",
		nurseHandler.GiveAccess,
	)

	return nil
}
