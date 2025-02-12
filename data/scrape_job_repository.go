package data

import (
	"backend-service/defn"
	"backend-service/util"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type ScrapeJobRepo struct {
	db Database
}

func NewScrapeJobRepo() *ScrapeJobRepo {
	return &ScrapeJobRepo{db: GetDatabaseConnection()}
}

func (repo *ScrapeJobRepo) Create(ctx context.Context, job defn.ScrapeJob) (string, *util.CustomError) {
	job.ID = util.ULID()

	var responseBytes []byte
	var err error
	if job.Response != nil {
		responseBytes, err = json.Marshal(job.Response)
		if err != nil {
			cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseCreateOperationFailed, defn.ErrDatabaseCreateOperationFailed, map[string]string{
				"error": err.Error(),
			})
			// log.Println(cerr)
			return "", cerr
		}
	}

	_, err = repo.db.Pool.Exec(ctx, "INSERT INTO scrape_job (id, depth, maxlimit, response) VALUES ($1, $2, $3, $4)",
		job.ID, job.Depth, job.Maxlimit, responseBytes)
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseCreateOperationFailed, defn.ErrDatabaseCreateOperationFailed, map[string]string{
			"error": err.Error(),
		})
		// log.Println(cerr)
		return "", cerr
	}
	return job.ID, nil
}

func (repo *ScrapeJobRepo) Update(ctx context.Context, id string, updates map[string]interface{}) (map[string]interface{}, *util.CustomError) {
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

	query := fmt.Sprintf("UPDATE scrape_job SET %s WHERE id = $1 RETURNING *", strings.Join(setClauses, ", "))

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

// func (repo *ScrapeJobRepo) Get(ctx context.Context, id string) (*defn.ScrapeJob, *util.CustomError) {

// 	return nil, nil
// }

func (repo *ScrapeJobRepo) GetJobWithTasks(ctx context.Context, pagesize int) (map[string]interface{}, *util.CustomError) {
	query := `
	SELECT 
		j.*, 
		COALESCE(jsonb_agg(t.*) FILTER (WHERE t.id IS NOT NULL), '[]') AS scrape_task
	FROM scrape_job j
	LEFT JOIN scrape_task t ON j.id = t.job_id
	GROUP BY j.id
	LIMIT $1;
	`

	rows, err := repo.db.Pool.Query(ctx, query, pagesize)
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseGetOperationFailed, defn.ErrDatabaseGetOperationFailed, map[string]string{
			"error": err.Error(),
		})
		// log.Println(cerr)
		return nil, cerr
	}
	defer rows.Close()

	// Fetch column names dynamically
	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = fd.Name
	}

	resp := []map[string]interface{}{}
	for rows.Next() {
		values, err := rows.Values() // Get all values in a slice
		if err != nil {
			cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseGetOperationFailed, defn.ErrDatabaseGetOperationFailed, map[string]string{
				"error": err.Error(),
			})
			// log.Println(cerr)
			return nil, cerr
		}

		// Store values dynamically in map
		jobRowResponse := make(map[string]interface{})
		for i, column := range columns {
			jobRowResponse[column] = values[i]
		}

		// Convert `tasks` column to `[]map[string]interface{}`
		if tasksJSON, ok := jobRowResponse["tasks"].([]byte); ok {
			var tasks []map[string]interface{}
			if err := json.Unmarshal(tasksJSON, &tasks); err != nil {
				cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseGetOperationFailed, defn.ErrDatabaseGetOperationFailed, map[string]string{
					"error": err.Error(),
				})
				// log.Println(cerr)
				return nil, cerr
			}
			jobRowResponse["tasks"] = tasks
		}
		resp = append(resp, jobRowResponse)
	}

	// Read the first row
	// if !rows.Next() {
	// 	cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseGetOperationFailed, defn.ErrDatabaseGetOperationFailed, map[string]string{
	// 		"error": "no job found with given id",
	// 	})
	// 	// log.Println(cerr)
	// 	return nil, cerr
	// }

	return map[string]interface{}{
		"scrape_job": resp,
	}, nil
}

func (repo *ScrapeJobRepo) GetJobWithTasksByID(ctx context.Context, jobID string) (map[string]interface{}, *util.CustomError) {
	query := `
	SELECT 
		j.*, 
		COALESCE(jsonb_agg(t.* ORDER BY t.created_on) FILTER (WHERE t.id IS NOT NULL), '[]') AS scrape_task
	FROM scrape_job j
	LEFT JOIN scrape_task t ON j.id = t.job_id
	WHERE j.id = $1
	GROUP BY j.id;
	`

	rows, err := repo.db.Pool.Query(ctx, query, jobID)
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseGetOperationFailed, defn.ErrDatabaseGetOperationFailed, map[string]string{
			"error": err.Error(),
		})
		// log.Println(cerr)
		return nil, cerr
	}
	defer rows.Close()

	// Fetch column names dynamically
	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = fd.Name
	}

	// Read the first row
	if !rows.Next() {
		return map[string]interface{}{
			"scrape_job": nil,
		}, nil
	}

	values, err := rows.Values() // Get all values in a slice
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseGetOperationFailed, defn.ErrDatabaseGetOperationFailed, map[string]string{
			"error": err.Error(),
		})
		// log.Println(cerr)
		return nil, cerr
	}

	// Store values dynamically in map
	jobResponse := make(map[string]interface{})
	for i, column := range columns {
		jobResponse[column] = values[i]
	}

	// Convert `tasks` column to `[]map[string]interface{}`
	if tasksJSON, ok := jobResponse["tasks"].([]byte); ok {
		var tasks []map[string]interface{}
		if err := json.Unmarshal(tasksJSON, &tasks); err != nil {
			cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseGetOperationFailed, defn.ErrDatabaseGetOperationFailed, map[string]string{
				"error": err.Error(),
			})
			// log.Println(cerr)
			return nil, cerr
		}
		jobResponse["tasks"] = tasks
	}

	return map[string]interface{}{
		"scrape_job": jobResponse,
	}, nil
}

// func (repo *ScrapeJobRepo) Delete(ctx context.Context, job defn.ScrapeJob) (defn.ScrapeJob, *util.CustomError) {
// 	return "", nil
// }
