package db

import (
	"database/sql"
)

// Create inserts a client and stores the generated id in c.ID.
// Pointer receiver so the generated id propagates back to the caller.
func (c *Client) Create(database *sql.DB) error {
	result, err := database.Exec(
		"INSERT INTO client (name, last_name, phone) VALUES (?, ?, ?)",
		c.Name, c.LastName, c.Phone,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = int(id)
	return nil
}

// ClientByID returns the client with the given id.
func ClientByID(database *sql.DB, id int) (Client, error) {
	var client Client
	err := database.QueryRow(
		"SELECT id, name, last_name, phone FROM client WHERE id = ?",
		id,
	).Scan(&client.ID, &client.Name, &client.LastName, &client.Phone)
	if err != nil {
		return Client{}, err
	}
	return client, nil
}

// ClientList returns all clients ordered by id.
func ClientList(database *sql.DB) ([]Client, error) {
	rows, err := database.Query(
		"SELECT id, name, last_name, phone FROM client ORDER BY id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var client Client
		if err := rows.Scan(&client.ID, &client.Name, &client.LastName, &client.Phone); err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return clients, nil
}

// Update updates the client identified by c.ID.
func (c Client) Update(database *sql.DB) error {
	_, err := database.Exec(
		"UPDATE client SET name = ?, last_name = ?, phone = ? WHERE id = ?",
		c.Name, c.LastName, c.Phone, c.ID,
	)
	return err
}

// ClientDelete removes the client with the given id.
func ClientDelete(database *sql.DB, id int) error {
	_, err := database.Exec("DELETE FROM client WHERE id = ?", id)
	return err
}
