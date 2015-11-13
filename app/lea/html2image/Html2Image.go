package html2image

/*
import (
	"github.com/leanote/leanote/app/lea"
	"github.com/leanote/leanote/app/lea/netutil"
	"bufio"
	"code.google.com/p/draw2d/draw2d"
//	"fmt"
	"image"
	"image/color"
	"image/png"
	"image/gif"
	"image/jpeg"
	"os"
//	"strings"
//	"time"
	"github.com/revel/revel"
	"code.google.com/p/go.net/html"
	"strings"
	"strconv"
)

type Html2Image struct {
	image *image.RGBA
	gc *draw2d.ImageGraphicContext

	// 试探
	gc2 *draw2d.ImageGraphicContext

	width float64 // 图片宽度
	height float64

	painWidth float64 // 画布宽度

	startX float64
	x float64
	y float64

	isFirstP bool // 是否是第一个段落?

	// 换行和段落的高度
	brY float64
	pY float64

	// 字体
	normalFontFamily draw2d.FontData
	boldFontFamily draw2d.FontData

	// preTag 之前的标签
	preTag *html.Node
}

func NewHtml2Image() *Html2Image {
	h := &Html2Image{}
	h.width = 440
	h.height = 10000 // 最开始设为很大, 不然加不了图片, 会影响速度
	i, gc := h.InitGc(h.width, h.height)
	h.gc = gc;
	h.image = i

	// 试探
	_, h.gc2 = h.InitGc(h.width, 100)

	h.startX = 10

	// 最初位置
	h.x =  h.startX
	h.y = 80

	h.isFirstP = true

	h.normalFontFamily = draw2d.FontData{"xihei", 4, draw2d.FontStyleNormal};
	h.boldFontFamily = draw2d.FontData{"heiti", 5, draw2d.FontStyleNormal};

	h.SetNormalFont()

	return h
}

func (this *Html2Image) InitGc(w, h float64) (* image.RGBA, *draw2d.ImageGraphicContext) {
	i := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	gc := draw2d.NewGraphicContext(i)

	gc.SetStrokeColor(image.Black)
	gc.SetFillColor(image.White)
	// fill the background
	// gc.Clear()

	draw2d.SetFontFolder(revel.BasePath + "/public/fonts/weibo")
	draw2d.Rect(gc, 0, 0, w, h) // 设置背景
	gc.FillStroke()
	gc.SetFillColor(image.Black)

	// 这个很耗时
//	gc.Translate(0, 0)
	return i, gc
}

func (this *Html2Image) SaveToPngFile(filePath string) bool {
//	m := this.image;
	m := this.image.SubImage(image.Rect(0, 0, int(this.width), int(this.y + 20)))
	// 需要截断之

	f, err := os.Create(filePath)
	if err != nil {
		return false
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	err = png.Encode(b, m)
	if err != nil {
		return false
	}
	err = b.Flush()
	if err != nil {
		return false
	}
	return true
}

// 字体大小

func (this *Html2Image) SetSmallFont() {
	this.gc.SetFontData(this.normalFontFamily)
	this.gc2.SetFontData(this.normalFontFamily)

	this.gc.SetFillColor(color.NRGBA{60, 60, 60, 255})
	this.gc.SetFontSize(12)
	this.gc2.SetFontSize(12)

	this.brY = 16
	this.pY = 30

	this.painWidth = this.width - 10
}

func (this *Html2Image) SetNormalFont() {
	this.gc.SetFillColor(image.Black)
	this.gc.SetFontData(this.normalFontFamily)
	this.gc2.SetFontData(this.normalFontFamily)

	this.gc.SetFontSize(14)
	this.gc2.SetFontSize(14)

	this.brY = 20
	this.pY = 30

	this.painWidth = this.width - 10
}
func (this *Html2Image) SetAColor() {
	this.gc.SetFillColor(color.NRGBA{66, 139, 202, 255})
}

// 标题
func (this *Html2Image) SetTitleFont() {
	this.gc.SetFillColor(image.Black)

	this.gc.SetFontData(this.boldFontFamily)
	this.gc2.SetFontData(this.boldFontFamily)

	this.gc.SetFontSize(24)
	this.gc2.SetFontSize(24)

	this.brY = 30
	this.pY = 60

	this.painWidth = this.width - 100
}

// h1, h2...
// h1
func (this *Html2Image) SetHeadFont(h string) {
	this.gc.SetFillColor(image.Black)

	this.gc.SetFontData(this.boldFontFamily)
	this.gc2.SetFontData(this.boldFontFamily)

	this.painWidth = this.width - 50

	if h == "h1" {
		this.gc.SetFontSize(20)
		this.gc2.SetFontSize(20)
		this.brY = 30
		this.pY = 60
	} else if h == "h2" {
		this.gc.SetFontSize(16)
		this.gc2.SetFontSize(16)
		this.brY = 25
		this.pY = 50
	} else if h == "h3" || h == "h4" {
		this.gc.SetFontSize(14)
		this.gc2.SetFontSize(14)
		this.brY = 20
		this.pY = 40
	}
}

// 新的段落
func (this *Html2Image) NewP() {
//	if !this.isFirstP {
		this.x = this.startX;
		this.y += this.pY;
//	} else {
//		this.isFirstP = false
//	}
}

// 新一行
func (this *Html2Image) NewBr() {
	this.x = this.startX;
	this.y += this.brY;
}

// 是否超出了
// 1 个汉字 17.875
// 1 个字母最多 17.15625
func (this *Html2Image) IsOver(r []rune) bool {
//	fmt.Println(string(r))
	// 还有多宽
	// 就算全拿汉字来说
	// 这里是优化，速度有提升
//	if this.painWidth - this.x > 17.875 * float64(len(r)) {
//		return false
//	}

//	println(text)
	width2 := this.gc2.FillStringAt(string(r), 0, 0)
	// 以下的方法可以极大节约时间
//	a, b, c, d := this.gc2.GetStringBounds(string(r))
//	width2 := c - a + 2

//	fmt.Println(width2)
//	fmt.Println(c - a)

	// 小于, 那么需要->大 到第一个不合适的位置
	if width2 + this.x <= this.painWidth {
		return false
	}
	return true
}
// 是否是字母
func (this *Html2Image) isAlpha(word rune) bool {
	if (word >= 'a' && word <= 'z') || (word >= 'A' && word <= 'Z') {
		return true
	}
	return false
}

// 是否是标点, 是标点就包含进来
func (this *Html2Image) includePunctuation(r []rune, end int) int {
	if len(r) == end {
		return end
	}
	c := r[end]
	if c == ',' || c == '.' || c == '?' || c == ':' || c == ';' || c == '!' || c == '，' || c == '。' || c == '？' || c == '：' || c == '；' || c == '！' {
		return end + 1
	}
	return end;
}

// 插入文本
// 要判断是否太长了, 太长了就截断
func (this *Html2Image) InsertText(text string, needTest bool, prefix string) {
	if needTest && this.painWidth - this.x < 2 {
		// 另起一行
		this.NewBr()
		this.InsertText(text, true, prefix)
		return;
	}

	r := []rune(text)
	// 试探吧, 可能需要截取
	if !needTest || !this.IsOver(r) {
		// 不用截
		// 可能有\n
		width := this.gc.FillStringAt(prefix + text, this.x, this.y)
		this.x += width + 1
	} else {
		// 刚开始加10个字, 之后一个一个来
		// 一个汉字, 或一个单词加
		wordStart := false
		wordStartPos := 0
		maxRI := len(r) - 1
		for i, word := range r { // i 是0, 1, 2, 3...
			// i是byte的位置, 一个汉字占3位
			// 是字母
			if this.isAlpha(word) {
				if !wordStart {
					wordStart = true
					wordStartPos = i
				}
			} else if(word == '\n' || word == '\r') {
				// 是否是\n
				i = this.includePunctuation(r, i)
				this.InsertText(string(r[0:i]), false, prefix)
				this.NewBr()
				// 之后的
				if maxRI != i {
					this.InsertText(string(r[i+1:]), true, prefix)
				}
				return;
			} else {
				// 单词没结束不计算
				wordStart = false
				if i > 0 {
					// 这里计算是否超出了, 包含自己在内
					if this.IsOver(r[0:i+1]) {
						// 那么...回退前一个
						end := i
						// 如果上一个是单词, 那么整个单词都不要, 取单词开头
						if this.isAlpha(r[i-1]) {
							end = wordStartPos
							// 这一行全是这个单词, 不太现实, 但有可能, 只能截断了
							if end == 0 {
								end = i
							}
						}

						// 这一段写上
//						println("------>" + string(r[0:end]))

						// 这里, 判断后面一个是否是标点符号
						end = this.includePunctuation(r, end)
						this.InsertText(string(r[0:end]), false, prefix)
						this.NewBr()
						// 之后的
						this.InsertText(string(r[end:]), true, prefix)

						return;
					} else {
						// 没超出, 不用计算, 但出要看是否是结尾了
						// 怎么可能会出现这种情况呢?, 第一步就试了
//						if i+1 == len(text) {
//							println("不可能")
//						}
					}
				}
			}
		} // for
		// !!
		// 如果是 go get code.google.com/p/graphics-go/graphics 最后是字母, 怎么办?
		if wordStart {
			// 这里计算是否超出了, 包含自己在内
			end := maxRI + 1
			// 如果上一个是单词, 那么整个单词都不要, 取单词开头
			if this.isAlpha(r[maxRI]) {
				end = wordStartPos
				// 这一行全是这个单词, 不太现实, 但有可能, 只能截断了
				if end == 0 {
					end = maxRI + 1
				}
			}

			// 这一段写上
//			println("-e----->" + string(r[0:end]))

			// 这里, 判断后面一个是否是标点符号
			end = this.includePunctuation(r, end)
			this.InsertText(string(r[0:end]), false, prefix)
			this.NewBr()
			// 之后的
			this.InsertText(string(r[end:]), true, prefix)

			return;
		}
	}
}

// 设置页脚, url文章链接
func (this *Html2Image) SetBottom(username, url string) {
	// 画一条线
	this.NewBr()
	this.gc.MoveTo(this.x, this.y)
    this.gc.LineTo(this.painWidth, this.y)
	this.gc.SetStrokeColor(color.NRGBA{200, 0, 0, 255})
	this.gc.SetLineWidth(2)
    this.gc.FillStroke()

    this.SetSmallFont()

    // 左侧写字
    this.NewP()
    this.InsertText("本文来自 " + username + " 的leanote笔记", true, "  ")
    this.NewBr()
    this.InsertText("个人博客: ", false, "  ")
    siteUrl, _ := revel.Config.String("site.url")
    if siteUrl == "" {
    	siteUrl = "http://leanote.com"
    }
    this.InsertA(siteUrl + "/blog/" + username, false)

	this.setLogo()
//    this.painWidth = this.width - 100
//    this.NewP()
//    this.InsertText("leanote, 不一样的笔记.", false, "  ")
//    this.NewBr()
//    this.InsertText("在这里你可以管理自己的知识", false, "  ")
//    this.NewBr()
//    this.InsertText("将知识分享给好友, 与好友一起协作知识", false, "  ")
//    this.NewBr()
//    this.InsertText("并且还可以将笔记设为博客公开", false, "  ")
//    this.InsertText(". 赶紧加入吧! leanote.com", false, "")
//
    // Logo
}

func (this *Html2Image) setImage(path string, x, y float64) {
    f1, err := os.Open(path)
    if err != nil {
    	return;
        panic(err)
    }

    var m1 image.Image
	_, ext := lea.SplitFilename(path)
    if ext == ".png" {
	    m1, err = png.Decode(f1)
    } else if ext == ".gif" {
	    m1, err = gif.Decode(f1)
    } else if ext == ".jpg" {
	    m1, err = jpeg.Decode(f1)
    }
	if err != nil {
	    return
        panic(err)
    }

    this.gc.Translate(x, y)
    this.gc.DrawImage(m1)
    this.gc.Translate(-x, -y)
}

// 画leanote logo
func (this *Html2Image) setLogo() {
	// 右上角的logo
	path := revel.BasePath + "/public/images/leanote/logo-20-a-6.png"
	println(path)
	this.setImage(path, 320, 10)

	// 右下角设置Logo
//	path = revel.BasePath + "/public/images/leanote/logo-60-a-6.png"
//	this.setImage(path, 320, this.y - 75)
}

// 插入链接
func (this *Html2Image) InsertA(text string, isNormal bool) {
	if text == "" {
		return
	}

	this.SetAColor()
	this.InsertText(text, true, "")

	// 还原
	if isNormal {
		this.SetNormalFont()
	} else {
		this.SetSmallFont()
	}
}

// 文章标题
func (this *Html2Image) InsertTitle(title string) {
	oldX := this.x
	oldY := this.y - 35

	// 插入之
	this.SetTitleFont()

	this.InsertText(title, true, "  ")

	// 还原字体大小
	this.SetNormalFont()

	this.NewBr()

	this.gc.MoveTo(oldX, oldY)
    this.gc.LineTo(this.x, this.y - 10)
	this.gc.SetStrokeColor(color.NRGBA{200, 0, 0, 255})
	this.gc.SetLineWidth(5)
    this.gc.FillStroke()
}

// 插入h1, h2, ... h4
func (this *Html2Image) InsertHead(n *html.Node) {
	this.SetHeadFont(n.Data); // h1, h2...
	this.NewP()
	// 把标题内容全都拿出
	var text = ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.Data == "p" || c.Data == "br"){
			text += " "
		} else {
			text += strings.TrimRight(c.Data, "\n")
		}
	}
	this.InsertText(text, true, "")
	this.SetNormalFont()
}

// 插入代码
func (this *Html2Image) InsertCode(n *html.Node) {
	this.NewP()
	oldX := this.x
	oldY := this.y - 20
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.Data == "p" || c.Data == "br"){
			this.NewBr()
		} else {
			this.InsertText(strings.TrimRight(c.Data, "\n"), true, "  ")
		}
	}
	this.NewBr()

	this.gc.MoveTo(oldX, oldY)
    this.gc.LineTo(this.x, this.y - 20)
	this.gc.SetStrokeColor(color.NRGBA{0, 200, 0, 255})
	this.gc.SetLineWidth(2)
    this.gc.FillStroke()
}

// 插入图片
// 这个path应该是url,
// http://abc.com/a.gif 需要先下载
// 或 /upload/a.gif
func (this *Html2Image) InsertImage(path string, needTrans bool, width uint) {
	if path == "" {
		return;
	}

	// 是url, 那么取网络图片之
	var ok bool
	if strings.HasPrefix(path, "http") || strings.HasPrefix(path, "//") {
		path, ok = netutil.WriteUrl(path, "/tmp")
		if !ok || path == ""{
			return
		}
	} else {
		path = revel.BasePath + "/public/" + path
	}

	// 需要转换, logo不需要转换
	if(needTrans) {
		painWidth := uint(this.painWidth - 10)
		if width > 0 && painWidth > width {
			painWidth = width
		}
		ok, path = lea.TransToGif(path, painWidth, false)
		if !ok || path == "" {
			return;
		}
	}

    f1, err := os.Open(path)
    if err != nil {
    	return;
        panic(err)
    }

    var m1 image.Image
	_, ext := lea.SplitFilename(path)
    if ext == ".png" {
	    m1, err = png.Decode(f1)
    } else {
	    m1, err = gif.Decode(f1)
    }
	if err != nil {
	    return
        panic(err)
    }

    // 如果之前是p, 那么不要有<br>
    if this.preTag.Data != "p" {
	    this.NewBr()
    }
    this.gc.Translate(this.x, this.y)
    this.gc.DrawImage(m1)
    // 还原
    this.gc.Translate(-this.x, -this.y) // 这个有用些
	this.y += float64(m1.Bounds().Dy()) - 20
	this.NewP()

    os.Remove(path)

    // 如果图片是文章第一个的话, 之后的需要p
    this.isFirstP = false
}

// 内容主体
func (this *Html2Image) InsertBody(htmlStr string) (ok bool) {
	reader := bufio.NewReader(strings.NewReader(htmlStr))
	doc, err := html.Parse(reader)
	if err != nil {
		return;
	}
	var f func(*html.Node, *html.Node, string)
	f = func(n *html.Node, p *html.Node, prefix string) {
//		if p != nil {
			// fmt.Println("Parent Data: " + p.Data)
//		}

		defer func() {
			if n.Type == html.ElementNode {
				this.preTag = n
			}
		}()

		// 标签
		if n.Type == html.ElementNode {
			if n.Data == "p" {
				this.NewP()
				// 遍历之后的
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c, n, prefix)
				}
				return;
			}

			// 也是一个段落, 只是要缩进
			if n.Data == "ul" || n.Data == "ol" {
				this.NewP()
				// 遍历之后的
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c, n, "")
				}
				return;
			}
			if n.Data == "li" {
				// 遍历之后的
				// 是否需要前缀
				needPrefix := true // 第一个肯定要
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if needPrefix {
						f(c, n, "     ")
						needPrefix = false
					} else {
						f(c, n, "")
					}

					if c.Type == html.ElementNode {
						if c.Data == "br" || c.Data == "p" {
							needPrefix = true
						} else {
							needPrefix = false
						}
					}
				}
				// 输完一行后再换行
				this.NewBr()
				return;
			}

			// 标题
			if n.Data == "h1" || n.Data == "h2" || n.Data == "h3" || n.Data == "h4" {
				this.InsertHead(n)
				return;
			}

			if n.Data == "pre" {
				// 把之后的全拿过来
				this.InsertCode(n)
				return;
			}

			// 图片
			// 得到src
			if n.Data == "img" {
				src := ""
				width := 0
				for _, attr := range n.Attr {
					if attr.Key == "src" {
						src = attr.Val
					} else if attr.Key == "width" {
						width, _ = strconv.Atoi(attr.Val)
					}
				}
				if src != "" {
					this.InsertImage(src, true, uint(width))
				}
				return;
			}

			// 链接
			// 如果链接里只有文本, 那么单独处理, 如果还有其它的, 不作链接处理
			if n.Data == "a" {
				if n.FirstChild == n.LastChild {
					this.InsertA(n.FirstChild.Data, true)
					return;
				}
			}

			// 空行
			if n.Data == "br" { // || n.Data == "div"
				this.NewBr()
			}
		}

		// 是文本, 输出之
		if n.Type == html.TextNode {
			data := strings.TrimSpace(n.Data);
			// <p>xx<br/>xxx</p>    这些空白也是TextNode    <p>
			if data != "" {
				this.InsertText(prefix + data, true, "")
			}
			return;
		}

		// 其余的

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c, n, prefix)
		}

		return;
	}
	f(doc, nil, "")

	return true
}

// 主函数
func ToImage(uid, username, noteId, title, htmlStr, toPath string) (ok bool) {
	h := NewHtml2Image()

	// 标题
	h.InsertTitle(title)

	// 主体
	ok = h.InsertBody(htmlStr)
	if(!ok) {
		return
	}

	// 页眉与页脚
	h.SetBottom(username, "")

	// 保存成png图片
	ok = h.SaveToPngFile(toPath)
	return
}

func TestFillString() {
	h := NewHtml2Image()
	str := `一个合格的 Techspace 需要有足够专业的器材、场地和资源，你可以和你的团队在里面进行激光切割、快速贴片甚至加工木材等操作，在相对独立的空间内又能同周围的同道友人互相激发切磋。国内现有的 Techspace 没几家，不久前我去深圳特地拜访了当地的 Techspace，很喜欢那里的氛围，希望国内其他地方也能有更多这类空间供创客发挥。
假如你有一个比较成型的想法，想在硬件领域做点事情，核心团队也基本组好，硬件软件交互基本都有专人了`
//	h.IsOver("W")
	h.InsertText("go get code.google.com/p/graphics\n-go/graphics", true, "")
//	h.InsertText("usr/bin/install: 无法创建一般文件'/usr/local/jpeg6/include/jconfig.", true)
//	h.InsertImage("/Users/life/Desktop/share.png")
//	h.NewP()
	h.InsertText(str, true, "")
//	h.InsertImage("/Users/life/Desktop/share.png")
	h.SaveToPngFile("/Users/life/Desktop/TestPath3.png")
}
*/

func ToImage(uid, username, noteId, title, htmlStr, toPath string) (ok bool) {
	return false
}
