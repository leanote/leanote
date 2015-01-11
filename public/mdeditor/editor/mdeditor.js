// leanote's markdown editor
// use require js
var markdownLang = "Markdown.en";
if(LEA.locale == "zh") {
	markdownLang = "Markdown.zh";
}
define('mdeditor', 
['Markdown.Converter', 'Markdown.Sanitizer', 'Markdown.Editor', markdownLang, 'Markdown.Extra', 
'underscore', 'scrollLink', 'mathJax', 'jquery.waitforimages', 'pretty'], 
function(){

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
/*
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

*/

(function () {
    // handle Tab keystroke
    $('#wmd-input').tabHandler();

    var converter1 = Markdown.getSanitizingConverter();
    Converter = converter1;

    // tell the converter to use Markdown Extra for tables, fenced_code_gfm, def_list
    Markdown.Extra.init(converter1, {extensions: ["tables", "fenced_code_gfm", "def_list"], highlighter: "prettify"});

    // To handle LaTeX expressions, to avoid the expression fail to work because of markdown syntax. inspired by stackeditor
    // This will handle $$LaTeX expression$$ only, so that $LaTeX expression$ could fail to handle either.
    bindMathJaxHooks(converter1);

    // 弹框显示markdown语法
    var markdownHelp = function () {
        window.open("http://leanote.com/blog/view/531b263bdfeb2c0ea9000002");
        return;
    }
    var options = {
        helpButton: { handler: markdownHelp },
        strings: Markdown.local.zh
    };

    var editor1 = new Markdown.Editor(converter1, null, options);
    MarkdownEditor = editor1;

    var scrollLink = getScrollLink(); 
    ScrollLink = scrollLink;
    scrollLink.onLayoutCreated();

    editor1.hooks.chain("onPreviewRefresh", function () {
        $("#left-column pre").addClass("prettyprint linenums");
        prettyPrint();

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
        if(typeof MathJax != "undefined") {
	        MathJax.Hub.Queue(["Typeset",MathJax.Hub,"wmd-preview"]);
	        MathJax.Hub.Queue(tryFinished);
        } else {
            scrollLink.onPreviewFinished();
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
        popupEditorDialog(options.strings.linkTitle, options.strings.linkAddress, 'fa fa-link', options.strings.linkExample);
        editorDialogCallback = callback;
        return true; // tell the editor that we'll take care of getting the link url
    });

    // Custom insert image dialog
    var editorDialogCallback = null;
    editor1.hooks.set("insertImageDialog", function(callback) {
        popupEditorDialog(options.strings.imageTitle, options.strings.imageAddress, 'fa fa-picture-o', options.strings.imageExample);
        editorDialogCallback = callback;
        return true; // tell the editor that we'll take care of getting the image url
    });

    $('#editorDialog').on('hidden.bs.modal', function(){
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

    $('#editorDialog').on('shown.bs.modal', function(){
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

    $(window).resize(function() {
        scrollLink.buildSections();
    });

    // 渲染编辑器
    mainHandler();

    function mainHandler(data) {
        var article = null;
        var cursorPosition = 0;

        // start editor.
        editor1.run(previewWrapper);

        // Load awesome font to button
        $('#wmd-bold-button > span').addClass('fa fa-bold');
        $('#wmd-italic-button > span').addClass('fa fa-italic');
        $('#wmd-link-button > span').addClass('fa fa-link');
        $('#wmd-quote-button > span').addClass('fa fa-quote-left');
        $('#wmd-code-button > span').addClass('fa fa-code');
        $('#wmd-image-button > span').addClass('fa fa-picture-o');
        $('#wmd-olist-button > span').addClass('fa fa-list-ol');
        $('#wmd-ulist-button > span').addClass('fa fa-list-ul');
        $('#wmd-heading-button > span').addClass('fa fa-list-alt');
        $('#wmd-hr-button > span').addClass('fa fa-minus');
        $('#wmd-undo-button > span').addClass('fa fa-undo');
        $('#wmd-redo-button > span').addClass('fa fa-repeat');
        
        $('#wmd-help-button > span').addClass('fa fa-question-circle');

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

        function getCurrentMode() {
            var currentMode = {isFullEditor: false, isFullReader: false, isEditorReader: false};
            return currentMode;
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
                    /*
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
                        */
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
    } // mainHander
})();


});