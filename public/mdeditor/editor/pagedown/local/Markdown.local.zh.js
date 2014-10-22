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
        linkTitle:"超链接",
        linkAddress:"请输入链接地址",
        linkdescription: "此处输入链接的描述",
        linkExample: 'http://example.com/ "可选标题"',
       
        image: "图片 <img> Ctrl+G",
        imageTitle:"图片",
        imageAddress:"请输入图片地址",
        imagedescription: "此处输入图片的描述",
        imageExample: 'http://example.com/images/diagram.jpg "可选标题"',
        
        quote: "段落引用 <blockquote> Ctrl+Q",
        quoteexample: "段落引用",

        code: "代码样例 <pre><code> Ctrl+K",
        codeexample: "此处输入代码",

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
