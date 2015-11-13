package lea

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 分离文件名与扩展名(包含.)
func SplitFilename(filename string) (baseName, ext string) {
	baseName = filename
	// 找到最后一个'.'
	ext = SubstringByte(filename, strings.LastIndex(filename, "."))
	baseName = strings.TrimRight(filename, ext)
	ext = strings.ToLower(ext)
	return
}

// 转换文件的格式
// toExt包含.
func TransferExt(path string, toExt string) string {
	dir := filepath.Dir(path) + "/" // 文件路径
	name := filepath.Base(path)     // 文件名 a.jpg
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
	return 0
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

func MkdirAll(dir string) bool {
	err := os.MkdirAll(dir, 0777)
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

func CopyDir(source string, dest string) (err error) {
	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				//				fmt.Println(err)
			}
		} else {
			// perform copy
			_, err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				//				fmt.Println(err)
			}
		}
	}
	return
}

func DeleteFile(path string) bool {
	err := os.Remove(path)
	if err != nil {
		return false
	}
	return true
}

func IsDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
	return false
}

// 获得文件str内容
func GetFileStrContent(path string) string {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(fileBytes)
}

func IsFileExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// 写入string内容
func PutFileStrContent(path, content string) bool {
	var f *os.File
	var err1 error
	defer (func() {
		if f != nil {
			f.Close()
		}
	})()
	f, err1 = os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666) //打开文件
	//	Log(err1)
	//	var n int
	_, err1 = io.WriteString(f, content) //写入文件(字符串)
	//	Log(content)
	//	Log(err1)
	//	Log(n)
	//	Log(path)

	if err1 != nil {
		Log(err1)
		return false
	}
	return true
}
