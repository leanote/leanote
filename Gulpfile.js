var gulp = require('gulp');
var clean = require('gulp-clean');
var uglify = require('gulp-uglify');
var rename = require('gulp-rename');
var minifyHtml = require("gulp-minify-html");
var concat = require('gulp-concat');
var replace = require('gulp-replace');
var inject = require('gulp-inject');
var gulpSequence = require('gulp-sequence');

var fs = require('fs');

var leanoteBase = './';
var base = leanoteBase + '/public'; // public base
var noteDev = leanoteBase + '/app/views/note/note-dev.html';
var noteProBase = leanoteBase + '/app/views/note';

// 合并Js, 这些js都是不怎么修改, 且是依赖
// 840kb, 非常耗时!!
gulp.task('concatDepJs', function() {
    var jss = [
        'js/jquery-1.9.0.min.js',
        'js/jquery.ztree.all-3.5-min.js',
        // 'tinymce/tinymce.full.min.js', // 使用打成的包, 加载速度快
        // 'libs/ace/ace.js',
        'js/jQuery-slimScroll-1.3.0/jquery.slimscroll-min.js',
        'js/contextmenu/jquery.contextmenu-min.js',
        'js/bootstrap-min.js',
        'js/object_id-min.js',
    ];

    for(var i in jss) {
        jss[i] = base + '/' + jss[i];
    }

    return gulp
        .src(jss)
        // .pipe(uglify()) // 压缩
        .pipe(concat('dep.min.js'))
        .pipe(gulp.dest(base + '/js'));
});

// 合并app js 这些js会经常变化 90kb
gulp.task('concatAppJs', function() {
    var jss = [
        'js/common.js',
        'js/app/note.js',
        'js/app/page.js', // 写作模式下, page依赖note
        'js/app/tag.js',
        'js/app/notebook.js',
        'js/app/share.js',
    ];

    for(var i in jss) {
        jss[i] = base + '/' + jss[i];
    }

    return gulp
        .src(jss)
        .pipe(uglify()) // 压缩
        .pipe(concat('app.min.js'))
        .pipe(gulp.dest(base + '/js'));
});

// plugins压缩
gulp.task('plugins', function() {
    // gulp.src(base + '/js/plugins/libs/*.js')
    //     .pipe(uglify()) // 压缩
    //     // .pipe(concat('main.min.js'))
    //     .pipe(gulp.dest(base + '/js/plugins/libs-min'));

    // 所有js合并成一个
     var jss = [
        'note_info',
        'tips',
        'history',
        'attachment_upload',
        'editor_drop_paste',
        'main'
    ];

    for(var i in jss) {
        jss[i] = base + '/js/plugins/' + jss[i] + '.js';
    }
    jss.push(base + '/js/plugins/libs-min/fileupload.js');

    return gulp.src(jss)
        .pipe(uglify()) // 压缩
        .pipe(concat('main.min.js'))
        .pipe(gulp.dest(base + '/js/plugins'));
});

// 合并requirejs和markdown为一个文件
gulp.task('concatMarkdownJs', function() {
    var jss = [
        'js/require.js',
        'md/main.min.js',
    ];

    for(var i in jss) {
        jss[i] = base + '/' + jss[i];
    }

    return gulp
        .src(jss)
        .pipe(uglify()) // 压缩
        .pipe(concat('markdown.min.js'))
        .pipe(gulp.dest(base + '/js'));
});

// / 合并requirejs和markdown为一个文件
gulp.task('concatMarkdownJsV2', function() {
    var jss = [
        'js/require.js',
        'md/main-v2.min.js',
    ];

    for(var i in jss) {
        jss[i] = base + '/' + jss[i];
    }

    return gulp
        .src(jss)
        .pipe(uglify()) // 压缩
        .pipe(concat('markdown-v2.min.js'))
        .pipe(gulp.dest(base + '/js'));
});

// note-dev.html -> note.html, 替换css, js
// TODO 加?t=2323232, 强制浏览器更新, 一般只需要把app.min.js上加
gulp.task('devToProHtml', function() {
    return gulp
        .src(noteDev)
        .pipe(replace(/<!-- dev -->[.\s\S]+?<!-- \/dev -->/g, '')) // 把dev 去掉
        .pipe(replace(/<!-- pro_dep_js -->/, '<script src="/js/dep.min.js"></script>')) // 替换
        .pipe(replace(/<!-- pro_app_js -->/, '<script src="/js/app.min.js"></script>')) // 替换
        // .pipe(replace(/<!-- pro_markdown_js -->/, '<script src="/js/markdown.min.js"></script>')) // 替换
        .pipe(replace(/<!-- pro_markdown_js -->/, '<script src="/js/markdown-v2.min.js"></script>')) // 替换
        .pipe(replace('/tinymce/tinymce.js', '/tinymce/tinymce.full.min.js')) // 替换
        .pipe(replace(/<!-- pro_tinymce_init_js -->/, "var tinyMCEPreInit = {base: '/public/tinymce', suffix: '.min'};")) // 替换
        .pipe(replace(/plugins\/main.js/, "plugins/main.min.js")) // 替换
        // 连续两个空行换成一个空行
        .pipe(replace(/\r\n\r\n/g, '\r\n'))
        .pipe(replace(/\r\n\r\n/g, '\r\n'))
        .pipe(replace(/\r\n\r\n/g, '\r\n'))
        .pipe(replace(/\r\n\r\n/g, '\r\n'))
        .pipe(replace(/\r\n\r\n/g, '\r\n'))
        .pipe(replace(/\r\n\r\n/g, '\r\n'))
        .pipe(replace('console.log(o);', ''))
        .pipe(replace('console.trace(o);', ''))
        // .pipe(minifyHtml()) // 不行, 压缩后golang报错
        .pipe(rename('note.html'))
        .pipe(gulp.dest(noteProBase));
});

// 只获取需要js i18n的key
var path = require('path');
gulp.task('i18n', function() {
    var keys = {};
    var reg = /getMsg\(["']+(.+?)["']+/g;
    // {rule: "required", msg: "inputNewPassword"},
    var reg2 = /msg: ?"?([0-9a-zA-Z]*)"?/g;
    function getKey(data) {
        while(ret = reg.exec(data)) {
            keys[ret[1]] = 1;
        }

        while(ret2 = reg2.exec(data)) {
            keys[ret2[1]] = 1;
        }
    }
    // 先获取需要的key
    function ls(ff) { 
        var files = fs.readdirSync(ff);  
        for(fn in files) {  
            var fname = ff + path.sep + files[fn];  
            var stat = fs.lstatSync(fname);  
            if(stat.isDirectory() == true) {
                ls(fname);
            } 
            else {
                if ((fname.indexOf('.html') > 0 || fname.indexOf('.js') > 0)) {
                    // console.log(fname);
                    // if (fname.indexOf('min.js') < 0) {
                        var data = fs.readFileSync(fname, "utf-8");
                        // 得到getMsg里的key
                        getKey(data);
                    // }
                }
            }  
        }  
    }

    console.log('parsing used keys');

    ls(base + '/admin');
    ls(base + '/blog');
    ls(base + '/md');
    ls(base + '/js');
    ls(base + '/album');
    ls(base + '/libs');
    ls(base + '/member');
    ls(base + '/tinymce');

    ls(leanoteBase + '/app/views');

    console.log('parsed');

    // msg.zh
    function getAllMsgs(fname) {
        var msg = {};

        var data = fs.readFileSync(fname, "utf-8");
        var lines = data.split('\n');
        for (var i = 0; i < lines.length; ++i) {
            var line = lines[i];
            // 忽略注释
            if (line[0] == '#' || line[1] == '#') {
                continue;
            }
            var lineArr = line.split('=');
            if (lineArr.length >= 2) {
               var key = lineArr[0];
               lineArr.shift();
               msg[key] = lineArr.join('=');
               // msg[lineArr[0]] = lineArr[1];
            }
        }
        return msg;
    }

    // msg.zh, msg.js
    function genI18nJsFile(fromFilename, otherNames, keys) {
        var msgs = {};
        otherNames.unshift(fromFilename);
        // console.log(fromFilename);
        // console.log(otherNames);
        otherNames.forEach(function (name) {
            var tmpMsgs = getAllMsgs(leanoteBase + '/messages/' + name);
            for (var i in tmpMsgs) {
                msgs[i] = tmpMsgs[i];
            }
        });

        var toMsgs = {};
        for (var i in msgs) {
            // 只要需要的
            if (i in keys) {
                toMsgs[i] = msgs[i];
            }
        }
        var str = 'var MSG=' + JSON.stringify(toMsgs) + ';';
        str += 'function getMsg(key, data) {var msg = MSG[key];if(msg) {if(data) {if(!isArray(data)) {data = [data];}' + 
                        'for(var i = 0; i < data.length; ++i) {' + 
                            'msg = msg.replace("%s", data[i]);' + 
                        '}' + 
                    '}' + 
                    'return msg;' + 
                '}' + 
                'return key;' + 
            '}';

        // 写入到文件中
        var toFilename = fromFilename + '.js';
        fs.writeFile(base + '/js/i18n/' + toFilename, str);
    }

    // 必须要的
    // keys.push();

    genI18nJsFile('blog.zh', [], keys);
    genI18nJsFile('blog.en', [], keys);
    genI18nJsFile('blog.fr', [], keys);
    genI18nJsFile('blog.pt', [], keys);

    genI18nJsFile('msg.fr', ['member.fr', 'markdown.fr', 'album.fr'], keys);
    genI18nJsFile('msg.zh', ['member.zh', 'markdown.zh', 'album.zh'], keys);
    genI18nJsFile('msg.en', ['member.en', 'markdown.en', 'album.en'], keys);
    genI18nJsFile('msg.pt', ['member.pt', 'markdown.pt', 'album.pt'], keys);
});

// 合并album需要的js
gulp.task('concatAlbumJs', function() {
    /*
    gulp.src(base + '/album/js/main.js')
        .pipe(uglify()) // 压缩
        .pipe(rename({suffix: '.min'}))
        .pipe(gulp.dest(base + '/album/js/'));
    */

    gulp.src(base + '/album/css/style.css')
        .pipe(rename({suffix: '-min'}))
        .pipe(minifycss())
        .pipe(gulp.dest(base + '/album/css'));

    var jss = [
        'js/jquery-1.9.0.min.js',
        'js/bootstrap-min.js',
        'js/plugins/libs-min/fileupload.js',
        'js/jquery.pagination.js',
        'album/js/main.js',
    ];

    for(var i in jss) {
        jss[i] = base + '/' + jss[i];
    }

    return gulp
        .src(jss)
        .pipe(uglify()) // 压缩
        .pipe(concat('main.all.js'))
        .pipe(gulp.dest(base + '/album/js'));
});

// tinymce
// please set the right path on your own env
var tinymceBase = '/Users/life/leanote/leanote-tools/tinymce_4.1.9_leanote_public';
gulp.task('tinymce', function() {
    // 先清理
    fs.unlink(tinymceBase + '/js/tinymce/tinymce.dev.js');
    fs.unlink(tinymceBase + '/js/tinymce/tinymce.jquery.dev.js');
    fs.unlink(tinymceBase + '/js/tinymce/tinymce.full.js');
    fs.unlink(tinymceBase + '/js/tinymce/tinymce.full.min.js');

    var cp = require('child_process');

    var bundleCmd = 'grunt bundle --themes leanote --plugins autolink,link,leaui_image,lists,hr,paste,searchreplace,leanote_nav,leanote_code,tabfocus,table,directionality,textcolor';
    // build
    cp.exec('grunt minify', {cwd: tinymceBase}, function(err, stdout, stderr) {
        console.log('stdout: ' + stdout);
        console.log('stderr: ' + stderr);

        // 将所有都合并成一起
        cp.exec(bundleCmd, {cwd: tinymceBase}, function(err, stdout, stderr) {
            console.log('stdout: ' + stdout);
            console.log('stderr: ' + stderr);
        });
    });
});

// 合并css, 无用
// Deprecated
gulp.task('concatCss', function() {
    return gulp
        .src([markdownRaw + '/css/default.css', markdownRaw + '/css/md.css'])
        .pipe(concat('all.css'))
        .pipe(gulp.dest(markdownMin));
});

// mincss
var minifycss = require('gulp-minify-css');
gulp.task('minifycss', function() {
    gulp.src(base + '/css/bootstrap.css')
        .pipe(rename({suffix: '-min'}))
        .pipe(minifycss())
        .pipe(gulp.dest(base + '/css'));

    gulp.src(base + '/css/font-awesome-4.2.0/css/font-awesome.css')
        .pipe(rename({suffix: '-min'}))
        .pipe(minifycss())
        .pipe(gulp.dest(base + '/css/font-awesome-4.2.0/css'));

    gulp.src(base + '/css/zTreeStyle/zTreeStyle.css')
        .pipe(rename({suffix: '-min'}))
        .pipe(minifycss())
        .pipe(gulp.dest(base + '/css/zTreeStyle'));

    gulp.src(base + '/md/themes/default.css')
        .pipe(rename({suffix: '-min'}))
        .pipe(minifycss())
        .pipe(gulp.dest(base + '/md/themes'));

    gulp.src(base + '/js/contextmenu/css/contextmenu.css')
        .pipe(rename({suffix: '-min'}))
        .pipe(minifycss())
        .pipe(gulp.dest(base + '/js/contextmenu/css'));

    // theme
    // 用codekit
    var as = ['default', 'simple', 'writting', /*'writting-overwrite', */ 'mobile'];
    /*
    for(var i = 0; i < as.length; ++i) {
        gulp.src(base + '/css/theme/' + as[i] + '.css')
            .pipe(minifycss())
            .pipe(gulp.dest(base + '/css/theme'));
    }
    */
});


gulp.task('concat', ['concatDepJs', 'concatAppJs', /* 'concatMarkdownJs', */'concatMarkdownJsV2']);
gulp.task('html', ['devToProHtml']);
gulp.task('default', ['concat', 'plugins', 'minifycss', 'i18n', 'concatAlbumJs', 'html']);
