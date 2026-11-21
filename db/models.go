package db

// Barber represents a barber shop barber.
type Barber struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Phone    string `json:"phone"`
}

// Client represents a barber shop client.
type Client struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Phone    string `json:"phone"`
}

// Service represents a service offered by the barber shop.
type Service struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ServiceType string `json:"service_type"`
}

// Sale represents a sale made to a client.
type Sale struct {
	ID       int    `json:"id"`
	ClientID int    `json:"client_id"`
	Date     string `json:"date"`
}

// Purchase represents a service purchased within a sale.
type Purchase struct {
	ID        int     `json:"id"`
	SaleID    int     `json:"sale_id"`
	ServiceID int     `json:"service_id"`
	UnitPrice float64 `json:"unit_price"`
}

// Admin represents an admin user.
type Admin struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // never serialized
}
