package controller

import (
	"github/codingMaster/internal/database"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Application struct {
	Store   database.Store
	Address string
}

func (app *Application) Moul() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/leabgin", app.tst)
	// project
	r.Post("/project/add", app.AddProject)
	r.Delete("/project/delete", app.DeleteProject)
	r.Put("/project/update", app.UpdateProject)
	r.Post("/project/upprove", app.Upprove)
	r.Post("/project/unpprove", app.Unpprove)
	r.Get("/project/get", app.GetProject)

	r.Post("/auth/login", app.Login)
	r.Post("/auth/signup", app.Singup)

	r.Post("/upload/pdf", app.UploadPDF)
	r.Get("/download/pdf", app.DownloadPDF)

	r.Post("/convert/json", app.ConvertJSON)

	r.Get("/download/csv", app.DownloadCSV)
	r.Get("/download/excel", app.DownloadExcel)
	return r
}

func (app *Application) Run(mux *chi.Mux) error {

	srv := &http.Server{
		Addr:         app.Address,
		Handler:      mux,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Server is live in Port %s", app.Address)
		return err
	}
	return nil
}
