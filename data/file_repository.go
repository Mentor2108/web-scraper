package data

import (
	"backend-service/defn"
	"backend-service/util"
	"context"
)

type FileRepo struct {
	db Database
}

func NewFileRepo() *FileRepo {
	return &FileRepo{db: GetDatabaseConnection()}
}

func (repo *FileRepo) Create(ctx context.Context, fileInfo defn.FileInfo) (string, *util.CustomError) {
	fileInfo.ID = util.ULID()

	_, err := repo.db.Pool.Exec(ctx, "INSERT INTO file_data (id, file_name, file_type, file_path, file_size) VALUES ($1, $2, $3, $4, $5)",
		fileInfo.ID, fileInfo.FileName, fileInfo.FileType, fileInfo.FilePath, fileInfo.FileSize)
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseCreateOperationFailed, defn.ErrDatabaseCreateOperationFailed, map[string]string{
			"error": err.Error(),
		})
		// log.Println(cerr)
		return "", cerr
	}
	return fileInfo.ID, nil
}

func (repo *FileRepo) GetFileById(ctx context.Context, fileId string) (map[string]interface{}, *util.CustomError) {
	rows, err := repo.db.Pool.Query(ctx, "SELECT * from file_data where id = $1", fileId)
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
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDatabaseGetOperationFailed, defn.ErrDatabaseGetOperationFailed, map[string]string{
			"error": "no file found with id " + fileId,
		})
		// log.Println(cerr)
		return nil, cerr
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
	fileResponse := make(map[string]interface{})
	for i, column := range columns {
		fileResponse[column] = values[i]
	}
	return fileResponse, nil
}
