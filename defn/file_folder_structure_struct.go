package defn

type FileFolderStructure struct {
	Name    string
	Folders []*FileFolderStructure
	Files   []*FileStructure
}

type FileStructure struct {
	FileName    string
	FileType    string
	FileContent []byte
}
