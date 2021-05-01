package main

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/qkgo/yin"

	dbmanager "github.com/digimakergo/omlog/dbmanager"

	//for DB connection

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, _ := sql.Open("sqlite3", "./godb.db")

	r := chi.NewRouter()
	r.Use(yin.SimpleLogger)
	r.Get("/posts", func(w http.ResponseWriter, r *http.Request) {
		res, _ := yin.Event(w, r)
		items := dbmanager.GetAllLogFromDB(db)
		res.SendJSON(items)
	})

	r.Post("/posts", func(w http.ResponseWriter, r *http.Request) {
		res, req := yin.Event(w, r)
		body := map[string]string{}
		req.BindBody(&body)
		res.SendStatus(204)
	})

	r.Get("/getone", func(w http.ResponseWriter, r *http.Request) {
		res, _ := yin.Event(w, r)
		items := dbmanager.GetLogFromDB(db, 1)
		res.SendJSON(items)
	})

	r.Post("/getone", func(w http.ResponseWriter, r *http.Request) {
		res, req := yin.Event(w, r)
		body := map[string]string{}
		req.BindBody(&body)
		res.SendStatus(204)
	})

	http.ListenAndServe(":3000", r)

}
