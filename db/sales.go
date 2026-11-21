package db

import (
	"database/sql"
)

// Create inserts a sale and stores the generated id in s.ID.
// Pointer receiver so the generated id propagates back to the caller.
func (s *Sale) Create(database *sql.DB) error {
	result, err := database.Exec(
		"INSERT INTO sale (client_id, date) VALUES (?, ?)",
		s.ClientID, s.Date,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	s.ID = int(id)
	return nil
}

// SaleByID returns the sale with the given id.
func SaleByID(database *sql.DB, id int) (Sale, error) {
	var sale Sale
	err := database.QueryRow(
		"SELECT id, client_id, date FROM sale WHERE id = ?",
		id,
	).Scan(&sale.ID, &sale.ClientID, &sale.Date)
	if err != nil {
		return Sale{}, err
	}
	return sale, nil
}

// SaleList returns all sales ordered by id.
func SaleList(database *sql.DB) ([]Sale, error) {
	rows, err := database.Query(
		"SELECT id, client_id, date FROM sale ORDER BY id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []Sale
	for rows.Next() {
		var sale Sale
		if err := rows.Scan(&sale.ID, &sale.ClientID, &sale.Date); err != nil {
			return nil, err
		}
		sales = append(sales, sale)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sales, nil
}

// Update updates the sale identified by s.ID.
func (s Sale) Update(database *sql.DB) error {
	_, err := database.Exec(
		"UPDATE sale SET client_id = ?, date = ? WHERE id = ?",
		s.ClientID, s.Date, s.ID,
	)
	return err
}

// SaleDelete removes the sale with the given id.
func SaleDelete(database *sql.DB, id int) error {
	_, err := database.Exec("DELETE FROM sale WHERE id = ?", id)
	return err
}

// Create inserts a purchase and stores the generated id in p.ID.
// Pointer receiver so the generated id propagates back to the caller.
func (p *Purchase) Create(database *sql.DB) error {
	result, err := database.Exec(
		"INSERT INTO purchase (sale_id, service_id, unit_price) VALUES (?, ?, ?)",
		p.SaleID, p.ServiceID, p.UnitPrice,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil
}

// PurchaseByID returns the purchase with the given id.
func PurchaseByID(database *sql.DB, id int) (Purchase, error) {
	var purchase Purchase
	err := database.QueryRow(
		"SELECT id, sale_id, service_id, unit_price FROM purchase WHERE id = ?",
		id,
	).Scan(&purchase.ID, &purchase.SaleID, &purchase.ServiceID, &purchase.UnitPrice)
	if err != nil {
		return Purchase{}, err
	}
	return purchase, nil
}

// PurchaseList returns all purchases ordered by id.
func PurchaseList(database *sql.DB) ([]Purchase, error) {
	rows, err := database.Query(
		"SELECT id, sale_id, service_id, unit_price FROM purchase ORDER BY id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []Purchase
	for rows.Next() {
		var purchase Purchase
		if err := rows.Scan(&purchase.ID, &purchase.SaleID, &purchase.ServiceID, &purchase.UnitPrice); err != nil {
			return nil, err
		}
		purchases = append(purchases, purchase)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return purchases, nil
}

// Update updates the purchase identified by p.ID.
func (p Purchase) Update(database *sql.DB) error {
	_, err := database.Exec(
		"UPDATE purchase SET sale_id = ?, service_id = ?, unit_price = ? WHERE id = ?",
		p.SaleID, p.ServiceID, p.UnitPrice, p.ID,
	)
	return err
}

// PurchaseDelete removes the purchase with the given id.
func PurchaseDelete(database *sql.DB, id int) error {
	_, err := database.Exec("DELETE FROM purchase WHERE id = ?", id)
	return err
}
