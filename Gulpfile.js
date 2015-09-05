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
        'tinymce/tinymce.full.min.js', // 使用打成的包, 加载速度快
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

// 合并requirejs和markdown为一个文件
gulp.task('concatMarkdownJs', function() {
    var jss = [
        'js/require.js',
        'dist/main.min.js',
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

// note-dev.html -> note.html, 替换css, js
// TODO 加?t=2323232, 强制浏览器更新, 一般只需要把app.min.js上加
gulp.task('devToProHtml', function() {
    return gulp
        .src(noteDev)
        .pipe(replace(/<!-- dev -->[.\s\S]+?<!-- \/dev -->/g, '')) // 把dev 去掉
        .pipe(replace(/<!-- pro_dep_js -->/, '<script src="/js/dep.min.js"></script>')) // 替换
        .pipe(replace(/<!-- pro_app_js -->/, '<script src="/js/app.min.js"></script>')) // 替换
        .pipe(replace(/<!-- pro_markdown_js -->/, '<script src="/js/markdown.min.js"></script>')) // 替换
        .pipe(replace(/<!-- pro_tinymce_init_js -->/, "var tinyMCEPreInit = {base: '/public/tinymce', suffix: '.min'};")) // 替换
        // 连续两个空行换成一个空行, 没用
        .pipe(replace(/\n\n/g, '\n'))
        .pipe(replace(/\n\n/g, '\n'))
        .pipe(replace('console.log(o);', ''))
        // .pipe(minifyHtml()) // 不行, 压缩后golang报错
        .pipe(rename('note.html'))
        .pipe(gulp.dest(noteProBase));
});

// tinymce

var tinymceBase = base + '/tinymce_4.1.9';
gulp.task('tinymce', function() {
    // 先清理
    fs.unlink(tinymceBase + '/js/tinymce/tinymce.dev.js');
    fs.unlink(tinymceBase + '/js/tinymce/tinymce.jquery.dev.js');
    fs.unlink(tinymceBase + '/js/tinymce/tinymce.full.js');
    fs.unlink(tinymceBase + '/js/tinymce/tinymce.full.min.js');

    var cp = require('child_process');

    var bundleCmd = 'grunt bundle --themes leanote --plugins autolink,link,leaui_image,leaui_mind,lists,hr,paste,searchreplace,leanote_nav,leanote_code,tabfocus,table,directionality,textcolor';
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
gulp.task('concatCss', function() {
    return gulp
        .src([markdownRaw + '/css/default.css', markdownRaw + '/css/md.css'])
        .pipe(concat('all.css'))
        .pipe(gulp.dest(markdownMin));
});

gulp.task('concat', ['concatDepJs', 'concatAppJs', 'concatMarkdownJs']);
gulp.task('html', ['devToProHtml']);
gulp.task('default', ['concat', 'html']);
