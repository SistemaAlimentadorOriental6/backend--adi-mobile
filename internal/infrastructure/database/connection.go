package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/microsoft/go-mssqldb"
	"github.com/joho/godotenv"
)

// InitDB inicializa la conexión a la base de datos MySQL
func InitDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: No se pudo cargar el archivo .env")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error al abrir MySQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error al conectar a MySQL: %v", err)
	}

	log.Println("✅ Conexión a MySQL establecida correctamente")
	return db
}

// InitSQLServer inicializa la conexión a SQL Server
func InitSQLServer() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: No se pudo cargar el archivo .env")
	}

	dbUser := os.Getenv("SQLSERVER_USER")
	dbPass := os.Getenv("SQLSERVER_PASSWORD")
	dbHost := os.Getenv("SQLSERVER_HOST")
	dbPort := os.Getenv("SQLSERVER_PORT")
	dbName := os.Getenv("SQLSERVER_DATABASE")

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		log.Fatalf("Error al abrir SQL Server: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error al conectar a SQL Server: %v", err)
	}

	log.Println("✅ Conexión a SQL Server establecida correctamente")
	return db
}
