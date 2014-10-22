// Usage:
//
// var myConverter = new Markdown.Editor(myConverter, null, { strings: Markdown.local.fr });

(function () {
        Markdown.local = Markdown.local || {};
        Markdown.local.zh = {
        bold: "Bold <strong> Ctrl+B",
        boldexample: "Bold",

        italic: "Italic <em> Ctrl+I",
        italicexample: "Italic",

		link: "Link <a> Ctrl+L",
        linkTitle:"Link",
        linkAddress:"Input Link",
        linkExample: 'http://example.com/ "Optional Title"',
       
        image: "Image <img> Ctrl+G",
        imageTitle:"Image",
        imageAddress:"Input Image Address",
        imageExample: 'http://example.com/images/diagram.jpg "Optional Title"',

        quote: "Blockquote <blockquote> Ctrl+Q",
        quoteexample: "Blockquote ",

        code: "Code <pre><code> Ctrl+K",
        codeexample: "Input Code",

        olist: "Ordered List <ol> Ctrl+O",
        ulist: "Unordered list <ul> Ctrl+U",
        litem: "item",

        heading: "Heading <h1>/<h2> Ctrl+H",
        headingexample: "Heading",

        hr: "Hr <hr> Ctrl+R",

        undo: "Undo - Ctrl+Z",
        redo: "Redo - Ctrl+Y",
        redomac: "Redomac - Ctrl+Shift+Z",

        help: "Markdown Help Ctrl+Alt+H"
    };
})();
