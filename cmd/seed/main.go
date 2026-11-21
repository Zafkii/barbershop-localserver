package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"barbershop-localserver/db"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run cmd/seed/main.go <db-path>")
		fmt.Println("Ejemplo: go run cmd/seed/main.go ./barbershop.db")
		os.Exit(1)
	}

	dbPath := os.Args[1]

	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("Error abriendo DB: %v", err)
	}
	defer database.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Usuario admin: ")
	username, _ := reader.ReadString('\n')
	username = username[:len(username)-1] // trim newline

	fmt.Print("Contraseña: ")
	password, _ := reader.ReadString('\n')
	password = password[:len(password)-1]

	if username == "" || password == "" {
		log.Fatal("Usuario y contraseña requeridos")
	}

	id, err := db.CreateAdmin(database, username, password)
	if err != nil {
		log.Fatalf("Error creando admin: %v", err)
	}

	fmt.Printf("Admin creado con ID: %d\n", id)
}