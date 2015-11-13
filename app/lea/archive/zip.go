package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// main functions shows how to TarGz a directory/file and
// UnTarGz a file
// Gzip and tar from source directory or file to destination file
// you need check file exist before you call this function

func Zip(srcDirPath string, destFilePath string) (ok bool) {
	defer func() { //必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			ok = false
		}
	}()

	fw, err := os.Create(destFilePath)

	if err != nil {
		panic(err)
	}
	defer fw.Close()

	// Tar writer
	tw := zip.NewWriter(fw)
	defer tw.Close()

	// Check if it's a file or a directory
	f, err := os.Open(srcDirPath)
	if err != nil {
		panic(err)
	}
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	if fi.IsDir() {
		// handle source directory
		//        fmt.Println("Cerating tar.gz from directory...")
		zipDir(srcDirPath, path.Base(srcDirPath), tw)
	} else {
		// handle file directly
		//        fmt.Println("Cerating tar.gz from " + fi.Name() + "...")
		zipFile(srcDirPath, fi.Name(), tw, fi)
	}
	ok = true
	return
}

// Deal with directories
// if find files, handle them with zipFile
// Every recurrence append the base path to the recPath
// recPath is the path inside of tar.gz
func zipDir(srcDirPath string, recPath string, tw *zip.Writer) {
	// Open source diretory
	dir, err := os.Open(srcDirPath)
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	// Get file info slice
	fis, err := dir.Readdir(0)
	if err != nil {
		panic(err)
	}
	for _, fi := range fis {
		// Append path
		curPath := srcDirPath + "/" + fi.Name()
		// Check it is directory or file
		if fi.IsDir() {
			// Directory
			// (Directory won't add unitl all subfiles are added)
			//            fmt.Printf("Adding path...%s\n", curPath)
			zipDir(curPath, recPath+"/"+fi.Name(), tw)
		} else {
			// File
			//            fmt.Printf("Adding file...%s\n", curPath)
		}

		zipFile(curPath, recPath+"/"+fi.Name(), tw, fi)
	}
}

// Deal with files
func zipFile(srcFile string, recPath string, tw *zip.Writer, fi os.FileInfo) {
	if fi.IsDir() {
		//    	fmt.Println("??")
		// Create tar header
		/*
		   fh, err := zip.FileInfoHeader(fi)
		   if err != nil {
		       panic(err)
		   }
		   fh.Name = recPath // + "/"
		   err = tw.WriteHeader(hdr)
		   tw.Create(recPath)
		*/
	} else {
		// File reader
		fr, err := os.Open(srcFile)
		if err != nil {
			panic(err)
		}
		defer fr.Close()

		// Write hander
		w, err2 := tw.Create(recPath)
		if err2 != nil {
			panic(err)
		}
		// Write file data
		_, err = io.Copy(w, fr)
		if err != nil {
			panic(err)
		}
	}
}

// Ungzip and untar from source file to destination directory
// you need check file exist before you call this function
func Unzip(srcFilePath string, destDirPath string) (ok bool, msg string) {
	ok = false
	msg = ""

	defer func() { //必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			msg = fmt.Sprintf("%v", err)
			ok = false
		}
	}()

	os.Mkdir(destDirPath, os.ModePerm)
	r, err := zip.OpenReader(srcFilePath)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	for _, f := range r.File {
		//		fmt.Println("FileName : ", f.Name); // j/aaa.zip
		rc, err := f.Open()
		if err != nil {
			panic(err)
		}

		// 把首文件夹去掉, 即j去掉, 分离出文件夹和文件名
		paths := strings.Split(f.Name, "/")
		prePath := ""
		filename := ""
		l := len(paths)
		//		fmt.Println(l)
		if l > 1 {
			// 去掉第1个文件夹
			if l == 2 {
				filename = paths[1]
			} else {
				filename = paths[l-1]
				prePath = strings.Join(paths[1:l-1], "/")
			}
		} else {
			filename = f.Name
		}
		//		fmt.Println(prePath)

		// 相对于目标文件件下的路径
		destPath := destDirPath + "/" + filename
		if prePath != "" {
			os.MkdirAll(destDirPath+"/"+prePath, os.ModePerm)
			destPath = destDirPath + "/" + prePath + "/" + filename
		}
		// Write data to file
		//        fmt.Println(destPath)
		fw, _ := os.Create(destPath)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(fw, rc)
		if err != nil {
			panic(err)
		}
	}

	ok = true
	return
}
