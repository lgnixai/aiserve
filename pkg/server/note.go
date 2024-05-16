package server

import (
	`encoding/json`
	`fmt`
	`net/http`
	`path/filepath`

	`github.com/gin-gonic/gin`
	`github.com/spf13/afero`
)

type FileNode struct {
	Name     string      `json:"name"`
	IsDir    bool        `json:"is_dir"`
	Children []*FileNode `json:"children,omitempty"`
}

func (s *Server) NoteTree(ctx *gin.Context) {
	fs := afero.NewBasePathFs(afero.NewOsFs(), "/ai/note")

	// 在内存文件系统中创建一些文件和文件夹
	afero.WriteFile(fs, "file1.txt", []byte("Hello, World!"), 0644)
	fs.Mkdir("dir1", 0755)
	afero.WriteFile(fs, "dir1/file2.txt", []byte("This is file2.txt"), 0644)
	fs.Mkdir("dir1/subdir", 0755)
	afero.WriteFile(fs, "dir1/subdir/file3.txt", []byte("This is file3.txt"), 0644)

	// 构建文件树
	root, err := buildFileTree(fs, "/")
	if err != nil {
		fmt.Println("Error building file tree:", err)
		return
	}

	// 将文件树转换为 JSON 格式并输出
	jsonData, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		fmt.Println("Error encoding to JSON:", err)
		return
	}
	ctx.JSON(http.StatusOK, jsonData)

}
func buildFileTree(fs afero.Fs, path string) (*FileNode, error) {
	fileInfo, err := fs.Stat(path)
	if err != nil {
		return nil, err
	}

	node := &FileNode{
		Name:  fileInfo.Name(),
		IsDir: fileInfo.IsDir(),
	}

	if node.IsDir {
		files, err := afero.ReadDir(fs, path)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			childPath := filepath.Join(path, file.Name())
			childNode, err := buildFileTree(fs, childPath)
			if err != nil {
				return nil, err
			}
			node.Children = append(node.Children, childNode)
		}
	}

	return node, nil
}
