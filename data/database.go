package data

import (
	"backend-service/util"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

type Database struct {
	Pool *pgxpool.Pool
}

var GetDatabaseConnection func() Database

func (db *Database) GetPgxPoolConnection(ctx context.Context) (*pgxpool.Conn, error) {
	return db.Pool.Acquire(ctx)
}

func ConnectDatabase(ctx context.Context) Database {
	log := util.GetGlobalLogger(ctx)

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("missing required environment variables for database connection")
	}
	// credentials := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=verify-full sslrootcert=%s sslkey=%s sslcert=%s", host, user, password, database, port, dBRootCertFilePath, dBKeyFilePath, dBCertFilePath)
	credentials := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", dbHost, dbUser, dbPassword, dbName, dbPort)

	conn, err := pgxpool.New(ctx, credentials)
	if err != nil {
		log.Fatalln("could not connect to database", err)
	} else if err := conn.Ping(ctx); err != nil {
		log.Fatalln("could not ping database", err)
	}

	GetDatabaseConnection = func() Database {
		return Database{Pool: conn}
	}
	return Database{Pool: conn}
}
