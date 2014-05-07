tinymce.PluginManager.add('leanote_code', function(editor, url) {
	var ed = editor;
	
	// 切换代码
	function toggleLang(value) {
		var node = ed.selection.getNode();
		var selectedContent = ed.selection.getContent(); // 包含了html标签
		var everBookmark = ed.selection.getBookmark(); // 光标, 为了处理后重新定位到那个位置
		var text;
		
		// 如果在pre下, 只需要改brush即可
		if(node.nodeName == 'PRE') {
			$(node).attr("class", 'brush:' + value);
		// 不能把BODY转成pre
		} else if(node.nodeName == "BODY") {
			return;
		// 否则, 当所选的, 或者当前行转成code
		} else {
			try {
				text = $(selectedContent).text();
			} catch(e) {
			}
			if(!text) {
				text = $(node).text();
				$(node).replaceWith('<pre class="brush:' + value + '"> '+ text + "</pre>");
			} else {
				ed.insertContent('<pre class="brush:' + value + '"> ' + text + "</pre>");
			}
		}
		ed.selection.moveToBookmark(everBookmark);
	}
	
	// 切换代码
	// 用户选择了/用户光标所在行
//	ed.addCommand('toggleCode', );
	
	//----------------
	function createListBoxChangeHandler() {
		return function() {
			var self = this;
			editor.on('nodeChange', function() {
				// 得到node
				var value = null;
				try {
					var node = editor.selection.getNode();
					if(node.nodeName == "PRE") {
						value = $.trim($(node).attr("class").split(":")[1]);
					}
				} catch(e) {
				}
				self.value(value);
			});
		};
	}
	    
    editor.addButton('leanote_code', function() {
    	var langs = ["CSS:css", 
    	             "C/C++:cpp", 
    	             "C#:c#",
    	             "Javascript:javascript", 
    	             "Java:java", 
    	             "PHP:php",
    	             "Python:python",
    	             "Shell:shell", 
    	             "Delphi:delphi",
    	             "Golang:golang",
    	             "Erlang:erlang",
    	             "Groovy:groovy",
    	             "ActionScript:actionScript"
    	             ];
		var items = [];
    	for(var i in langs) {
    		var each = langs[i].split(":");
    		items.push({text:each[0], value:each[1]});
    	}
		return {
			type: 'listbox',
			text: "代码语言",
			tooltip: "ctrl+shift+c 切换代码",
			values: items,
			fixedWidth: true,
			onselect: function(e) {
				// 这里, change或不转成code
				if (e.control.settings.value) {
					toggleLang(e.control.settings.value);
				}
			},
			onPostRender: createListBoxChangeHandler(items)
		};
	});
    
    //----------------
	// 切换代码
	// 用户选择了/用户光标所在行
	ed.addCommand('toggleCode', function() {
		var node = ed.selection.getNode();
		var selectedContent = ed.selection.getContent(); // 包含了html标签
		var everBookmark = ed.selection.getBookmark(); // 光标, 为了处理后重新定位到那个位置
		var text;
		try {
			text = $.trim($(selectedContent).text());
		} catch(e) {
		}
		// 可能不是一个完整的html, 可能是一个文本此时.html()无
		if(!text) {
			text = $.trim(selectedContent);
		}
		// 如果是用户选择了某内容, 其父<pre>不会包含进来, 此时还是要判断node的父是否是PRE, 
		if(text) {
			// 是pre, 那么 去掉pre之
			if(node.nodeName == "PRE") {
				$(node).replaceWith("<p>" + $(node).html() + "</p>");
			} else {
				// 不是, 那么替换成<pre>
				ed.insertContent("<pre>" + text + "</pre>");
			}
		} else {
			// 没有选择内容, 那么是鼠标的所处行, 同样检测该行是否是
			if(node.nodeName == "PRE") {
				$(node).replaceWith("<p>" + $(node).html().replace(/\n/g, "<br />") + "</p>");
				
			} else {
				// 不是, 那么替换成<pre>
				$(node).replaceWith("<pre>" + $(node).html() + "</pre>");
			}
		}
		ed.selection.moveToBookmark(everBookmark);
	});
    
    ed.addShortcut('ctrl+shift+c', '', 'toggleCode');
	ed.addShortcut('command+shift+c', '', 'toggleCode');
});