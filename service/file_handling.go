package service

import (
	"backend-service/data"
	"backend-service/defn"
	"backend-service/util"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func ParseFolderStructureAndSaveFile(ctx context.Context, folderPath string, fileFolderMap *defn.FileFolderStructure) ([]map[string]interface{}, *util.CustomError) {
	if fileFolderMap == nil {
		return []map[string]interface{}{}, nil
	}

	log := util.GetGlobalLogger(ctx)

	savedFilesMap := []map[string]interface{}{}

	if strings.EqualFold(fileFolderMap.Name, "") {
		cerr := util.NewCustomError(ctx, "no-folder-name-provided", errors.New("no folder name was provided"))
		log.Println(cerr)
		return nil, cerr
	}

	for _, folder := range fileFolderMap.Folders {
		savedFiles, cerr := ParseFolderStructureAndSaveFile(ctx, filepath.Join(folderPath, fileFolderMap.Name), folder)
		if savedFiles != nil {
			savedFilesMap = append(savedFilesMap, savedFiles...)
		}
		if cerr != nil {
			return savedFilesMap, cerr
		}
	}

	for _, file := range fileFolderMap.Files {
		savedFile, cerr := SaveFile(ctx, filepath.Join(folderPath, fileFolderMap.Name), file)
		if cerr != nil {
			return savedFilesMap, cerr
		}

		savedFilesMap = append(savedFilesMap, savedFile)
	}

	return savedFilesMap, nil
}

func SaveFile(ctx context.Context, folderPath string, file *defn.FileStructure) (map[string]interface{}, *util.CustomError) {
	log := util.GetGlobalLogger(ctx)

	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeDirectoryCreationFailed, defn.ErrDirectoryCreationFailed, map[string]string{
			"error": err.Error(),
			"path":  folderPath,
		})
		log.Println(cerr)
		return nil, cerr
	}

	filePath := filepath.Join(folderPath, file.FileName+file.FileType)
	if err := os.WriteFile(filePath, file.FileContent, os.ModePerm); err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeFileWriteFailed, defn.ErrFileWriteFailed, map[string]string{
			"error": err.Error(),
			"path":  filePath,
		})
		log.Println(cerr)
		return nil, cerr
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeFileStatFailed, defn.ErrFileStatFailed, map[string]string{
			"error": err.Error(),
			"path":  filePath,
		})
		log.Println(cerr)
		return nil, cerr
	}

	fileMetaData := defn.FileInfo{
		FileName: file.FileName,
		FilePath: filePath,
		FileSize: fileInfo.Size(),
		FileType: file.FileType,
	}

	fileRepo := data.NewFileRepo()
	if fileId, cerr := fileRepo.Create(ctx, fileMetaData); cerr != nil {
		log.Println(cerr)
		return nil, cerr
	} else {
		return map[string]interface{}{
			"id":   fileId,
			"name": file.FileName + file.FileType,
		}, nil
	}
}

func GetFile(ctx context.Context) {

}
