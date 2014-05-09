new function($) {
  $.fn.setCursorPosition = function(pos) {
    if ($(this).get(0).setSelectionRange) {
      $(this).get(0).setSelectionRange(pos, pos);
    } else if ($(this).get(0).createTextRange) {
      var range = $(this).get(0).createTextRange();
      range.collapse(true);
      range.moveEnd('character', pos);
      range.moveStart('character', pos);
      range.select();
    }
    $(this).focus();
  }
  $.fn.tabHandler = function() {
    $(this).keydown(function(e) {
        if(e.keyCode === 9) { // tab was pressed
            // get caret position/selection
            var start = this.selectionStart;
            var end = this.selectionEnd;

            var $this = $(this);
            var value = $this.val();

            // set textarea value to: text before caret + four spaces + text after caret
            $this.val(value.substring(0, start)
                        + "    "
                        + value.substring(end));

            // put caret at right position again (add four for the spaces)
            this.selectionStart = this.selectionEnd = start + 4;

            // prevent the focus lose
            e.preventDefault();
        }
    });
  }
}(jQuery);

// full screen api
(function() {
    var
        fullScreenApi = {
            supportsFullScreen: false,
            isFullScreen: function() { return false; },
            requestFullScreen: function() {},
            cancelFullScreen: function() {},
            fullScreenEventName: '',
            prefix: ''
        },
        browserPrefixes = 'webkit moz o ms khtml'.split(' ');
 
    // check for native support
    if (typeof document.cancelFullScreen != 'undefined') {
        fullScreenApi.supportsFullScreen = true;
    } else {
        // check for fullscreen support by vendor prefix
        for (var i = 0, il = browserPrefixes.length; i < il; i++ ) {
            fullScreenApi.prefix = browserPrefixes[i];
 
            if (typeof document[fullScreenApi.prefix + 'CancelFullScreen' ] != 'undefined' ) {
                fullScreenApi.supportsFullScreen = true;
 
                break;
            }
        }
    }
 
    // update methods to do something useful
    if (fullScreenApi.supportsFullScreen) {
        fullScreenApi.fullScreenEventName = fullScreenApi.prefix + 'fullscreenchange';
 
        fullScreenApi.isFullScreen = function() {
            switch (this.prefix) {
                case '':
                    return document.fullScreen;
                case 'webkit':
                    return document.webkitIsFullScreen;
                default:
                    return document[this.prefix + 'FullScreen'];
            }
        }
        fullScreenApi.requestFullScreen = function(el) {
            return (this.prefix === '') ? el.requestFullScreen() : el[this.prefix + 'RequestFullScreen']();
        }
        fullScreenApi.cancelFullScreen = function(el) {
            return (this.prefix === '') ? document.cancelFullScreen() : document[this.prefix + 'CancelFullScreen']();
        }
    }
 
    // jQuery plugin
    if (typeof jQuery != 'undefined') {
        jQuery.fn.requestFullScreen = function() {
 
            return this.each(function() {
                if (fullScreenApi.supportsFullScreen) {
                    fullScreenApi.requestFullScreen(this);
                }
            });
        };
    }
    // export api
    window.fullScreenApi = fullScreenApi;
})();


(function () {

    var cmdMarkdownUrl = '/mdeditor/';

    // handle Tab keystroke
    $('#wmd-input').tabHandler();

    var converter1 = Markdown.getSanitizingConverter();

    // tell the converter to use Markdown Extra for tables, fenced_code_gfm, def_list
    Markdown.Extra.init(converter1, {extensions: ["tables", "fenced_code_gfm", "def_list"], highlighter: "prettify"});

    // To handle LaTeX expressions, to avoid the expression fail to work because of markdown syntax. inspired by stackeditor
    // This will handle $$LaTeX expression$$ only, so that $LaTeX expression$ could fail to handle either.
    bindMathJaxHooks(converter1);

    var markdownHelp = function () {
        var w = window.open(cmdMarkdownUrl);
        w.isEditablePage = false;
    }
    var options = {
        helpButton: { handler: markdownHelp },
        strings: Markdown.local.zh
    };

    var editor1 = new Markdown.Editor(converter1, null, options);

    var scrollLink = getScrollLink(); 
    scrollLink.onLayoutCreated();


    editor1.hooks.chain("onPreviewRefresh", function () {

        $('.prettyprint').each(function(){
            $(this).addClass('linenums');
        });
        prettyPrint(); // print code syntax for code snippet if there is.

        if ($('body').hasClass('theme-white')) {
            $('table').each(function() {
                $(this).addClass('table table-striped-white table-bordered');
            });
        } else {
            $('table').each(function() {
                $(this).addClass('table table-striped-black table-bordered');
            });
        }

        // Call onPreviewFinished callbacks when all async preview are finished, make sure sync actions have been ABOVE the line below.
        var counter = 0;
        var nbAsyncPreviewCallback = 2; // 1 for waitForImages below and 1 for MathJax below, they are both time consuming task, if only they are both done, begin to caculate md section and scroll bar.
        function tryFinished() {
            if(++counter === nbAsyncPreviewCallback) {
                scrollLink.onPreviewFinished();
            }
        }
        // We assume images are loading in the preview
        $("#wmd-preview").waitForImages(tryFinished);
        // TODO: could we cache the result to speed up ? This action is slow, especially, when there are multiple LaTeX expression on the page, google solution.
        MathJax.Hub.Queue(["Typeset",MathJax.Hub,"wmd-preview"]);
        MathJax.Hub.Queue(tryFinished);

        if (window.isEditablePage) { // Editing on markdown help page won't change local storage
            var preSaveArticle = $('#wmd-input').val();
            var savedArticle = $.localStorage('article');
            if (preSaveArticle != savedArticle) {
                $.localStorage('article', preSaveArticle);
            }
        }
    });
    scrollLink.onEditorConfigure(editor1);

    function popupEditorDialog(title, body, imageClass, placeholder) {
        $('#editorDialog').find('.modal-body input').val("");
        $('#editorDialog').find('.modal-body input').attr("placeholder", placeholder);
        $('#editorDialog').find('#editorDialog-title').text(title);
        $('#editorDialog').find('.modal-body p').text(body);
        $('#editorDialog').find('.modal-body i').removeClass().addClass(imageClass);
        $('#editorDialog').modal({keyboard : true});
    }

    // Custom insert link dialog
    editor1.hooks.set("insertLinkDialog", function(callback) {
        popupEditorDialog('链接', '请输入链接地址', 'icon-link icon-2x', 'http://example.com/ "可选标题"');
        editorDialogCallback = callback;
        return true; // tell the editor that we'll take care of getting the link url
    });

    // Custom insert image dialog
    var editorDialogCallback = null;
    editor1.hooks.set("insertImageDialog", function(callback) {
        popupEditorDialog('图片', '请输入图片地址', 'icon-picture icon-2x', 'http://example.com/images/diagram.jpg "可选标题"');
        editorDialogCallback = callback;
        return true; // tell the editor that we'll take care of getting the image url
    });

    $('#editorDialog').on('hidden', function(){
        if (editorDialogCallback) {
            var url = $('#editorDialog-confirm').data('url');
            if (url) {
                $('#editorDialog-confirm').removeData('url');
                editorDialogCallback(url);
            } else {
                editorDialogCallback(null);
            }
        }
    });

    $('#editorDialog-confirm').click(function(event) {
        var url = $('#editorDialog').find('.modal-body input').val();
        if (url) {
            $(this).data('url', url);
        }
        $('#editorDialog').modal('hide');
    });

    $('#editorDialog').on('shown', function(){
        $('#editorDialog').find('.modal-body input').focus();
    });


    // Make preview if it's inactive in 500ms to reduce the calls in onPreviewRefresh chains above and cpu cost.
    documentContent = undefined;
    var previewWrapper;
    previewWrapper = function(makePreview) {
        var debouncedMakePreview = _.debounce(makePreview, 500);
        return function() {
            if(documentContent === undefined) {
                makePreview();
                documentContent = '';
            } else {
                debouncedMakePreview();
            }
        };
    };

    // To make sure there is no overflow(scroll bar) on the whole page.
    function calculateEditorPreviewHeight() {
        var height = $(window).height() - $('.preview-container').position().top - 20;
        $('#wmd-input').height(height);
        $('.preview-container').height(height);
        $("#wmd-preview").height("auto");
    }

    $(window).resize(function() {
        calculateEditorPreviewHeight();
        scrollLink.buildSections();
    });


    // load md help doc from server.
    var mdUrl = 'editor/md-help';
    var isExternalUrl = false;
    if (window.isEditablePage === undefined) { // if window.isEditablePage is not undefined, means it's markdown help page.
        var keyUrl = '?url=';
        var indexOfKeyUrl = window.location.href.indexOf(keyUrl);
        if (indexOfKeyUrl != -1) {
            // jiawzhang NOTICE: make sure nginx and 'sh ~/uwsgi/uwsgi.sh' is running on ec2.
            mdUrl = 'http://www.zuoyebuluo.com?callback=?&url=' + encodeURIComponent(window.location.href.substring(indexOfKeyUrl + keyUrl.length));
            isExternalUrl = true;
            window.isEditablePage = false;
        } else {
            window.isEditablePage = true;
        }
    }

    $.ajax({
        type: 'GET',
        url: mdUrl,
        async: false,
        dataType: isExternalUrl ? 'jsonp' : '',
        success: function(data) {
            mainHandler(data);
        },
        error: function(e) {
            mainHandler('');
        }
    });

    function mainHandler(data) {
        if (data) {
            var article = null;
            var cursorPosition = 0;
            if (!window.isEditablePage) { // markdown help page is loading the certain text, regardless of local storage.
                article = '[『Cmd 技术渲染的沙箱页面，修改无法保存，点击此处编写自己的文档』](http://ghosertblog.github.io/mdeditor/ "中文在线 Markdown 编辑器")\n\n' + data;
            } else {
                var article = $.localStorage('article');
                if (!article) {
                    article = data;
                } else {
                    cursorPosition = article.length; // go to the end of the article, if the article is not help doc.
                }
            }
            // Populate editor value
            $('#wmd-input').val(article);

            // start editor.
            editor1.run(previewWrapper);

            // Load awesome font to button
            $('#wmd-bold-button > span').addClass('icon-bold muted');
            $('#wmd-italic-button > span').addClass('icon-italic muted');
            $('#wmd-link-button > span').addClass('icon-link muted');
            $('#wmd-quote-button > span').addClass('icon-quote-left muted');
            $('#wmd-code-button > span').addClass('icon-code muted');
            $('#wmd-image-button > span').addClass('icon-picture muted');
            $('#wmd-olist-button > span').addClass('icon-list-ol muted');
            $('#wmd-ulist-button > span').addClass('icon-list-ul muted');
            $('#wmd-heading-button > span').addClass('icon-list-alt muted');
            $('#wmd-hr-button > span').addClass('icon-minus muted');
            $('#wmd-undo-button > span').addClass('icon-undo muted');
            $('#wmd-redo-button > span').addClass('icon-repeat muted');
            $('#wmd-help-button > span').addClass('icon-question-sign muted');


            // create additional new buttons.
            $('#wmd-help-button').after('<li id="wmd-editor-full-button" class="wmd-button" title="编辑模式 Ctrl+M"><span class="icon-resize-full muted"></span></li>');
            $('#wmd-help-button').after('<li id="wmd-editor-small-button" class="wmd-button" title="预览模式 Ctrl+M"><span class="icon-resize-small muted"></span></li>');
            $('#wmd-help-button').css('margin-left', '50px');
            $('#wmd-editor-full-button').css('margin-left', '50px');
            $('#wmd-editor-small-button').css('margin-left', '50px');


            // hide #wmd-editor-small-button when initialization.
            $('#wmd-editor-small-button').hide();
            // hide #preview-editor-small-button when initialization.
            $('#preview-reader-small-button').hide();


            function buttonBinding(rowClassName, spanClassName) {
                // change color when hovering.
                $(rowClassName).hover(function() {
                    $(spanClassName).animate({color: '#F9F9F5'}, 400);
                },
                function() {
                    $(spanClassName).animate({color: '#BBBBBB'}, 400);
                });

                // enlarge the icon when hovering.
                $(spanClassName).hover(function() {
                    $(this).addClass('icon-large');
                },
                function() {
                    $(this).removeClass('icon-large');
                });
            }
            buttonBinding('.wmd-button-row', '.wmd-button > span');
            buttonBinding('.preview-button-row', '.preview-button > span');

            function clearAndNewFile() {
                var answer = confirm('新建文件将会清除当前的文件内容，请确认当前内容已保存');
                if (answer) {
                    $('#wmd-input').val('\n\n\n> *本文使用 [Cmd](http://ghosertblog.github.io/mdeditor/ "中文在线 Markdown 编辑器") 编写*');
                    $('#wmd-input').setCursorPosition(0);
                    editor1.refreshPreview();
                }
            }

            // new file button handler
            $('#preview-new-button').on('click', clearAndNewFile);


            function getCurrentMode() {
                var hash = window.location.hash;
                var currentMode = {isFullEditor: false, isFullReader: false, isEditorReader: false};
                if (hash == '#full-editor') {
                    currentMode.isFullEditor = true;
                } else if (hash == '#full-reader') {
                    currentMode.isFullReader = true;
                } else { // normal mode
                    currentMode.isEditorReader = true;
                }
                return currentMode;
            }

            function switchWmdButtonColor(buttonRow, buttonSpan, colorIn, colorOut) {
                $(buttonRow).unbind('hover');
                $(buttonSpan).animate({color: colorOut}, 400); // This resolve the no color changed issue.
                $(buttonSpan).removeClass('icon-large');
                $(buttonRow).hover(function() {
                    $(buttonSpan).animate({color: colorIn}, 400);
                },
                function() {
                    $(buttonSpan).animate({color: colorOut}, 400);
                });
            }

            // Loading theme setting from local storage and then apply.
            var siteThemeClassName = $.localStorage('siteThemeClassName');
            if (!siteThemeClassName) {
                siteThemeClassName = 'theme-white';
            }
            function applySiteTheme(siteThemeClassName) {
                $('.theme').each(function() {
                    $(this).removeClass('theme-white theme-black').addClass(siteThemeClassName);
                });
                var currentMode = getCurrentMode();
                if (currentMode.isFullReader) { // if it's full-reader page, change button color as well.
                    if (siteThemeClassName == 'theme-black') {
                        switchWmdButtonColor('.preview-button-row', '.preview-button > span', '#F9F9F5', '#BBBBBB');
                    } else {
                        switchWmdButtonColor('.preview-button-row', '.preview-button > span', '#2C3E50', '#999999');
                    }
                }
                if (siteThemeClassName == 'theme-white') {
                    $('table').each(function() {
                        $(this).removeClass('table-striped-black').addClass('table-striped-white');
                    });
                } else {
                    $('table').each(function() {
                        $(this).removeClass('table-striped-white').addClass('table-striped-black');
                    });
                }
                $.localStorage('siteThemeClassName', siteThemeClassName);
            }
            applySiteTheme(siteThemeClassName);

            function switchSiteTheme() {
                var currentMode = getCurrentMode();
                if ($('body').hasClass('theme-white')) {
                    applySiteTheme('theme-black');
                } else {
                    applySiteTheme('theme-white');
                }
            }

            // theme button handler
            $('#preview-theme-button').on('click', switchSiteTheme);


            // test whether the browser support fullscreen.
            if (fullScreenApi.supportsFullScreen) {
                $('#preview-fullscreen-button').on('click', function() {
                    if (fullScreenApi.isFullScreen()) {
                        fullScreenApi.cancelFullScreen();
                    } else {
                        fullScreenApi.requestFullScreen(document.documentElement);
                    }
                });
            } else {
                $('#preview-fullscreen-button').on('click', function() {
                    alert('您的浏览器不支持自动全屏，请尝试按 F11 切换全屏');
                });
            }

            /* ============================= Handle customized shortcut key binding. ========================================= */
            browserType = {
                isIE: /msie/.test(window.navigator.userAgent.toLowerCase()),
                isIE_5or6: /msie 6/.test(window.navigator.userAgent.toLowerCase()) || /msie 5/.test(window.navigator.userAgent.toLowerCase()),
                isOpera: /opera/.test(window.navigator.userAgent.toLowerCase()),
                isFirefox: /firefox/.test(window.navigator.userAgent.toLowerCase()),
                isChrome: /(chrome|chromium)/.test(window.navigator.userAgent.toLowerCase())
            };

            var keyEvent = 'keydown';
            if (browserType.isOpera || browserType.isFirefox) {
                keyEvent = 'keypress';
            }

            $(document).on(keyEvent, function(key) {
                // Check to see if we have a button key and, if so execute the callback.
                if ((key.ctrlKey || key.metaKey) && !key.shiftKey) {

                    var currentMode = getCurrentMode();

                    var keyCode = key.charCode || key.keyCode;
                    var keyCodeStr = String.fromCharCode(keyCode).toLowerCase();

                    switch (keyCodeStr) {
                        case "m":
                            if (!key.altKey) { // 'ctrl + m' for switching normal/full editor
                                if (currentMode.isEditorReader) {
                                    switchFullEditorMode();
                                } else if (currentMode.isFullEditor) {
                                    switchNormalModeFromFullEditorMode();
                                }
                            } else { // 'ctrl + alt + m' for switching normal/full reader
                                if (currentMode.isEditorReader) {
                                    switchFullReaderMode();
                                } else if (currentMode.isFullReader) {
                                    switchNormalModeFromFullReaderMode();
                                }
                            }
                            break;
                        case "j":
                            if (key.altKey) { // 'ctrl + alt + j' for switching site theme.
                                switchSiteTheme();
                                break;
                            }
                        case "h":
                            if (key.altKey) { // 'ctrl + alt + h' for markdown help.
                                markdownHelp();
                                break;
                            }
                        case "n":
                            if (key.altKey) { // 'ctrl + alt + n' for markdown help.
                                clearAndNewFile();
                                break;
                            }
                        default:
                            return;
                    }

                    if (key.preventDefault) {
                        key.preventDefault();
                    }

                    if (window.event) {
                        window.event.returnValue = false;
                    }
                }
            });

            /* ============================= Begin to Handle modes ========================================= */
            var switchFullEditorMode = function() {

                // hide all first.
                $('#container').hide();

                // show the hidden #editor-reader-full
                $('#editor-reader-full').removeClass('editor-reader-full-hidden').addClass('editor-reader-full-shown');

                // change full-small button
                $('#wmd-editor-small-button').show();
                $('#wmd-editor-full-button').hide();

                // Reset pagedown, pageup on textarea to make sure these keystroke works.
                $('#wmd-input').keydown(function(e) {
                    var targetScrollTop = null;
                    var pageHeight = $(window).height() - $('#wmd-panel-editor').position().top;
                    var lineHeight = parseInt($(this).css('line-height'));
                    pageHeight = pageHeight - lineHeight * 2; // subtract two lineHeights.
                    var currentScrollTop = $('#wmd-panel-editor').scrollTop();
                    if(e.keyCode === 33) { // pageup was pressed
                        targetScrollTop = currentScrollTop - pageHeight;
                    }
                    if(e.keyCode === 34) { // pagedown was pressed
                        targetScrollTop = currentScrollTop + pageHeight;
                    }
                    if (targetScrollTop != null) {
                        $('#wmd-panel-editor').animate({scrollTop: targetScrollTop}, 1); // set 1 here, since set 0 will lead to issue on IE.
                        var cursorPosition = scrollLink.getCursorPositionForPageDownUpInFullEditorMode(targetScrollTop);
                        if (cursorPosition != -1) {
                            var totalLength = $('#wmd-input').val().length;
                            if (cursorPosition > totalLength) {
                                cursorPosition = totalLength;
                            }
                            $('#wmd-input').setCursorPosition(cursorPosition);
                        }
                        e.preventDefault();
                    }
                });

                // reset hover colors on wmd button.
                if ($('body').hasClass('theme-white')) {
                    switchWmdButtonColor('.wmd-button-row', '.wmd-button > span', '#2C3E50', '#999999');
                }

                // Add wmd-button-bar and wmd-panel-editor to editor-reader-full
                $('#wmd-button-bar').removeClass().addClass('wmd-button-bar-full-shown');
                $('#editor-reader-full').append($('#wmd-button-bar'));

                $('#wmd-panel-editor').removeClass().addClass('wmd-panel-editor-full-shown');
                $('#editor-reader-full').append($('#wmd-panel-editor'));

                // make sure the width of #md-section-helper is same to #wmd-input, since scrollLink.buildSections() below will still caculate the sections.
                // (for move cursor correctly when pagedown/pageup pressed in full editor mode.)
                $('#wmd-input').css('max-width', '850px').focus();
                $('#md-section-helper').attr('style', 'max-width:850px; overflow:hidden; word-wrap:break-word; resize: none;');

                // binding new resize event.
                $(window).unbind('resize');
                function calculateEditFullHeight() {
                    // jquery autosize plugin to make sure textarea is autosized in full editor mode.
                    $('#wmd-input').autosize();
                    $('#wmd-input').trigger('autosize.resize');
                    $('#wmd-panel-editor').height($(window).height() - $('#wmd-panel-editor').position().top);
                }
                calculateEditFullHeight();
                scrollLink.buildSections();
                // '#wmd-input' should be autosize() when windows resize, otherwise, the display is not correct.
                $(window).resize(function() {
                    calculateEditFullHeight();
                    scrollLink.buildSections();
                });

                // change location hash to 'full-editor'
                window.location.hash = '#full-editor';
            }
            $('#wmd-editor-full-button').on('click', switchFullEditorMode);

            // This is basically the reversal action of switchFullEditorMode.
            var switchNormalModeFromFullEditorMode = function() {

                $('#wmd-input').unbind('keydown');

                // reset hover colors on wmd button.
                switchWmdButtonColor('.wmd-button-row', '.wmd-button > span', '#F9F9F5', '#BBBBBB');

                $('#wmd-button-bar').removeClass().addClass('pull-left');
                $('#editor-nav-bar').prepend($('#wmd-button-bar'));

                $('#wmd-panel-editor').removeClass().addClass('wmd-panel-editor');
                $('#left-column').prepend($('#wmd-panel-editor'));

                $('#md-section-helper').removeAttr('style');
                $('#wmd-panel-editor').removeAttr('style');
                $('#wmd-input').trigger('autosize.destroy');
                $('#wmd-input').removeAttr('style');

                $('#wmd-editor-small-button').hide();
                $('#wmd-editor-full-button').show();
                $('#editor-reader-full').removeClass('editor-reader-full-shown').addClass('editor-reader-full-hidden');
                $('#container').show();

                // binding new resize event.
                $(window).unbind('resize');
                calculateEditorPreviewHeight();
                scrollLink.buildSections();
                $(window).resize(function() {
                    calculateEditorPreviewHeight();
                    scrollLink.buildSections();
                });

                // change location hash to ''
                window.location.hash = '';
                $('#wmd-input').focus();
            }
            $('#wmd-editor-small-button').on('click', switchNormalModeFromFullEditorMode);
            
            var switchFullReaderMode = function() {

                // hide all first.
                $('#container').hide();

                // show the hidden #editor-reader-full
                $('#editor-reader-full').removeClass('editor-reader-full-hidden').addClass('editor-reader-full-shown');
                $('#editor-reader-full').css('position', 'static'); // to make sure the page scrollbar is present.

                // change full-small button
                $('#preview-reader-small-button').show();
                $('#preview-reader-full-button').hide();

                // Hide some buttons on preview
                $('#preview-new-button').hide();

                // reset hover colors on wmd button.
                if ($('body').hasClass('theme-white')) {
                    switchWmdButtonColor('.preview-button-row', '.preview-button > span', '#2C3E50', '#999999');
                }

                // Add toolbar to editor-reader-full
                $('#editor-reader-full').append('<div id="reader-full-toolbar"></div>');
                $('#reader-full-toolbar').append($('#preview-button-row'));

                $('#editor-reader-full').append($('#wmd-preview'));
                $('#wmd-preview').css('border', '0 none').css('margin', '0 auto').css('max-width', '850px').css('overflow-x', 'hidden').css('padding-top', '50px').css('padding-bottom', '50px').focus();

                $('.preview-button').css('font-size', '25px');
                $('#preview-button-row > li:visible').css('display', 'block').css('margin-bottom', '20px');
                $('#reader-full-toolbar').css('position', 'fixed').css('right', 40).css('top', 30);

                // unbind original resize event.
                $(window).unbind('resize');

                // change location hash to 'full-reader'
                window.location.hash = '#full-reader';
            }
            $('#preview-reader-full-button').on('click', switchFullReaderMode);

            // This is basically the reversal action of switchFullReaderMode.
            var switchNormalModeFromFullReaderMode = function() {

                // reset hover colors on wmd button.
                switchWmdButtonColor('.preview-button-row', '.preview-button > span', '#F9F9F5', '#BBBBBB');

                // restore items.
                $('.preview-button').removeAttr('style');
                $('#preview-button-row > li:visible').removeAttr('style');
                $('#preview-button-bar').prepend($('#preview-button-row'));
                $('#reader-full-toolbar').remove();

                $('#wmd-panel-preview').prepend($('#wmd-preview'));
                $('#wmd-preview').removeAttr('style');
                $('#wmd-preview').css('height', 'auto');

                // Show some buttons on preview
                $('#preview-new-button').show();
                // change full-small button
                $('#preview-reader-small-button').hide();
                $('#preview-reader-full-button').show();

                $('#editor-reader-full').removeAttr('style');
                $('#editor-reader-full').removeClass('editor-reader-full-shown').addClass('editor-reader-full-hidden');
                $('#container').show();

                // restore binding resize event.
                $(window).unbind('resize');
                calculateEditorPreviewHeight();
                scrollLink.buildSections();
                $(window).resize(function() {
                    calculateEditorPreviewHeight();
                    scrollLink.buildSections();
                });

                // change location hash to ''
                window.location.hash = '';
                $('#wmd-input').focus();
            }
            $('#preview-reader-small-button').on('click', switchNormalModeFromFullReaderMode);

            calculateEditorPreviewHeight();

            // Switch mode if there is.
            var currentMode = getCurrentMode();
            if (currentMode.isFullEditor) {
                $('#wmd-input').setCursorPosition(cursorPosition);
                switchFullEditorMode();
            } else if (currentMode.isFullReader) { // Don't set focus on '#wmd-input', otherwise, when first time press pagedown key on full reader page, Firefox can't scroll.
                switchFullReaderMode();
            } else { // normal mode
                $('#wmd-input').setCursorPosition(cursorPosition);
            }


            // Popup service change notification
            $('#notification-confirm').click(function(event) {
                var operation = $(this).data('operation');
                if (operation) {
                    $(this).removeData('operation');
                    operation();
                } else {
                    $('#notification').modal('hide');
                    window.location.reload();
                }
            });
            function _popupConfirm(title, content, operation) {
                $('#notification').find('#notification-title').html(title);
                $('#notification').find('.modal-body p').html(content);
                $('#notification').modal({keyboard : true});
                $('#notification-confirm').data('operation', operation); // store function here.
            }
            _popupConfirm("服务迁移通告", "为更好的提升用户体验，本站服务已迁移到独立服务器：<a href='http://www.zybuluo.com/mdeditor'>http://www.zybuluo.com/mdeditor</a>, 请点击【确定】按钮前往，如果您在本站有编辑中的文稿，请点击【取消】按钮先保存文稿，然后前往新站，本站服务会在近期关停，感谢一路支持。", function () {
                window.location = 'http://www.zybuluo.com/mdeditor';
            });
        }
    }
})();
