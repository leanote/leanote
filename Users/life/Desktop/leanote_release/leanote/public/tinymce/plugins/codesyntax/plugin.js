/**
 * life
 * 用户所在位置是<pre>下
 * 用户选中代码, 点击 <>, 则将代码贴上, 并且只贴纯文本
 */
(function() {
	tinymce.create('tinymce.plugins.codesyntax', {
		/**
		 * Initializes the plugin, this will be executed after the plugin has been created.
		 * This call is done before the editor instance has finished it's initialization so use the onInit event
		 * of the editor instance to intercept that event.
		 *
		 * @param {tinymce.Editor} ed Editor instance that the plugin is initialized in.
		 * @param {string} url Absolute URL to where the plugin is located.
		 */
		init : function(ed, url) {
			// Register the command so that it can be invoked by using tinyMCE.activeEditor.execCommand('mceExample');
			ed.addCommand('mceSyntaxHL', function() {
				ed.windowManager.open({
					title: "插入代码",
					file : url + '/dialog.htm',
					width : 520,// + parseInt(ed.getLang('.delta_width', 0)),
					height : 350// + parseInt(ed.getLang('.delta_height', 0)),
				});
			});
			
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

			// Register example button
			ed.addButton('codesyntax', {
				title : '插入代码',
				cmd : 'mceSyntaxHL',
				// icon: "code"
				image : url + '/img/file-xml.png'
			});
			
			// Register example button
			ed.addButton('toggleCode', {
				title : '文本<->代码切换',
				cmd : 'toggleCode',
				icon: "code"
				// image : url + '/img/highlight.gif'
			});
			
			ed.addShortcut('ctrl+shift+c', '', 'toggleCode');
			ed.addShortcut('command+shift+c', '', 'toggleCode');
			
			// 应该放在其它地方的
			// indent outdent
			ed.on('keydown', function(e) {
				var num = e.which ? e.which : e.keyCode;
		    	if (num == 9) { // tab pressed
		    		if(!e.shiftKey) {
//				            ed.execCommand('mceInsertRawHTML', false, "&nbsp;&nbsp;&nbsp;&nbsp;"); // inserts tab
		                ed.execCommand('Indent');
		    		} else {
		    			// shift + tab
		    			// delete 4 个空格
		                ed.execCommand('Outdent');
		    		}
		            e.preventDefault();
		            e.stopPropagation();   			
		            return false;
		       }
			});
		},

		/**
		 * Creates control instances based in the incomming name. This method is normally not
		 * needed since the addButton method of the tinymce.Editor class is a more easy way of adding buttons
		 * but you sometimes need to create more complex controls like listboxes, split buttons etc then this
		 * method can be used to create those.
		 *
		 * @param {String} n Name of the control to create.
		 * @param {tinymce.ControlManager} cm Control manager to use inorder to create new control.
		 * @return {tinymce.ui.Control} New control instance or null if no control was created.
		 */
		createControl : function(n, cm) {
			return null;
		}
	});


	/**
	 * Parses Syntax Highlighters parameters within the class attribute and
	 * returns them in an associative array.
	 *
	 * @param params string or Syntax Highlighter parameters to parse
	 *
	 * @return array
	 */
	function getParameters(params) {
		var parsedItems = params.replace(/ /g, '').replace(/;([^;]*)$/, '').split(';');
		var paramArray = new Array();
		for(var i = 0; i < parsedItems.length; i++)
		{
			var option = parsedItems[i].split(':');
			paramArray[option[0]] = option[1];
		}
		return paramArray;
	}

	// Register plugin
	tinymce.PluginManager.add('codesyntax', tinymce.plugins.codesyntax);
})();