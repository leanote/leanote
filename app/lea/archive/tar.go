package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
)

// main functions shows how to TarGz a directory/file and
// UnTarGz a file
// Gzip and tar from source directory or file to destination file
// you need check file exist before you call this function

func main() {
	/*
	   os.Mkdir("/home/ty4z2008/tar", 0777)
	   w, err := CopyFile("/home/ty4z2008/tar/1.pdf", "/home/ty4z2008/src/1.pdf")
	   //targetfile,sourcefile
	   if err != nil {
	       fmt.Println(err.Error())
	   }
	   fmt.Println(w)

	   TarGz("/home/ty4z2008/tar/1.pdf", "/home/ty4z2008/test.tar.gz") //压缩
	   //UnTarGz("/home/ty4z2008/1.tar.gz", "/home/ty4z2008")     //解压
	   os.RemoveAll("/home/ty4z2008/tar")
	*/
	//    TaZip("/Users/life/Desktop/j", "/Users/life/Desktop/aaa.tar.gz")
	Zip("/Users/life/Desktop/j", "/Users/life/Desktop/aaa.zip")
	fmt.Println("ok")
}

func TarGz(srcDirPath string, destFilePath string) (ok bool) {
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

	// Gzip writer
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// Tar writer
	tw := tar.NewWriter(gw)
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
		tarGzDir(srcDirPath, path.Base(srcDirPath), tw)
	} else {
		// handle file directly
		//        fmt.Println("Cerating tar.gz from " + fi.Name() + "...")
		tarGzFile(srcDirPath, fi.Name(), tw, fi)
	}
	ok = true
	return
}

// Deal with directories
// if find files, handle them with tarGzFile
// Every recurrence append the base path to the recPath
// recPath is the path inside of tar.gz
func tarGzDir(srcDirPath string, recPath string, tw *tar.Writer) {
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
			tarGzDir(curPath, recPath+"/"+fi.Name(), tw)
		} else {
			// File
			//            fmt.Printf("Adding file...%s\n", curPath)
		}

		tarGzFile(curPath, recPath+"/"+fi.Name(), tw, fi)
	}
}

// Deal with files
func tarGzFile(srcFile string, recPath string, tw *tar.Writer, fi os.FileInfo) {
	if fi.IsDir() {
		//    	fmt.Println("??")
		// Create tar header
		hdr := new(tar.Header)
		// if last character of header name is '/' it also can be directory
		// but if you don't set Typeflag, error will occur when you untargz
		hdr.Name = recPath // + "/"
		//        fmt.Println(hdr.Name)
		hdr.Typeflag = tar.TypeDir
		//        hdr.Size = 0
		//hdr.Mode = 0755 | c_ISDIR
		//        hdr.Mode = int64(fi.Mode()) // 加这个会有错误!!!
		//        hdr.ModTime = fi.ModTime() // 加这个会有错误!!

		// Write hander
		err := tw.WriteHeader(hdr)
		if err != nil {
			panic(err)
		}
	} else {
		// File reader
		fr, err := os.Open(srcFile)
		if err != nil {
			panic(err)
		}
		defer fr.Close()

		// Create tar header
		hdr := new(tar.Header)
		hdr.Name = recPath
		//        fmt.Println(hdr.Name)
		hdr.Size = fi.Size()
		hdr.Mode = int64(fi.Mode())
		hdr.ModTime = fi.ModTime()

		// Write hander
		err = tw.WriteHeader(hdr)
		if err != nil {
			panic(err)
		}

		// Write file data
		_, err = io.Copy(tw, fr)
		if err != nil {
			panic(err)
		}
	}
}

// Ungzip and untar from source file to destination directory
// you need check file exist before you call this function
func UnTarGz(srcFilePath string, destDirPath string) {
	//    fmt.Println("UnTarGzing " + srcFilePath + "...")
	// Create destination directory
	os.Mkdir(destDirPath, os.ModePerm)

	fr, err := os.Open(srcFilePath)
	if err != nil {
		panic(err)
	}
	defer fr.Close()

	// Gzip reader
	gr, err := gzip.NewReader(fr)
	if err != nil {
		panic(err)
	}
	defer gr.Close()

	// Tar reader
	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// End of tar archive
			break
		}
		//handleError(err)
		//        fmt.Println("UnTarGzing file..." + hdr.Name)
		// Check if it is diretory or file
		if hdr.Typeflag != tar.TypeDir {
			// Get files from archive
			// Create diretory before create file
			os.MkdirAll(destDirPath+"/"+path.Dir(hdr.Name), os.ModePerm)
			// Write data to file
			fw, _ := os.Create(destDirPath + "/" + hdr.Name)
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(fw, tr)
			if err != nil {
				panic(err)
			}
		}
	}
	//    fmt.Println("Well done!")
}

// Copyfile
func CopyFile(dstName, srcName string) (written int64, err error) {
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
