package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"barbershop-localserver/db"
)

func registerAuthRoutes(router *http.ServeMux, database *sql.DB) {
	// POST /auth/login
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

		token, err := db.CreateSession(database, adminID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to create session")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"admin_id": adminID,
			"username": req.Username,
			"token":    token,
			"message":  "login successful",
		})
	})

	// GET /auth/verify - validates token and extends expiry
	router.HandleFunc("GET /auth/verify", func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r)
		if token == "" {
			writeError(w, http.StatusUnauthorized, "missing token")
			return
		}

		adminID, err := db.VerifySession(database, token)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if adminID == 0 {
			writeError(w, http.StatusUnauthorized, "invalid or expired session")
			return
		}

		// Extend session by 7 days
		db.RefreshSession(database, token)

		writeJSON(w, http.StatusOK, map[string]any{
			"admin_id": adminID,
			"valid":    true,
		})
	})

	// POST /auth/logout
	router.HandleFunc("POST /auth/logout", func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r)
		if token != "" {
			db.DeleteSession(database, token)
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"message": "logged out",
		})
	})
}

func extractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}