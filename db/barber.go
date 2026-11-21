package db

import (
	"database/sql"
)

// Create inserts a barber and stores the generated id in b.ID.
// Pointer receiver so the generated id propagates back to the caller.
func (b *Barber) Create(database *sql.DB) error {
	result, err := database.Exec(
		"INSERT INTO barber (name, last_name, phone) VALUES (?, ?, ?)",
		b.Name, b.LastName, b.Phone,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	b.ID = int(id)
	return nil
}

// BarberByID returns the barber with the given id.
func BarberByID(database *sql.DB, id int) (Barber, error) {
	var barber Barber
	err := database.QueryRow(
		"SELECT id, name, last_name, phone FROM barber WHERE id = ?",
		id,
	).Scan(&barber.ID, &barber.Name, &barber.LastName, &barber.Phone)
	if err != nil {
		return Barber{}, err
	}
	return barber, nil
}

// BarberList returns all barbers ordered by id.
func BarberList(database *sql.DB) ([]Barber, error) {
	rows, err := database.Query(
		"SELECT id, name, last_name, phone FROM barber ORDER BY id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var barbers []Barber
	for rows.Next() {
		var barber Barber
		if err := rows.Scan(&barber.ID, &barber.Name, &barber.LastName, &barber.Phone); err != nil {
			return nil, err
		}
		barbers = append(barbers, barber)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return barbers, nil
}

// Update updates the barber identified by b.ID.
func (b Barber) Update(database *sql.DB) error {
	_, err := database.Exec(
		"UPDATE barber SET name = ?, last_name = ?, phone = ? WHERE id = ?",
		b.Name, b.LastName, b.Phone, b.ID,
	)
	return err
}

// BarberDelete removes the barber with the given id.
func BarberDelete(database *sql.DB, id int) error {
	_, err := database.Exec("DELETE FROM barber WHERE id = ?", id)
	return err
}
