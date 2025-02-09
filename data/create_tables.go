package data

import (
	"backend-service/util"
	"context"
)

func (db *Database) InitialiseDatabaseTables(ctx context.Context) error {
	log := util.GetGlobalLogger(ctx)
	if err := db.createScrapeJobTable(ctx); err != nil {
		return err
	}
	log.Println("scrape_job table successfully created")

	if err := db.createScrapeTaskTable(ctx); err != nil {
		return err
	}
	log.Println("scrape_task table successfully created")

	if err := db.createFileDataTable(ctx); err != nil {
		return err
	}
	log.Println("file_data table successfully created")
	return nil
}

func (db *Database) createScrapeJobTable(ctx context.Context) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS scrape_job(
		id VARCHAR(26) PRIMARY KEY CONSTRAINT ulid_size	CHECK (char_length(id) = 26),
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

func (db *Database) createScrapeTaskTable(ctx context.Context) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS scrape_task(
		id VARCHAR(26) PRIMARY KEY,
		job_id VARCHAR(26) NOT NULL,
		url VARCHAR(255) NOT NULL,
		depth integer default 1 NOT NULL,
		maxlimit integer default 1 NOT NULL,
		level integer NOT NULL,
		response jsonb,
		created_on timestamp default NOW(),
		CONSTRAINT ulid_size CHECK (char_length(id) = 26),
		CONSTRAINT foreign_key_size CHECK (char_length(job_id) = 26),
		CONSTRAINT fk_scrapeJob_scrapeTask FOREIGN KEY(job_id) REFERENCES scrape_job(id)
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
