package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"barbershop-localserver/db"
)

func registerServiceRoutes(router *http.ServeMux, database *sql.DB) {
	router.HandleFunc("GET /services", func(w http.ResponseWriter, r *http.Request) {
		services, err := db.ServiceList(database)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if services == nil {
			services = []db.Service{}
		}
		writeJSON(w, http.StatusOK, services)
	})

	router.HandleFunc("POST /services", func(w http.ResponseWriter, r *http.Request) {
		var service db.Service
		if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if service.Name == "" {
			writeError(w, http.StatusBadRequest, "name is required")
			return
		}
		if err := service.Create(database); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, service)
	})

	router.HandleFunc("GET /services/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(r)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		service, err := db.ServiceByID(database, id)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "service not found")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, service)
	})

	router.HandleFunc("PUT /services/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(r)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if _, err := db.ServiceByID(database, id); err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "service not found")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		var service db.Service
		if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		service.ID = id
		if err := service.Update(database); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, service)
	})

	router.HandleFunc("DELETE /services/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(r)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if _, err := db.ServiceByID(database, id); err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "service not found")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := db.ServiceDelete(database, id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	})
}
