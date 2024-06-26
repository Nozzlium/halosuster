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
		Prefork:     false,
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

	userRepo := repository.NewUserRepository(
		db,
	)
	patientRepo := repository.NewPatientRepository(
		db,
	)
	recordRepo := repository.NewRecordRepository(
		db,
	)

	userService := service.NewUserService(
		userRepo,
		int(cfg.BCryptSalt),
		cfg.JWTSecret,
	)
	patientService := service.NewPatientService(
		patientRepo,
	)
	recordService := service.NewRecordService(
		recordRepo,
	)

	userHandler := handler.NewUserHandler(
		userService,
	)
	patientHandler := handler.NewPatientHandler(
		patientService,
	)
	recordHandler := handler.NewRecordHandler(
		recordService,
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
		userHandler.LoginNurse,
	)
	userNurseProtected := userNurse.
		Use(middleware.Protected()).
		Use(middleware.SetClaimsData())
	userNurseProtected.Post(
		"/register",
		userHandler.RegisterNurse,
	)
	userNurseProtected.Put(
		"/:userId",
		userHandler.Update,
	)
	userNurseProtected.Delete(
		"/:userId",
		userHandler.Delete,
	)
	userNurseProtected.Post(
		"/:userId/access",
		userHandler.GrantNurseAccess,
	)

	user := v1.Group("/user")
	user.Use(middleware.Protected())
	user.Get("", userHandler.FindAll)

	patient := v1.Group(
		"/medical/patient",
	)
	patient.Use(middleware.Protected()).
		Use(middleware.SetClaimsData())
	patient.Post(
		"",
		patientHandler.Create,
	)
	patient.Get(
		"",
		patientHandler.FindAll,
	)

	record := v1.Group(
		"/medical/record",
	)
	record.Use(middleware.Protected()).
		Use(middleware.SetClaimsData())
	record.Post(
		"",
		recordHandler.Create,
	)

	return nil
}
