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
	var uploadedImages []map[string]interface{}

	if strings.EqualFold(fileFolderMap.Name, "") {
		cerr := util.NewCustomError(ctx, "no-folder-name-provided", errors.New("no folder name was provided"))
		log.Println(cerr)
		return nil, cerr
	}

	for _, folder := range fileFolderMap.Folders {
		savedFiles, cerr := ParseFolderStructureAndSaveFile(ctx, filepath.Join(folderPath, fileFolderMap.Name), folder)
		if strings.EqualFold(folder.Name, "images") {
			if savedFiles != nil {
				uploadedImages = savedFiles
			}
			if cerr != nil {
				return append(savedFilesMap, map[string]interface{}{"uploaded_images": uploadedImages}), cerr
			}
		} else {
			if savedFiles != nil {
				savedFilesMap = append(savedFilesMap, savedFiles...)
			}
			if cerr != nil {
				return savedFilesMap, cerr
			}
		}
	}

	for _, file := range fileFolderMap.Files {
		savedFile, cerr := SaveFile(ctx, filepath.Join(folderPath, fileFolderMap.Name), file)
		if cerr != nil {
			return savedFilesMap, cerr
		}

		if uploadedImages != nil {
			savedFile["uploaded_images"] = uploadedImages
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
	fileExists := false
	if _, err := os.Stat(filePath); !errors.Is(err, os.ErrNotExist) {
		fileExists = true
	}

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
	if !fileExists {
		if fileId, cerr := fileRepo.Create(ctx, fileMetaData); cerr != nil {
			log.Println(cerr)
			return nil, cerr
		} else {
			return map[string]interface{}{
				"id":   fileId,
				"name": file.FileName + file.FileType,
			}, nil
		}
	} else {
		if fileId, cerr := fileRepo.UpdateFileSizeByFilePath(ctx, fileMetaData); cerr != nil {
			log.Println(cerr)
			return nil, cerr
		} else {
			return map[string]interface{}{
				"id":   fileId,
				"name": file.FileName + file.FileType,
			}, nil
		}
	}
}

func GetFile(ctx context.Context, fileId string) ([]byte, *defn.FileInfo, *util.CustomError) {
	log := util.GetGlobalLogger(ctx)

	fileMap, cerr := data.NewFileRepo().GetFileById(ctx, fileId)
	if cerr != nil {
		log.Println(cerr)
		return nil, nil, cerr
	}

	fileInfo := &defn.FileInfo{
		FileName: fileMap["file_name"].(string),
		FilePath: fileMap["file_path"].(string),
		FileType: fileMap["file_type"].(string),
	}

	fileContent, err := os.ReadFile(fileInfo.FilePath)
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeFileReadFailed, defn.ErrFileReadFailed, map[string]string{
			"error": err.Error(),
			"path":  fileInfo.FilePath,
		})
		log.Println(cerr)
		return nil, nil, cerr
	}
	return fileContent, fileInfo, nil
}
