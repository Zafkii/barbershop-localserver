package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"barbershop-localserver/db"
)

func registerAuthRoutes(router *http.ServeMux, database *sql.DB) {
	router.HandleFunc("POST /auth/login", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if req.Username == "" || req.Password == "" {
			writeError(w, http.StatusBadRequest, "username and password are required")
			return
		}

		adminID, err := db.VerifyAdmin(database, req.Username, req.Password)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if adminID == 0 {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"admin_id": adminID,
			"username": req.Username,
			"message":  "login successful",
		})
	})
}