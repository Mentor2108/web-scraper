package data

import (
	"backend-service/util"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func (db *Database) GetPgxPoolConnection(ctx context.Context) (*pgxpool.Conn, error) {
	return db.Pool.Acquire(ctx)
}

func (db *Database) InitialiseDatabaseTables(ctx context.Context) error {
	if err := db.createUserProfileTable(ctx); err != nil {
		return err
	}
	log := util.GetGlobalLogger(ctx)

	log.Println("UserProfile table successfully created")
	if err := db.createPortfolioTable(ctx); err != nil {
		return err
	}
	return nil
}

func (db *Database) createUserProfileTable(ctx context.Context) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS UserProfile (
			id varchar(26) PRIMARY KEY CONSTRAINT ulid_size	CHECK (char_length(id) = 26),
			first_name varchar(100) NOT NULL,
			last_name varchar(100),
			email varchar(100) NOT NULL,
			password varchar(100) NOT NULL
		);`
	if _, err := db.Pool.Exec(ctx, createTableSQL); err != nil {
		util.GetGlobalLogger(ctx).Println("Failed to execute create query", err)
		return err
	}
	return nil
}

func (db *Database) createPortfolioTable(ctx context.Context) error {
	// createTableSQL := `CREATE TABLE IF NOT EXISTS portfolio_items(
	// 	id VARCHAR(26) PRIMARY KEY CONSTRAINT ulid_size	CHECK (char_length(id) = 26),
	// 	title VARCHAR(100),
	// 	description TEXT,
	// 	url VARCHAR(255)
	// );`
	// if _, err := db.Pool.Exec(ctx, createTableSQL); err != nil {
	// 	util.GetGlobalLogger(ctx).Println("Failed to execute create query", err)
	// 	return err
	// }
	return nil
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

	return Database{Pool: conn}
}
