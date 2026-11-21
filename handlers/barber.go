package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"barbershop-localserver/db"
)

// NewRouter creates the HTTP mux with all barber shop API routes.
func NewRouter(database *sql.DB) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /health", healthHandler)

	registerAuthRoutes(router, database)
	registerBarberRoutes(router, database)
	registerClientRoutes(router, database)
	registerServiceRoutes(router, database)
	registerSaleRoutes(router, database)
	registerPurchaseRoutes(router, database)

	return router
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, code int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, code int, message string) {
	writeJSON(w, code, map[string]string{"error": message})
}

func parseID(r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

func registerBarberRoutes(router *http.ServeMux, database *sql.DB) {
	router.HandleFunc("GET /barbers", func(w http.ResponseWriter, r *http.Request) {
		barbers, err := db.BarberList(database)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if barbers == nil {
			barbers = []db.Barber{}
		}
		writeJSON(w, http.StatusOK, barbers)
	})

	router.HandleFunc("POST /barbers", func(w http.ResponseWriter, r *http.Request) {
		var barber db.Barber
		if err := json.NewDecoder(r.Body).Decode(&barber); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if barber.Name == "" {
			writeError(w, http.StatusBadRequest, "name is required")
			return
		}
		if err := barber.Create(database); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, barber)
	})

	router.HandleFunc("GET /barbers/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(r)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		barber, err := db.BarberByID(database, id)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "barber not found")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, barber)
	})

	router.HandleFunc("PUT /barbers/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(r)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if _, err := db.BarberByID(database, id); err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "barber not found")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		var barber db.Barber
		if err := json.NewDecoder(r.Body).Decode(&barber); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		barber.ID = id
		if err := barber.Update(database); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, barber)
	})

	router.HandleFunc("DELETE /barbers/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(r)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if _, err := db.BarberByID(database, id); err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "barber not found")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := db.BarberDelete(database, id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	})
}
