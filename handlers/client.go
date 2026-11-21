package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"barbershop-localserver/db"
)

func registerClientRoutes(router *http.ServeMux, database *sql.DB) {
	router.HandleFunc("GET /clients", func(w http.ResponseWriter, r *http.Request) {
		clients, err := db.ClientList(database)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if clients == nil {
			clients = []db.Client{}
		}
		writeJSON(w, http.StatusOK, clients)
	})

	router.HandleFunc("POST /clients", func(w http.ResponseWriter, r *http.Request) {
		var client db.Client
		if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if client.Name == "" {
			writeError(w, http.StatusBadRequest, "name is required")
			return
		}
		if err := client.Create(database); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, client)
	})

	router.HandleFunc("GET /clients/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(r)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		client, err := db.ClientByID(database, id)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "client not found")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, client)
	})

	router.HandleFunc("PUT /clients/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(r)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if _, err := db.ClientByID(database, id); err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "client not found")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		var client db.Client
		if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		client.ID = id
		if err := client.Update(database); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, client)
	})

	router.HandleFunc("DELETE /clients/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(r)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if _, err := db.ClientByID(database, id); err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "client not found")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := db.ClientDelete(database, id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	})
}
