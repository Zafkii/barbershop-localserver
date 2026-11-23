package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"barbershop-localserver/db"
	"barbershop-localserver/handlers"
)

func main() {
	// Configuración de la DB - parámetros en código
	dbPath := os.Getenv("BARBERHOP_DB")
	if dbPath == "" {
		dbPath = "./barbershop.db"
	}

	// Verificar si es primera ejecución (archivo no existe)
	isFirstRun := false
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		isFirstRun = true
		fmt.Printf("Primera ejecución detectada. Creando base de datos en: %s\n", dbPath)
	}

	// Abrir DB (esto ejecuta migraciones automáticamente)
	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	// Verificar si hay admin(s) en la tabla
	adminExists, err := hasAdmin(database)
	if err != nil {
		log.Fatalf("failed to check admin table: %v", err)
	}

	// Si es primera ejecución O no hay admins, crear admin inicial
	if isFirstRun || !adminExists {
		fmt.Println("Configurando administrador inicial...")
		if err := ensureInitialAdmin(database); err != nil {
			log.Fatalf("failed to create initial admin: %v", err)
		}
		fmt.Println("Administrador inicial creado correctamente.")
	}

	// Iniciar servidor
	router := handlers.NewRouter(database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	fmt.Printf("Servidor escuchando en %s (DB: %s)\n", addr, dbPath)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// hasAdmin verifica si existe al menos un admin en la tabla
func hasAdmin(database *sql.DB) (bool, error) {
	var count int
	err := database.QueryRow("SELECT COUNT(*) FROM admin").Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ensureInitialAdmin crea el admin inicial
// Prioridad: 1) Variables de entorno BARBERHOP_ADMIN_USER/BARBERHOP_ADMIN_PASS
//           2) Prompt interactivo
func ensureInitialAdmin(database *sql.DB) error {
	username := os.Getenv("BARBERHOP_ADMIN_USER")
	password := os.Getenv("BARBERHOP_ADMIN_PASS")

	// Si no vienen por env, pedir interactivamente
	if username == "" || password == "" {
		reader := bufio.NewReader(os.Stdin)

		if username == "" {
			fmt.Print("Usuario admin: ")
			u, _ := reader.ReadString('\n')
			username = u[:len(u)-1] // trim newline
		}

		if password == "" {
			fmt.Print("Contraseña: ")
			p, _ := reader.ReadString('\n')
			password = p[:len(p)-1]
		}
	}

	if username == "" || password == "" {
		return fmt.Errorf("usuario y contraseña son requeridos")
	}

	_, err := db.CreateAdmin(database, username, password)
	return err
}