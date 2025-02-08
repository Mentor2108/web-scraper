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

func (repo *FileRepo) GetFileById(ctx context.Context, id string) (map[string]interface{}, *util.CustomError) {
	return nil, nil
}
