/**
 * leanote code plugin
 */

// tinymce.PluginManager.requireLangPack('leanote_code');

tinymce.PluginManager.add('leanote_code', function(editor, url) {
	var me = this;
	var ed = editor;

	// clearUndo没必要, 支持
	// depreciated
	function clearUndo() {
		// 必须要setTimeout
		ed.undoManager.clear();
		setTimeout(function() {
			ed.undoManager.clear();
		});
	}

	function html2BreakLine(html) {
		if(!html) {
			return html;
		}
		if(typeof html == "object") {
			html = $(html).html();
		}
		// 先<br />=>\n, 再<p>=>\n, 再把连续的两个\n => 1个\n
		return html.replace(/\<br *\/*\>/gi,"\n").replace(/<\/(p|li|div|ul|ol|hr)>/, "\n").replace(/(<([^>]+)>)/gi, "").replace(/\n\n/g, "\n");
	}
	function html2BreakLineForPre(html) {
		if(!html) {
			return html;
		}
		if(typeof html == "object") {
			html = $(html).html();
		}
		return html.replace(/\n/g, "<br />"); // .replace(/\<br *\/*\>/gi,"\n").replace(/<\/(p|li|div|ul|ol|hr)>/, "\n").replace(/(<([^>]+)>)/gi, "");
	}

	// 在toggle成pre或ace时
	// 最后没有元素, 或者元素不是p, 则在最后插入之
	function insertPIfNotExists() {
		var children = $('#editorContent').children();
		var lastChild = children && children.length > 0 ? children[children.length - 1] : null;
		if (!lastChild || lastChild.tagName != 'P') {
			$('#editorContent').append('<p><br data-mce-bogus="1"></p>');
		}
	}

	// brush 刷子
	function toggleCode(brush) {
		if (LEA.readOnly) {
			return;
		}
		ed = tinymce.activeEditor;
		var node = ed.selection.getNode();

		var selectedContent = ed.selection.getContent(); // 包含了html标签
		// var everBookmark = ed.selection.getBookmark(); // 光标, 为了处理后重新定位到那个位置
		var text;

		var aceEditorAndPre = LeaAce.isInAce(node);
		var aceEditor = false;
		var $pre = false;
		if(aceEditorAndPre) {
			aceEditor = aceEditorAndPre[0];
			$pre = aceEditorAndPre[1];
		}

		// 去掉
		// 当pre->text时会有遗留, 这里干脆清除之
		$('#editorContent .toggle-raw').remove();

		var brushClasses = 'class="brush:' + brush + '"';

		// log('curNode:')
		// log(node);
		// s = ed.selection;
		// log(s);
		// log($pre.get(0));
		
		if(brush && brush != "convert") {
			if(aceEditor) {
				aceEditor.session.setMode("ace/mode/" + brush);
			}
			// 不支持ace的情况
			if(!$pre && node.nodeName  == "PRE") {
				$pre = $(node);
			}
			if($pre) {
				var everBrush = LeaAce.getPreBrush($pre);
				$pre.removeClass(everBrush).addClass("brush:" + brush);
				return;
			}
		// 不能把BODY转成pre
		} else if(brush && (node.nodeName == "BODY" || $(node).attr('id') == 'editorContent')) {
			return;
		}

		// 对于不支持ace的情况
		if(!LeaAce.canAce()) {
			if(node.nodeName != "PRE") {
				node = $(node).closest("pre").get(0);
			}
			if(node && node.nodeName == "PRE") {
				var $pre = $(node);
				var val = $pre.html();
				if(val) {
					val = val.replace(/\n/g, "<br />");
				}
				$pre.replaceWith("<p>" + val + "</p>");
			} else {
				try {
					text = $.trim($(selectedContent).text());
				} catch(e) {
				}
				// 可能不是一个完整的html, 可能是一个文本此时.html()无
				if(!text) {
					text = $.trim(selectedContent);
				}
				var pre = null;
				var id = LeaAce.getAceId();
				if(text) {
					// 不是, 那么替换成<pre>
					text = html2BreakLineForPre(text);
					pre = '<pre id="' + id + '" ' + brushClasses + '>' + text + '</pre>';
					ed.insertContent(pre);
				} else {
					if(node) {
						text = html2BreakLineForPre(node);
						pre = '<pre id="' + id + '" ' + brushClasses + '>' + text + '</pre>';
						$(node).replaceWith(pre);
					} else {
						pre = '<pre id="' + id + '" ' + brushClasses + '>' + text + '</pre>';
						ed.insertContent(pre);
					}
				}
				if(pre) {
					insertPIfNotExists();
					/*
					var rng = ed.selection.getRng();
					var $pre = $(pre);
					rng.setStart($pre.get(0), 0);
					rng.setEnd($pre.get(0), 0);
					*/
				}
			}
			// ed.selection.moveToBookmark(everBookmark);
			return;
		}

		// 支持ace情况

		// 如果是用户选择了某内容, 其父<pre>不会包含进来, 此时还是要判断node的父是否是PRE, 
		var id = LeaAce.getAceId(); 

		// 防止ace处理时添加额外的历史
		// tinymce.activeEditor.undoManager.add();
		LeaAce.disableAddHistory();

		if(aceEditor) {
			var val = aceEditor.getValue();
			// TODO 实体转成&lt;&rg;
			val = val.replace(/</g, "&lt;");
			val = val.replace(/>/g, "&gt;");
			val = val.replace(/\n/g, "<br />");
			$pre.replaceWith("<p>" + val + "</p>");
			aceEditor.destroy();
		} else {
			if(node.nodeName == "PRE") {
				// $(node).find('.toggle-raw').remove();
				// 表示在pre下, 但不是aceEditor, toggle后的
				var $pre = $(node);
				var val = $pre.html();
				if(val) {
					val = val.replace(/\n/g, "<br />");
				}
				$pre.replaceWith("<p>" + val + "</p>");
				return;
			}
			var text = selectedContent;
			if(!text && (node.nodeName == "BODY" || $(node).attr('id') == 'editorContent')) {
				return;
			}

			if(text) {
				// 不是, 那么替换成<pre>
				// log(text);
				text = html2BreakLine(text);
				// log(text);
				ed.insertContent('<pre id="' + id + '" ' + brushClasses + '>' + text + '</pre>');
			} else {
				// 不是, 那么替换成<pre>
				text = html2BreakLine(node);
				$(node).replaceWith("<pre id='" + id  + "'" + brushClasses + ">" + text + "</pre>");
			}

			var aceEditor = LeaAce.initAce(id);
			if(aceEditor) {
				aceEditor.focus();
				if(brush && brush != "convert") {
					aceEditor.session.setMode("ace/mode/" + brush);
				}

				insertPIfNotExists();
			}
		}
		// ed.selection.moveToBookmark(everBookmark);
		LeaAce.resetAddHistory();
	};
		
	// 切换代码
	// 用户选择了/用户光标所在行
	
	//----------------
	function createListBoxChangeHandler() {
		return function() {
			var self = this;
			editor.on('nodeChange', function() {
				// 得到node
				var value = null;
				try {
					var node = editor.selection.getNode();
					if(node.nodeName != "PRE") {
						node = $(node).closest("pre").get(0);
					}
					if(node) {
						var aceEditorAndPre = LeaAce.isInAce(node);
						var aceEditor = false;
						var $pre = false;
						if(aceEditorAndPre || node.nodeName == "PRE") {
							if(aceEditorAndPre) {
								aceEditor = aceEditorAndPre[0];
								$pre = aceEditorAndPre[1];
							} else {
								$pre = $(node);
							}
							var brush = LeaAce.getPreBrush($pre);
							value = $.trim(brush.split(":")[1]);

							self.diableValue("convert", false);
						} else {
							self.diableValue("convert", true);
						}
					}
				} catch(e) {
					log(e);
				}
				if(value != "convert") {
					self.value(value);
				}
			});
		};
	}
	    
    editor.addButton('leanote_code', function() {
    	var langs = [
    		"Convert Code:convert",
			"CSS:css", 
            "HTML:html",
            "Javascript:javascript", 
            "C/C++:c_cpp", 
            "C#:csharp",
            "Java:java", 
            "Objective-c:objectivec", 
            "PHP:php",
            "Python:python",
            "Ruby:ruby",
            "Shell:sh", 
            "Delphi:delphi",
            "Golang:golang",
            "Erlang:erlang",
            "Groovy:groovy",
            "Latex:latex",
            "Xml:xml",
            "ActionScript:actionScript",
         ];
		var items = [];
    	for(var i in langs) {
    		var each = langs[i].split(":");
    		items.push({text:each[0], value:each[1]});
    	}
		return {
			type: 'listbox',
			text: "Language",
			tooltip: "`ctrl/cmd+shift+c` toggle code",
			values: items,
			fixedWidth: true,
			onselect: function(e) {
				// 这里, change或不转成code
				if (e.control.settings.value) {
					toggleCode(e.control.settings.value);
				}
			},
			onPostRender: createListBoxChangeHandler(items)
		};
	});

	editor.addButton('leanote_inline_code', {
		icon: 'code',
		tooltip: 'Inline Code',
		stateSelector: 'code',
		onclick: function() {
			editor.execCommand('mceToggleFormat', false, 'code');
		}
	});

	if(LeaAce.canAce()) {
		editor.addButton('leanote_ace_pre', {
			icon: 'ace-pre',
			tooltip: 'Toggle ace with raw html',
			active: LeaAce.isAce === false,
			onclick: function() {
				// 表示之前是ace->pre状态
				// 现在转成ace
				if(LeaAce.isAce === false) {
					this.active(false);
					LeaAce.isAce = true;
					LeaAce.initAceFromContent(editor);
				// 转成pre
				} else {
					this.active(true);
					LeaAce.allToPre(editor);
					LeaAce.isAce = false;
				}
			}
		});
	}

    //----------------
	// 切换代码
	// 用户选择了/用户光标所在行
	ed.addCommand('toggleCode', toggleCode);
    
    ed.addShortcut('ctrl+shift+c', '', 'toggleCode');
	ed.addShortcut('meta+shift+c', '', 'toggleCode');

	// life
	if(LeaAce.canAce()) {
		editor.on('keydown', function(e) {
			// paste时
			var ace = LeaAce.nowIsInAce();
			if(ace) {
				setTimeout(function() {
					ace[0].focus();
				});
				return true;
			}

			// todo
			/*
			var keyCode = e.keyCode;
			if(keyCode == 8 || keyCode == 46) {
				// alert(keyCode);
				// 如果之前是Ace
				// var node = editor.selection.getNode();
				// log(node);
				log(e);
				return false;
			}
			*/
		});
	}

	// indent outdent
	ed.on('keydown', function(e) {
		var num = e.which ? e.which : e.keyCode;
    	if (num == 9) { // tab pressed
    		if(!e.shiftKey) {
	    		// var node = ed.selection.getNode();
	    		/*
				if(node.nodeName == "PRE") {
                    ed.execCommand('mceInsertHTML', false, '\x09'); // inserts tab
				} else {
				*/
				// 如果是在li下的, 就不要控制
				var node = ed.selection.getNode();
				if (node && (node.nodeName == 'LI' || $(node.closest('li')).length > 0)) {
					return true;
				}
				ed.insertContent("&nbsp;&nbsp;&nbsp;&nbsp;");
	            e.preventDefault();
	            e.stopPropagation();   			
	            return false;
                // ed.execCommand('mceInsertHTML', false, "&nbsp;&nbsp;&nbsp;&nbsp;"); // inserts 空格
				// }
    		} else {
    			// delete 4 个空格
				// ed.execCommand('Outdent');
    		}
       }
	});
});