package lea

/*
import (
	"github.com/Terry-Mao/paint"
    "github.com/Terry-Mao/paint/wand"
    "fmt"
    "os"
)

// 传源路径, 在该路径下写入另一个gif
// maxWidth 最大宽度, == 0表示不改变宽度
// 成功后删除
func TransToGif(path string, maxWidth uint, afterDelete bool) (ok bool, transPath string) {
	ok = false
	transPath = path
	wand.Genesis()
    defer wand.Terminus()

    w := wand.NewMagickWand()
    defer w.Destroy()

    if err := w.ReadImage(path); err != nil {
    	fmt.Println(err);
    	return;
    }

    width := w.ImageWidth()
    height := w.ImageHeight()
    if maxWidth != 0 {
	    if width > maxWidth {
	    	// 等比缩放
	    	height = height * maxWidth/width
	    	width = maxWidth
	    }
    }

	w.SetImageFormat("GIF");

    if err := paint.Thumbnail(w, width, height); err != nil {
    	fmt.Println(err);
    	return;
    }

    // 判断是否是gif图片, 是就不用转换了
	baseName, ext := SplitFilename(path)
    var toPath string
    if(ext == ".gif") {
	    toPath = baseName + "_2" + ext
    } else {
	    toPath = TransferExt(path, ".gif");
    }

    if err := w.WriteImage(toPath); err != nil {
    	fmt.Println(err);
    	return;
    }

    if afterDelete {
    	os.Remove(path)
    }

    ok = true
    transPath = toPath

    return
}

// 缩小到
// 未用
func Reset(path string, maxWidth uint) (ok bool, transPath string){
	wand.Genesis()
    defer wand.Terminus()

    w := wand.NewMagickWand()
    defer w.Destroy()

    if err := w.ReadImage(path); err != nil {
    	fmt.Println(err);
    	return;
    }

    width := w.ImageWidth()
    height := w.ImageHeight()
    if maxWidth != 0 {
	    if width > maxWidth {
	    	// 等比缩放
	    	height = height * maxWidth/width
	    	width = maxWidth
	    }
    }
    if err := paint.Thumbnail(w, width, height); err != nil {
    	fmt.Println(err);
    	return;
    }

    toPath := TransferExt(path, ".gif");
    if err := w.WriteImage(toPath); err != nil {
    	fmt.Println(err);
    	return;
    }

    ok = true
    transPath = toPath

    return
}
*/

func TransToGif(path string, maxWidth uint, afterDelete bool) (ok bool, transPath string) {
	return ok, path
}
