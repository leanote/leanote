package lea

import (
	"strings"
	"path/filepath"
	"os"
	"io"
)

// 分离文件名与扩展名(包含.)
func SplitFilename(filename string) (baseName, ext string) {
	baseName = filename
	// 找到最后一个'.'
	ext = SubstringByte(filename, strings.LastIndex(filename, "."))
	baseName = strings.TrimRight(filename, ext)
	ext = strings.ToLower(ext)
	return;
}

// 转换文件的格式
// toExt包含.
func TransferExt(path string, toExt string) string {
	dir := filepath.Dir(path) + "/"  // 文件路径
	name := filepath.Base(path) // 文件名 a.jpg
	// 获取文件名与路径
	baseName, _ := SplitFilename(name)
	return dir + baseName + toExt
}

func GetFilename(path string) string {
	return filepath.Base(path)
}

// file size
// length in bytes
func GetFilesize(path string) int64 {
	fileinfo, err := os.Stat(path)
	if err == nil {
		return fileinfo.Size()
	}
	return 0;
}

// 清空dir下所有的文件和文件夹
// RemoveAll会清空本文件夹, 所以还要创建之
func ClearDir(dir string) bool {
	err := os.RemoveAll(dir)
	if err != nil {
		return false
	}
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return false
	}
	return true
}

// list dir's all file, return filenames
func ListDir(dir string) []string {
	f, err := os.Open(dir)
	if err != nil {
		return nil
	}
	names, _ := f.Readdirnames(0)
	return names
}

func CopyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func IsDirExists(path string) bool {
    fi, err := os.Stat(path)
    if err != nil {
        return os.IsExist(err)
    }else{
        return fi.IsDir()
    }
    return false
}