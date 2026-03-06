package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	host := "localhost"
	port := "5432"
	user := "root"
	password := "root"
	dbname := "patounes_db"
	sslmode := "disable"

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	log.Println("Tentative de connexion à PostgreSQL...")

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("Erreur lors de l'ouverture de la connexion: %w", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("Impossible de se connecter à la base de données: %w", err)
	}

	log.Println("Connexion à PostgreSQL réussie !")
	log.Printf("Base de données: %s@%s:%s/%s", user, host, port, dbname)

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	return nil
}

func Close() error {
	if DB != nil {
		log.Println("Fermeture de la connexion à la base de données...")
		return DB.Close()
	}
	return nil
}
