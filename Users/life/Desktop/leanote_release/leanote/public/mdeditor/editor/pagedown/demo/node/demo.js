// NOTE: This is just a demo -- in a production environment,
// be sure to spend a few more thoughts on sanitizing user input.
// (also, you probably wouldn't use a get request)

var http = require("http"),
    url = require("url"),
    querystring = require("querystring"),
    Converter = require("../../Markdown.Converter").Converter,
    getSanitizingConverter = require("../../Markdown.Sanitizer").getSanitizingConverter,
    conv = new Converter(),
    saneConv = getSanitizingConverter();

http.createServer(function (req, res) {

    var route = url.parse(req.url);
    if (route.pathname !== "/") {
        res.writeHead(404);
        res.end("Page not found");
        return;
    }
    
    var query = querystring.parse(route.query);

    res.writeHead(200, { "Content-type": "text/html" });
    res.write("<html><body>");
    
    var markdown = query.md || "## Hello!\n\n<marquee>I'm walking</marquee>\n\nVisit [Stack Overflow](http://stackoverflow.com)\n\n<b><i>This is never closed!";

    res.write("<h1>Your output, sanitized:</h1>\n" + saneConv.makeHtml(markdown))
    res.write("<h1>Your output, unsanitized:</h1>\n" + conv.makeHtml(markdown))
        
    res.write(
        "<h1>Enter Markdown</h1>\n" +
        "<form method='get' action='/'>" +
            "<textarea cols=50 rows=10 name='md'>" +
                markdown.replace(/</g, "&lt;") +
            "</textarea><br>" +
            "<input type='submit' value='Convert!'>" +
        "</form>"
    );
    
    res.end("</body></html>");

}).listen(8000);
