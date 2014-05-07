function getScrollLink() {

    var scrollLink = {
        extensionId: "scrollLink",
        extensionName: "Scroll Link",
        optional: true,
        settingsBloc: [
            '<p>Binds together editor and preview scrollbars.</p>',
            '<blockquote class="muted"><b>NOTE:</b>',
            '   The mapping between Markdown and HTML is based on the position of the title elements (h1, h2, ...) in the page.',
            '   Therefore, if your document does not contain any title, the mapping will be linear and consequently less accurate.',
            '</bloquote>'
        ].join("")
    };

    var msie = /msie/.test(navigator.userAgent.toLowerCase());

    var mdSectionList = [];
    var htmlSectionList = [];
    function pxToFloat(px) {
        return parseFloat(px.substring(0, px.length - 2));
    }
    var buildSections = _.debounce(function() {

        // Try to find Markdown sections by looking for titles
        var editorElt = $("#wmd-input");
        mdSectionList = [];
        // This textarea is used to measure sections height
        var textareaElt = $("#md-section-helper");
        // It has to be the same width as wmd-input
        textareaElt.width(editorElt.width());
        // Consider wmd-input top padding (will be used for 1st and last section) 
        var padding = pxToFloat(editorElt.css('padding-top'));
        var offset = 0, mdSectionOffset = 0;
        
        function addMdSection(sectionText) {
            var sectionHeight = padding;
            if(sectionText !== undefined) {
                textareaElt.val(sectionText);
                if (msie) {
                    sectionHeight += textareaElt.prop('scrollHeight') + 6;
                } else {
                    sectionHeight += textareaElt.prop('scrollHeight');
                }
            }
            var newSectionOffset = mdSectionOffset + sectionHeight;
            mdSectionList.push({
                startOffset: mdSectionOffset,
                endOffset: newSectionOffset,
                height: sectionHeight,
                length: sectionText == undefined ? 0 : sectionText.length + 1
            });
            mdSectionOffset = newSectionOffset;
            padding = 0;
        }
        // Create MD sections by finding title patterns (excluding gfm blocs)
        var text = editorElt.val() + "\n\n";
        text.replace(/^```.*\n[\s\S]*?\n```|(^.+[ \t]*\n=+[ \t]*\n+|^.+[ \t]*\n-+[ \t]*\n+|^\#{1,6}[ \t]*.+?[ \t]*\#*\n+)/gm, function(match, title, matchOffset) {
            if(title) {
                // We just found a title which means end of the previous section
                // Exclude last \n of the section
                var sectionText = undefined;
                if(matchOffset > offset) {
                    sectionText = text.substring(offset, matchOffset - 1);
                }
                addMdSection(sectionText);
                offset = matchOffset;
            }
            return "";
        });
        // Last section
        // Consider wmd-input bottom padding and exclude \n\n previously added
        padding += pxToFloat(editorElt.css('padding-bottom'));
        addMdSection(text.substring(offset, text.length - 2));

        // Try to find corresponding sections in the preview
        var previewElt = $(".preview-container");
        htmlSectionList = [];
        var htmlSectionOffset = 0;
        var previewScrollTop = previewElt.scrollTop();
        // Each title element is a section separator
        var wmdPreviewTop = $(".preview-container").position().top;
        $("#wmd-preview").children("h1,h2,h3,h4,h5,h6").each(function() {
            // Consider div scroll position and header element top margin
            var newSectionOffset = $(this).position().top - wmdPreviewTop + previewScrollTop + pxToFloat($(this).css('margin-top'));
            htmlSectionList.push({
                startOffset: htmlSectionOffset,
                endOffset: newSectionOffset,
                height: newSectionOffset - htmlSectionOffset
            });
            htmlSectionOffset = newSectionOffset;
        });
        // Last section
        var scrollHeight = previewElt.prop('scrollHeight');
        htmlSectionList.push({
            startOffset: htmlSectionOffset,
            endOffset: scrollHeight,
            height: scrollHeight - htmlSectionOffset
        });

        // apply Scroll Link
        lastEditorScrollTop = -10;
        lastPreviewScrollTop = -10;
        isScrollPreview = false;
        runScrollLink();
    }, 500);

    // -10 to be sure the gap is > 9
    var lastEditorScrollTop = -10;
    var lastPreviewScrollTop = -10;
    var isScrollPreview = false;
    // 每次滑动都会调用
    var runScrollLink = _.debounce(function() {
        if(mdSectionList.length === 0 || mdSectionList.length !== htmlSectionList.length) {
            return;
        }
        var editorElt = $("#wmd-input");
        var editorScrollTop = editorElt.scrollTop();
        var previewElt = $(".preview-container");
        var previewScrollTop = previewElt.scrollTop();
        
        function animate(srcScrollTop, srcSectionList, destElt, destSectionList, currentDestScrollTop, callback) {
            // Find the section corresponding to the offset
            var sectionIndex = undefined;
            var srcSection = _.find(srcSectionList, function(section, index) {
                sectionIndex = index;
                return srcScrollTop < section.endOffset;
            });
            if(srcSection === undefined) {
                // Something wrong in the algorithm...
                return;
            }
            var posInSection = (srcScrollTop - srcSection.startOffset) / srcSection.height;
            var destSection = destSectionList[sectionIndex];
            var destScrollTop = destSection.startOffset + destSection.height * posInSection;
            destScrollTop = _.min([
                destScrollTop,
                destElt.prop('scrollHeight') - destElt.outerHeight()
            ]);
            if(Math.abs(destScrollTop - currentDestScrollTop) <= 1) {
                // Skip the animation in case it's not necessary
                callback(currentDestScrollTop);
                return;
            }
            
            /*
            var d = 100; // 时间间隔
            var top = destScrollTop;
            var nowTop = destElt.scrollTop();
            log(srcScrollTop + " => " + destScrollTop)
    		for(var i = 0; i < d; i++) {
    			setTimeout(
    			(function(top) {
    				return function() {
			            destElt.slimScroll({ scrollTo: top, onlyScrollBar: true}); 
    				}
    			})(nowTop + 1.0*i*(top-nowTop)/d), i);
    		}
    		// 最后必然执行
    		setTimeout(function() {
//                callback(destScrollTop);
    		}, d+5);
    		*/
            
            destElt.animate({
                scrollTop: destScrollTop
            }, 500, function() {
                callback(destScrollTop);
            });
            setTimeout(
            		function() {
			            // 仅仅是移动bar而不移动内容
			            destElt.slimScroll({ scrollTo: destScrollTop, onlyScrollBar: true}); 
            		}, 500);
        }
        
        // Perform the animation if diff > 9px
//        console.log(editorScrollTop , lastEditorScrollTop);
        if(isScrollPreview === false && Math.abs(editorScrollTop - lastEditorScrollTop) > 1) {
            // Animate the preview
            lastEditorScrollTop = editorScrollTop;
//            console.log("--")
            animate(editorScrollTop, mdSectionList, previewElt, htmlSectionList, previewScrollTop, function(destScrollTop) {
                lastPreviewScrollTop = destScrollTop;
//                console.log(previewScrollTop, destScrollTop)
            });
        }
        else if(Math.abs(previewScrollTop - lastPreviewScrollTop) > 1) {
            // Animate the editor
            lastPreviewScrollTop = previewScrollTop;
            animate(previewScrollTop, htmlSectionList, editorElt, mdSectionList, editorScrollTop, function(destScrollTop) {
                lastEditorScrollTop = destScrollTop;
            });
        }
    }, 500);

    scrollLink.buildSections = function() {
        buildSections();
    };

    // 触发scroll时调用
    scrollLink.onLayoutCreated = function() {
        $(".preview-container").scroll(function() {
            isScrollPreview = true;
            runScrollLink();
        });
        $("#wmd-input").scroll(function() {
            isScrollPreview = false;
            runScrollLink();
        });
    };

    scrollLink.onEditorConfigure = function(editor) {
        editor.getConverter().hooks.chain("postConversion", function(text) {
            // To avoid losing scrolling position before elements are fully
            // loaded
            var previewElt = $("#wmd-preview");
            previewElt.height(previewElt.height());
            return text;
        });
    };

    scrollLink.onPreviewFinished = function() {
        // Now set the correct height
        $("#wmd-preview").height("auto");
        buildSections();
    };

    scrollLink.getCursorPositionForPageDownUpInFullEditorMode = function(scrollTop) {
        var sectionIndex = undefined;
        var srcSection = _.find(mdSectionList, function(section, index) {
            sectionIndex = index;
            return scrollTop < section.endOffset;
        });
        if(srcSection === undefined) {
            // Something wrong in the algorithm...
            return -1;
        }
        var caretPosition = 0;
        for (var index = 0; index <= sectionIndex; index++) {
            caretPosition = caretPosition + mdSectionList[index].length;
        }
        return caretPosition;
    }

    return scrollLink;
}

