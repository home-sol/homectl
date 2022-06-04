package fs

import (
	"io/ioutil"
	"os"
	"path"
)

type FileSystem struct {
	baseDir string
}

func Cwd() (*FileSystem, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &FileSystem{baseDir: cwd}, nil
}

func FromDir(dir string) (*FileSystem, error) {
	return &FileSystem{baseDir: dir}, nil
}

func (fs *FileSystem) GetRelativePath(relativePath string) string {
	return path.Join(fs.baseDir, relativePath)
}

func (fs *FileSystem) IsDirectory(dir string) (bool, error) {
	fileInfo, err := os.Stat(fs.GetRelativePath(dir))
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

// FileExists checks if a file exists and is not a directory
func (fs *FileSystem) FileExists(filename string) bool {
	fileInfo, err := os.Stat(fs.GetRelativePath(filename))
	if os.IsNotExist(err) {
		return false
	}
	return !fileInfo.IsDir()
}

func (fs *FileSystem) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(fs.GetRelativePath(filename))
}

func (fs *FileSystem) Remove(name string) error {
	return os.Remove(fs.GetRelativePath(name))
}
