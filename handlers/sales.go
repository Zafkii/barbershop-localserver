package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"barbershop-localserver/db"
)

func registerSaleRoutes(router *http.ServeMux, database *sql.DB) {
	router.HandleFunc("GET /sales", func(w http.ResponseWriter, r *http.Request) {
		sales, err := db.SaleList(database)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if sales == nil {
			sales = []db.Sale{}
		}
		writeJSON(w, http.StatusOK, sales)
	})

	router.HandleFunc("POST /sales", func(w http.ResponseWriter, r *http.Request) {
		var sale db.Sale
		if err := json.NewDecoder(r.Body).Decode(&sale); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if sale.ClientID < 1 {
			writeError(w, http.StatusBadRequest, "client_id is required")
			return
		}
		if err := sale.Create(database); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, sale)
	})

	router.HandleFunc("GET /sales/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(r)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		sale, err := db.SaleByID(database, id)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "sale not found")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, sale)
	})

	router.HandleFunc("PUT /sales/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(r)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if _, err := db.SaleByID(database, id); err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "sale not found")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		var sale db.Sale
		if err := json.NewDecoder(r.Body).Decode(&sale); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		sale.ID = id
		if err := sale.Update(database); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, sale)
	})

	router.HandleFunc("DELETE /sales/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(r)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if _, err := db.SaleByID(database, id); err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "sale not found")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := db.SaleDelete(database, id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	})
}

func registerPurchaseRoutes(router *http.ServeMux, database *sql.DB) {
	router.HandleFunc("GET /purchases", func(w http.ResponseWriter, r *http.Request) {
		purchases, err := db.PurchaseList(database)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if purchases == nil {
			purchases = []db.Purchase{}
		}
		writeJSON(w, http.StatusOK, purchases)
	})

	router.HandleFunc("POST /purchases", func(w http.ResponseWriter, r *http.Request) {
		var purchase db.Purchase
		if err := json.NewDecoder(r.Body).Decode(&purchase); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if purchase.SaleID < 1 {
			writeError(w, http.StatusBadRequest, "sale_id is required")
			return
		}
		if purchase.ServiceID < 1 {
			writeError(w, http.StatusBadRequest, "service_id is required")
			return
		}
		if err := purchase.Create(database); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, purchase)
	})

	router.HandleFunc("GET /purchases/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(r)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		purchase, err := db.PurchaseByID(database, id)
		if err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "purchase not found")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, purchase)
	})

	router.HandleFunc("PUT /purchases/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(r)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if _, err := db.PurchaseByID(database, id); err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "purchase not found")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		var purchase db.Purchase
		if err := json.NewDecoder(r.Body).Decode(&purchase); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		purchase.ID = id
		if err := purchase.Update(database); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, purchase)
	})

	router.HandleFunc("DELETE /purchases/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(r)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		if _, err := db.PurchaseByID(database, id); err != nil {
			if err == sql.ErrNoRows {
				writeError(w, http.StatusNotFound, "purchase not found")
				return
			}
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := db.PurchaseDelete(database, id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	})
}
