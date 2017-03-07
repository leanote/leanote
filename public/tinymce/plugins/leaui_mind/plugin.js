/**
 * leaui album image manager plugin
 * copyright leaui
 * leaui.com
 */
var LEAUI_MIND = {};
tinymce.PluginManager.add('leaui_mind', function(editor, url) {
	
	function showDialog() {
		var dom = editor.dom;

		var content = editor.selection.getContent();
		// get images and attrs
		var p = /<img.*?\/>/g;
		var images = content.match(p);
		var newNode = document.createElement("p");
		LEAUI_MIND = {};
		for(var i in images) {
			newNode.innerHTML = images[i];
			var imgElm = newNode.firstChild;
			if(imgElm && imgElm.nodeName == "IMG") {
				LEAUI_MIND.json = dom.getAttrib(imgElm, 'data-mind-json');
				break;
			}
		}

		function GetTheHtml(){
			var lang = editor.settings.language;
			var u = '//leanote.com/public/libs/mind/edit.html';
			// u = 'http://localhost:9000/public/libs/mind/edit.html';
			var html = '<iframe id="leauiIfr" src="'+ u + '?' + new Date().getTime() + '&lang=' + lang + '" frameborder="0"></iframe>';
			return html;
		}
		
		var w = $(document).width() - 10;
		var h = $(document).height() - 100;
		
		win = editor.windowManager.open({
			title: "Mind Map",
			width : w,
			height : h,
			html: GetTheHtml(),
			buttons: [
				{
					text: 'Cancel',
					onclick: function() {
						this.parent().parent().close();
					}
				},
				{
				text: 'Insert Mind Map',
				subtype: 'primary',
				onclick: function(e) {
					var me = this;
					var _iframe = document.getElementById('leauiIfr').contentWindow;
					var km = _iframe.km;
					km.exportData('png').then(function(data) {
						var json = JSON.stringify(km.exportJson());
						// console.log(json);
						var img = '<img src="' + data + '" data-mce-src="-" data-mind-json=\'' + json + '\'>';
						editor.insertContent(img);
						
						me.parent().parent().close();
					});
					return;
					
					var _div =_iframe.document.getElementById('preview');
					var ii = _div.childNodes; 
					//console.log(ii);
					var datas = [];
					for(var i = 0; i < ii.length; ++i) {
						var e = ii[i]; 
						//console.log(e);
						// 有些没有image
						if(e.firstChild && e.firstChild.nodeName == "IMG") {
							var img = e.firstChild;
							var d = {};
							d.src = img.getAttribute("src");
							d.width = img.getAttribute("data-width");
							d.height = img.getAttribute("data-height");
							d.title = img.getAttribute("data-title");

							datas.push(d);
						}
					};

					for(var i in datas) {
						var data = datas[i];
						var src = data.src;
						// the network image
						var trueSrc;
						if(src.indexOf("http://") != -1 || src.indexOf("https://") != -1) {
							trueSrc = src;
						} else {
							trueSrc = url + "/" + src;
						}
						data.src = trueSrc;
						
						var renderImage = function(data) {
							// 这里, 如果图片宽度过大, 这里设置成500px
							var back = (function(data2, i) {
								var d = {};
								var imgElm;
								// 先显示loading...
								d.id = '__mcenew' + i;
								d.src = "http://leanote.com/images/loading-24.gif";
								imgElm = dom.createHTML('img', d);
								editor.insertContent(imgElm);
								imgElm = dom.get(d.id);
								
								return function(wh) {
									if(wh && wh.width) {
										if(wh.width > 600) {
											wh.width = 600;
										}
										data2.width = wh.width;
									}
									dom.setAttrib(imgElm, 'src', data2.src);
									// dom.setAttrib(imgElm, 'width', data2.width);
									dom.setAttrib(imgElm, 'title', data2.title);
									
									dom.setAttrib(imgElm, 'id', null);
								}
							})(data, i);
							getImageSize(data.src, back);
						}
						
						// outputImage?fileId=123232323
						var fileId = "";
						fileIds = trueSrc.split("fileId=")
						if(fileIds.length == 2 && fileIds[1].length == "53aecf8a8a039a43c8036282".length) {
							fileId = fileIds[1];
						}
						if(fileId) {
							// 得到fileId, 如果这个笔记不是我的, 那么肯定是协作的笔记, 那么需要将图片copy给原note owner
							// 博客设置中不用没有Note
							var curNote;
							if(Note && Note.getCurNote) {
								curNote = Note.getCurNote();
							}
							if(curNote && curNote.UserId != UserInfo.UserId) {
								(function(data) {
									ajaxPost("/file/copyImage", {userId: UserInfo.UserId, fileId: fileId, toUserId: curNote.UserId}, function(re) {
										if(reIsOk(re) && re.Id) {
											var urlPrefix = UrlPrefix; // window.location.protocol + "//" + window.location.host;
											data.src = urlPrefix + "/file/outputImage?fileId=" + re.Id;
										}
										renderImage(data);
									});
								})(data);
							} else {
								renderImage(data);
							}
						} else {
							renderImage(data);
						}
						
					} // end for
					
					this.parent().parent().close();
				}
				}]
		});
	}
	
	editor.addButton('leaui_mind', {
		icon: 'mind',
		tooltip: 'Insert/edit mind map',
		onclick: showDialog,
		stateSelector: 'img[data-mind-json]'
	});
});
