package data

import (
	"backend-service/defn"
	"backend-service/util"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type ScrapeTaskRepo struct {
	db Database
}

func NewScrapeTaskRepo() *ScrapeTaskRepo {
	return &ScrapeTaskRepo{db: GetDatabaseConnection()}
}

func (repo *ScrapeTaskRepo) Create(ctx context.Context, task defn.ScrapeTask) (string, *util.CustomError) {
	task.ID = util.ULID()

	var responseBytes []byte
	var err error
	if task.Response != nil {
		responseBytes, err = json.Marshal(task.Response)
		if err != nil {
			cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseCreateOperationFailed, defn.ErrDatabaseCreateOperationFailed, map[string]string{
				"error": err.Error(),
			})
			// log.Println(cerr)
			return "", cerr
		}
	}

	_, err = repo.db.Pool.Exec(ctx, "INSERT INTO scrape_task (id, job_id, url, depth, maxlimit, level, response) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		task.ID, task.JobId, task.URL, task.Depth, task.Maxlimit, task.Level, responseBytes)
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseCreateOperationFailed, defn.ErrDatabaseCreateOperationFailed, map[string]string{
			"error": err.Error(),
		})
		// log.Println(cerr)
		return "", cerr
	}
	return task.ID, nil
}

func (repo *ScrapeTaskRepo) Update(ctx context.Context, id string, updates map[string]interface{}) (map[string]interface{}, *util.CustomError) {
	if len(updates) == 0 {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseUpdateOperationFailed, defn.ErrDatabaseUpdateOperationFailed, map[string]string{
			"error": "no fields found for updating",
		})
		// log.Println(cerr)
		return nil, cerr
	}

	setClauses := []string{}
	args := []interface{}{id}
	argIndex := 2

	for field, value := range updates {
		if field == "response" {
			jsobBytes, err := json.Marshal(value)
			if err != nil {
				cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseUpdateOperationFailed, defn.ErrDatabaseUpdateOperationFailed, map[string]string{
					"error": err.Error(),
				})
				// log.Println(cerr)
				return nil, cerr
			}
			value = jsobBytes
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf("UPDATE scrape_task SET %s WHERE id = $1 RETURNING *", strings.Join(setClauses, ", "))

	rows, err := repo.db.Pool.Query(ctx, query, args...)
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseUpdateOperationFailed, defn.ErrDatabaseUpdateOperationFailed, map[string]string{
			"error": err.Error(),
		})
		// log.Println(cerr)
		return nil, cerr
	}
	defer rows.Close()

	// Get column names dynamically
	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = fd.Name
	}

	// Read the first and only row
	if !rows.Next() {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseUpdateOperationFailed, defn.ErrDatabaseUpdateOperationFailed, map[string]string{
			"error": "no rows found",
		})
		// log.Println(cerr)
		return nil, cerr
	}

	values, err := rows.Values() // Get all values in a slice
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseUpdateOperationFailed, defn.ErrDatabaseUpdateOperationFailed, map[string]string{
			"error": err.Error(),
		})
		// log.Println(cerr)
		return nil, cerr
	}

	// Map column names to values
	updatedData := make(map[string]interface{})
	for i, column := range columns {
		updatedData[column] = values[i]
	}

	return updatedData, nil
}

func (repo *ScrapeTaskRepo) Get(ctx context.Context, id string) (*defn.ScrapeJob, *util.CustomError) {
	return nil, nil
}

// func (repo *ScrapeTaskRepo) Delete(ctx context.Context, job defn.ScrapeJob) (defn.ScrapeJob, *util.CustomError) {
// 	return "", nil
// }
