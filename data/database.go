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

var GetDatabaseConnection func() Database

func (db *Database) GetPgxPoolConnection(ctx context.Context) (*pgxpool.Conn, error) {
	return db.Pool.Acquire(ctx)
}

func (db *Database) InitialiseDatabaseTables(ctx context.Context) error {
	log := util.GetGlobalLogger(ctx)
	// if err := db.createUserProfileTable(ctx); err != nil {
	// 	return err
	// }

	// log.Println("UserProfile table successfully created")
	if err := db.createScrapeJobTable(ctx); err != nil {
		return err
	}
	log.Println("scrape_job table successfully created")

	if err := db.createFileDataTable(ctx); err != nil {
		return err
	}
	log.Println("file_data table successfully created")
	return nil
}

func (db *Database) createUserProfileTable(ctx context.Context) error {
	// createTableSQL := `CREATE TABLE IF NOT EXISTS UserProfile (
	// 		id varchar(26) PRIMARY KEY CONSTRAINT ulid_size	CHECK (char_length(id) = 26),
	// 		first_name varchar(100) NOT NULL,
	// 		last_name varchar(100),
	// 		email varchar(100) NOT NULL,
	// 		password varchar(100) NOT NULL
	// 	);`
	// if _, err := db.Pool.Exec(ctx, createTableSQL); err != nil {
	// 	util.GetGlobalLogger(ctx).Println("Failed to execute create query", err)
	// 	return err
	// }
	return nil
}

func (db *Database) createScrapeJobTable(ctx context.Context) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS scrape_job(
		id VARCHAR(26) PRIMARY KEY CONSTRAINT ulid_size	CHECK (char_length(id) = 26),
		url VARCHAR(255) NOT NULL,
		depth integer default 1 NOT NULL,
		maxlimit integer default 1 NOT NULL,
		response jsonb,
		created_on timestamp default NOW()
	);`
	if _, err := db.Pool.Exec(ctx, createTableSQL); err != nil {
		util.GetGlobalLogger(ctx).Println("Failed to execute create query", err)
		return err
	}
	return nil
}

func (db *Database) createFileDataTable(ctx context.Context) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS file_data(
		id VARCHAR(26) PRIMARY KEY CONSTRAINT ulid_size	CHECK (char_length(id) = 26),
		file_name VARCHAR(255) NOT NULL,
		file_type VARCHAR(127) NOT NULL,
		file_path VARCHAR(511) NOT NULL,
		file_size numeric NOT NULL,
		created_on timestamp default NOW()
	);`
	if _, err := db.Pool.Exec(ctx, createTableSQL); err != nil {
		util.GetGlobalLogger(ctx).Println("Failed to execute create query", err)
		return err
	}
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

	GetDatabaseConnection = func() Database {
		return Database{Pool: conn}
	}
	return Database{Pool: conn}
}
