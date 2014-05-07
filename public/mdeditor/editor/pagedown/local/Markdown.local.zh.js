// Usage:
//
// var myConverter = new Markdown.Editor(myConverter, null, { strings: Markdown.local.fr });

(function () {
        Markdown.local = Markdown.local || {};
        Markdown.local.zh = {
        bold: "粗体 <strong> Ctrl+B",
        boldexample: "粗体文本",

        italic: "斜体 <em> Ctrl+I",
        italicexample: "斜体文本",

        link: "超链接 <a> Ctrl+L",
        linkdescription: "此处输入链接的描述",
        linkdialog: "<p><b>输入超链接</b></p><p>http://example.com/ \"可选标题\"</p>",

        quote: "段落引用 <blockquote> Ctrl+Q",
        quoteexample: "段落引用",

        code: "代码样例 <pre><code> Ctrl+K",
        codeexample: "此处输入代码",

        image: "图片 <img> Ctrl+G",
        imagedescription: "此处输入图片的描述",
        imagedialog: "<p><b>插入图片</b></p><p>http://example.com/images/diagram.jpg \"可选标题\"<br><br><a href='http://www.google.com/search?q=free+image+hosting' target='_blank'>需要免费的图片主机？</a></p>",

        olist: "有序列表 <ol> Ctrl+O",
        ulist: "无序列表 <ul> Ctrl+U",
        litem: "列表项",

        heading: "标题 <h1>/<h2> Ctrl+H",
        headingexample: "标题",

        hr: "水平线 <hr> Ctrl+R",

        undo: "撤销 - Ctrl+Z",
        redo: "重复 - Ctrl+Y",
        redomac: "重复 - Ctrl+Shift+Z",

        help: "Markdown 语法帮助 Ctrl+Alt+H"
    };
})();
