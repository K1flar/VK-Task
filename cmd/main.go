package main

import (
	_ "film_library/docs"
	"film_library/internal/config"
	"film_library/internal/handlers"
	"film_library/internal/logger"
	"film_library/internal/repositories/postgres"
	"film_library/internal/services"
	adminmw "film_library/pkg/middlewares/admin_mw"
	"film_library/pkg/middlewares/auth"
	loggermw "film_library/pkg/middlewares/logger_mw"
	"film_library/pkg/mux"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Swagger Film library
// @version 1.0
// @description This is a VK task.

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	log := logger.New()

	cfg, err := config.New("./configs/local.yaml")
	exitOnErr(log, err)

	repository, err := postgres.New(&cfg.Database)
	exitOnErr(log, err)

	service := services.New(repository, log, cfg)

	handler := handlers.New(service, log)

	router := mux.New()

	router.HandleFunc("GET /swagger/", httpSwagger.Handler())

	router.Use(loggermw.New(log))
	router.HandleFunc("POST /api/register", handler.Register)
	router.HandleFunc("POST /api/login", handler.Login)

	router.Group(func(r *mux.Mux) {
		r.Use(auth.New(log, cfg))

		r.HandleFunc("GET /api/actors", handler.GetActorsWithFilms)
		r.HandleFunc("GET /api/films", handler.GetFilms)

		r.Group(func(adminRouter *mux.Mux) {
			adminRouter.Use(adminmw.New(log))

			adminRouter.HandleFunc("POST /api/actor", handler.CreateActor)
			adminRouter.HandleFunc("POST /api/actors/{filmID}", handler.AddActorsToFilm)
			adminRouter.HandleFunc("PUT /api/actor/name/{id}/{name}", handler.UpdateActorFullName)
			adminRouter.HandleFunc("PUT /api/actor/gender/{id}/{gender}", handler.UpdateActorGender)
			adminRouter.HandleFunc("PUT /api/actor/birthday/{id}/{birthday}", handler.UpdateActorBirthday)
			adminRouter.HandleFunc("PUT /api/actor/{id}", handler.UpdateActor)
			adminRouter.HandleFunc("DELETE /api/actor/{id}", handler.DeleteActor)
			adminRouter.HandleFunc("DELETE /api/actor/{id}/{filmID}", handler.DeleteActorFromFilm)

			adminRouter.HandleFunc("POST /api/film", handler.CreateFilm)
			adminRouter.HandleFunc("PUT /api/film/name/{id}/{name}", handler.UpdateFilmName)
			adminRouter.HandleFunc("PUT /api/film/description/{id}", handler.UpdateFilmDescription)
			adminRouter.HandleFunc("PUT /api/film/date/{id}/{date}", handler.UpdateFilmReleaseDate)
			adminRouter.HandleFunc("PUT /api/film/{id}/{rating}", handler.UpdateFilmRating)
			adminRouter.HandleFunc("PUT /api/film/{id}", handler.UpdateFilm)
			adminRouter.HandleFunc("DELETE /api/film/{id}", handler.DeleteFilm)
		})
	})

	done := make(chan struct{})
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port), router)
		exitOnErr(log, err)

		done <- struct{}{}
	}()

	fmt.Println("Starting server...")
	<-done
}

func exitOnErr(log *slog.Logger, err error) {
	if err == nil {
		return
	}

	log.Error(err.Error())
	os.Exit(-1)
}
