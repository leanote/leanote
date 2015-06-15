// Included from: js/tinymce/plugins/paste/classes/Clipboard.js

/**
 * Clipboard.js
 *
 * Copyright, Moxiecode Systems AB
 * Released under LGPL License.
 *
 * License: http://www.tinymce.com/license
 * Contributing: http://www.tinymce.com/contributing
 */

/**
 * This class contains logic for getting HTML contents out of the clipboard.
 *
 * We need to make a lot of ugly hacks to get the contents out of the clipboard since
 * the W3C Clipboard API is broken in all browsers: Gecko/WebKit/Blink. We might rewrite
 * this the way those API:s stabilize.
 *
 * Current implementation steps:
 *  1. On keydown with paste keys Ctrl+V or Shift+Insert create
 *     a paste bin element and move focus to that element.
 *  2. Wait for the browser to fire a "paste" event and get the contents out of the paste bin.
 *  3. Check if the paste was successful if true, process the HTML.
 *  (4). If the paste was unsuccessful use IE execCommand, Clipboard API, document.dataTransfer old WebKit API etc.
 * 
 * @class tinymce.pasteplugin.Clipboard
 * @private
 */
define("tinymce/pasteplugin/Clipboard", [
	"tinymce/Env",
	"tinymce/util/VK",
	"tinymce/pasteplugin/Utils"
], function(Env, VK, Utils) {
	return function(editor) {
		var self = this, pasteBinElm, lastRng, keyboardPasteTimeStamp = 0;
		var pasteBinDefaultContent = '%MCEPASTEBIN%', keyboardPastePlainTextState;

		/**
		 * Pastes the specified HTML. This means that the HTML is filtered and then
		 * inserted at the current selection in the editor. It will also fire paste events
		 * for custom user filtering.
		 *
		 * @param {String} html HTML code to paste into the current selection.
		 */
		 function copyImage(src, ids) {
			ajaxPost("/file/copyHttpImage", {src: src}, function(ret) {
				if(reIsOk(ret)) {
					// 将图片替换之
					var src = urlPrefix + "/" + ret.Item;
					var dom = editor.dom
					for(var i in ids) {
						var id = ids[i];
						var imgElm = dom.get(id);
						dom.setAttrib(imgElm, 'src', src);
					}
				}
			});
		}
		// 粘贴HTML
		// 当在pre下时不能粘贴成HTML
		// life add text
		function pasteHtml(html, text) {
			var args, dom = editor.dom;

			// Remove all data images from paste for example from Gecko
			if (!editor.settings.paste_data_images) {
				html = html.replace(/<img[^>]+src=\"data:image[^>]+>/g, '');
			}

			args = editor.fire('BeforePastePreProcess', {content: html}); // Internal event used by Quirks
			args = editor.fire('PastePreProcess', args);
			html = args.content;

			if (!args.isDefaultPrevented()) {
				// User has bound PastePostProcess events then we need to pass it through a DOM node
				// This is not ideal but we don't want to let the browser mess up the HTML for example
				// some browsers add &nbsp; to P tags etc
				if (editor.hasEventListeners('PastePostProcess') && !args.isDefaultPrevented()) {
					// We need to attach the element to the DOM so Sizzle selectors work on the contents
					var tempBody = dom.add(editor.getBody(), 'div', {style: 'display:none'}, html);
					args = editor.fire('PastePostProcess', {node: tempBody});
					dom.remove(tempBody);
					html = args.node.innerHTML;
				}
				
				if (!args.isDefaultPrevented()) {
					// life
					var node = editor.selection.getNode();
					if(node.nodeName == "PRE") {
						if(!text) {
							try {
								text = $(html).text();
							} catch(e) {
							}
						}
						// HTML不能粘贴
						// 其它有错误.... TODO
						// 若有HTML, paste到其它地方有js错误
						// 貼html时自动会删除
						// 纯HTML编辑也会
						text = text.replace(/</g, "&lt;");
						text = text.replace(/>/g, "&gt;");
						// firefox下必须这个
						editor.insertRawContent(text);
						// 之前用insertRawContent()有问题, ace paste下, TODO
						// editor.insertContent(text);
					} else {
						// life 这里得到图片img, 复制到leanote下
						if(!self.copyImage) {
							editor.insertContent(html);
						} else {
							var urlPrefix = UrlPrefix;
							var needCopyImages = {}; // src => [id1,id2]
							var time = (new Date()).getTime();
							try {
								var $html = $("<div>" + html + "</div");
								var $imgs = $html.find("img");
								for(var i = 0; i < $imgs.length; ++i) {
									var $img = $imgs.eq(i)
									var src = $img.attr("src");
									// 是否是外链
									if(src.indexOf(urlPrefix) == -1) {
										time++;
										var id = "__LEANOTE_IMAGE_" + time;
										$img.attr("id", id);
										if(needCopyImages[src]) {
											needCopyImages[src].push(id);
										} else {
											needCopyImages[src] = [id];
										}
									}
								}
								editor.insertContent($html.html());
								
								for(var src in needCopyImages) {
									var ids = needCopyImages[src];
									copyImage(src, ids);
								}
							} catch(e) {
								editor.insertContent(html);
							}
						}
					}
				}
			}
		}

		/**
		 * Pastes the specified text. This means that the plain text is processed
		 * and converted into BR and P elements. It will fire paste events for custom filtering.
		 *
		 * @param {String} text Text to paste as the current selection location.
		 */
		// life text2
		function pasteText(text) {
			var text2 = text;
			text = editor.dom.encode(text).replace(/\r\n/g, '\n');

			var startBlock = editor.dom.getParent(editor.selection.getStart(), editor.dom.isBlock);

			// Create start block html for example <p attr="value">
			var forcedRootBlockName = editor.settings.forced_root_block;
			var forcedRootBlockStartHtml;
			if (forcedRootBlockName) {
				forcedRootBlockStartHtml = editor.dom.createHTML(forcedRootBlockName, editor.settings.forced_root_block_attrs);
				forcedRootBlockStartHtml = forcedRootBlockStartHtml.substr(0, forcedRootBlockStartHtml.length - 3) + '>';
			}

			if ((startBlock && /^(PRE|DIV)$/.test(startBlock.nodeName)) || !forcedRootBlockName) {
				text = Utils.filter(text, [
					[/\n/g, "<br>"]
				]);
			} else {
				text = Utils.filter(text, [
					[/\n\n/g, "</p>" + forcedRootBlockStartHtml],
					[/^(.*<\/p>)(<p>)$/, forcedRootBlockStartHtml + '$1'],
					[/\n/g, "<br />"]
				]);

				if (text.indexOf('<p>') != -1) {
					text = forcedRootBlockStartHtml + text;
				}
			}

			pasteHtml(text, text2);
		}
		
		/**
		 * Creates a paste bin element and moves the selection into that element. It will also move the element offscreen
		 * so that resize handles doesn't get produced on IE or Drag handles or Firefox.
		 */
		function createPasteBin() {
			var dom = editor.dom, body = editor.getBody(), viewport = editor.dom.getViewPort(editor.getWin());
			var scrollY = editor.inline ? body.scrollTop : viewport.y, height = editor.inline ? body.clientHeight : viewport.h;

			removePasteBin();

			// Create a pastebin
			pasteBinElm = dom.add(editor.getBody(), 'div', {
				id: "mcepastebin",
				contentEditable: true,
				"data-mce-bogus": "1",
				style: 'position: absolute; top: ' + (scrollY + 20) + 'px;' +
					'width: 10px; height: ' + (height - 40) + 'px; overflow: hidden; opacity: 0'
			}, pasteBinDefaultContent);

			// Move paste bin out of sight since the controlSelection rect gets displayed otherwise
			dom.setStyle(pasteBinElm, 'left', dom.getStyle(body, 'direction', true) == 'rtl' ? 0xFFFF : -0xFFFF);

			// Prevent focus events from bubbeling fixed FocusManager issues
			dom.bind(pasteBinElm, 'beforedeactivate focusin focusout', function(e) {
				e.stopPropagation();
			});

			lastRng = editor.selection.getRng();
			pasteBinElm.focus();
			editor.selection.select(pasteBinElm, true);
		}

		/**
		 * Removes the paste bin if it exists.
		 */
		function removePasteBin() {
			if (pasteBinElm) {
				editor.dom.unbind(pasteBinElm);
				editor.dom.remove(pasteBinElm);

				if (lastRng) {
					editor.selection.setRng(lastRng);
				}
			}

			keyboardPastePlainTextState = false;
			pasteBinElm = lastRng = null;
		}

		/**
		 * Returns the contents of the paste bin as a HTML string.
		 *
		 * @return {String} Get the contents of the paste bin.
		 */
		function getPasteBinHtml() {
			return pasteBinElm ? pasteBinElm.innerHTML : pasteBinDefaultContent;
		}

		/**
		 * Gets various content types out of the Clipboard API. It will also get the
		 * plain text using older IE and WebKit API:s.
		 *
		 * @param {ClipboardEvent} clipboardEvent Event fired on paste.
		 * @return {Object} Object with mime types and data for those mime types.
		 */
		function getClipboardContent(clipboardEvent) {
			var data = {}, clipboardData = clipboardEvent.clipboardData || editor.getDoc().dataTransfer;

			if (clipboardData && clipboardData.types) {
				data['text/plain'] = clipboardData.getData('Text');

				for (var i = 0; i < clipboardData.types.length; i++) {
					var contentType = clipboardData.types[i];
					data[contentType] = clipboardData.getData(contentType);
				}
			}

			return data;
		}

		function inAcePrevent() {
			// 这个事件是从哪触发的? 浏览器自带的
			// life ace 如果在pre中, 直接返回 TODO
			var ace = LeaAce.nowIsInAce();
			if(ace) {
				// log("in aceEdiotr 2 paste");
				// 原来这里focus了
				setTimeout(function() {
					ace[0].focus();
				});
				return true;
			}
			return false;
		}

		editor.on('keydown', function(e) {
			if (e.isDefaultPrevented()) {
				return;
			}

			// Ctrl+V or Shift+Insert
			if ((VK.metaKeyPressed(e) && e.keyCode == 86) || (e.shiftKey && e.keyCode == 45)) {

				if(inAcePrevent()) {
					return;
				}

				keyboardPastePlainTextState = e.shiftKey && e.keyCode == 86;

				// Prevent undoManager keydown handler from making an undo level with the pastebin in it
				e.stopImmediatePropagation();

				keyboardPasteTimeStamp = new Date().getTime();

				// IE doesn't support Ctrl+Shift+V and it doesn't even produce a paste event
				// so lets fake a paste event and let IE use the execCommand/dataTransfer methods
				if (Env.ie && keyboardPastePlainTextState) {
					e.preventDefault();
					editor.fire('paste', {ieFake: true});
					return;
				}

				createPasteBin();
			}
		});
		
		// 当url改变时, 得到图片的大小 copy from leanote_image
		function getImageSize(url, callback) {
			var img = document.createElement('img');
		
			function done(width, height) {
				img.parentNode.removeChild(img);
				callback({width: width, height: height});
			}
		
			img.onload = function() {
				done(img.clientWidth, img.clientHeight);
			};
		
			img.onerror = function() {
				done();
			};
		
			img.src = url;
		
			var style = img.style;
			style.visibility = 'hidden';
			style.position = 'fixed';
			style.bottom = style.left = 0;
			style.width = style.height = 'auto';
		
			document.body.appendChild(img);
		}
		
		// 上传图片
		function pasteImage(event) {
			// use event.originalEvent.clipboard for newer chrome versions
			  var items = (event.clipboardData  || event.originalEvent.clipboardData).items; // 可能有多个file, 找到属于图片的file
			  log(JSON.stringify(items)); // will give you the mime types
			  // find pasted image among pasted items
			  var blob;
			  for (var i = 0; i < items.length; i++) {
			    if (items[i].type.indexOf("image") === 0) {
			      blob = items[i].getAsFile();
			    }
			  }
			  // load image if there is a pasted image
			  if (blob) {
			    var reader = new FileReader();
			    reader.onload = function(event) {
			      	// 上传之
			      	var c = new FormData;
				    c.append("from", "pasteImage");
				    c.append("file", blob);
				    c.append("noteId", Note.curNoteId); // life
				    // var d;
				    // d = $.ajaxSettings.xhr();
				    // d.withCredentials = i;var d = {};
				    
					// 先显示loading...
					var editor = tinymce.EditorManager.activeEditor; 
					var dom = editor.dom;
					var d = {};						
					d.id = '__mcenew';
					d.src = "http://leanote.com/images/loading-24.gif"; // 写死了
					editor.insertContent(dom.createHTML('img', d));
					var imgElm = dom.get('__mcenew');
				    $.ajax({url: "/file/pasteImage", contentType:false, processData:false , data: c, type: "POST"}
				    	).done(function(re) {
				    		if(!re || typeof re != "object" || !re.Ok) {
				    			// 删除
				    			dom.remove(imgElm);
				    			return;
				    		}
				    		// 这里, 如果图片宽度过大, 这里设置成500px
							var urlPrefix = UrlPrefix; // window.location.protocol + "//" + window.location.host;
							var src = urlPrefix + "/file/outputImage?fileId=" + re.Id;
							getImageSize(src, function(wh) {
								// life 4/25
								if(wh && wh.width) {
									if(wh.width > 600) {
										wh.width = 600;
									}
									d.width = wh.width;
									dom.setAttrib(imgElm, 'width', d.width);
								}
								dom.setAttrib(imgElm, 'src', src);
							});
							dom.setAttrib(imgElm, 'id', null);
				    	});
			    };
			    reader.readAsDataURL(blob);
			    return true;
			  }
			  return false;
		}

		editor.on('paste', function(e) {
			if(inAcePrevent()) {
				return;
			}

			var clipboardContent = getClipboardContent(e);
			var isKeyBoardPaste = new Date().getTime() - keyboardPasteTimeStamp < 100;
			var plainTextMode = self.pasteFormat == "text" || keyboardPastePlainTextState;

			// Not a keyboard paste prevent default paste and try to grab the clipboard contents using different APIs
			if (!isKeyBoardPaste) {
				e.preventDefault();
			}

			// Try IE only method if paste isn't a keyboard paste
			if (Env.ie && (!isKeyBoardPaste || e.ieFake)) {
				createPasteBin();

				editor.dom.bind(pasteBinElm, 'paste', function(e) {
					e.stopPropagation();
				});

				editor.getDoc().execCommand('Paste', false, null);
				clipboardContent["text/html"] = getPasteBinHtml();
				removePasteBin();
			}

			setTimeout(function() {
				var html = getPasteBinHtml();

				// WebKit has a nice bug where it clones the paste bin if you paste from for example notepad
				if (pasteBinElm && pasteBinElm.firstChild && pasteBinElm.firstChild.id === 'mcepastebin') {
					plainTextMode = true;
				}

				removePasteBin();

				if (html == pasteBinDefaultContent || !isKeyBoardPaste) {
					html = clipboardContent['text/html'] || clipboardContent['text/plain'] || pasteBinDefaultContent;

					if (html == pasteBinDefaultContent) {
						if (!isKeyBoardPaste) {
							// editor.windowManager.alert('Please use Ctrl+V/Cmd+V keyboard shortcuts to paste contents.');
						}
						return;
					}
				}

				if (plainTextMode) {
					pasteText(clipboardContent['text/plain'] || Utils.innerText(html));
				} else {
					// life
					pasteHtml(html, clipboardContent['text/plain']);
				}
			}, 0);
			
			//-----------
			// paste image
			try {
				/*
				if(pasteImage(e)) {
					return;
				}
				*/
			} catch(e) {};

		});
		
		

		self.pasteHtml = pasteHtml;
		self.pasteText = pasteText;
	};
});
