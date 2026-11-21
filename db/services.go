package db

import (
	"database/sql"
)

// Create inserts a service and stores the generated id in s.ID.
// Pointer receiver so the generated id propagates back to the caller.
func (s *Service) Create(database *sql.DB) error {
	result, err := database.Exec(
		"INSERT INTO service (name, service_type) VALUES (?, ?)",
		s.Name, s.ServiceType,
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

// ServiceByID returns the service with the given id.
func ServiceByID(database *sql.DB, id int) (Service, error) {
	var service Service
	err := database.QueryRow(
		"SELECT id, name, service_type FROM service WHERE id = ?",
		id,
	).Scan(&service.ID, &service.Name, &service.ServiceType)
	if err != nil {
		return Service{}, err
	}
	return service, nil
}

// ServiceList returns all services ordered by id.
func ServiceList(database *sql.DB) ([]Service, error) {
	rows, err := database.Query(
		"SELECT id, name, service_type FROM service ORDER BY id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var service Service
		if err := rows.Scan(&service.ID, &service.Name, &service.ServiceType); err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return services, nil
}

// Update updates the service identified by s.ID.
func (s Service) Update(database *sql.DB) error {
	_, err := database.Exec(
		"UPDATE service SET name = ?, service_type = ? WHERE id = ?",
		s.Name, s.ServiceType, s.ID,
	)
	return err
}

// ServiceDelete removes the service with the given id.
func ServiceDelete(database *sql.DB, id int) error {
	_, err := database.Exec("DELETE FROM service WHERE id = ?", id)
	return err
}
