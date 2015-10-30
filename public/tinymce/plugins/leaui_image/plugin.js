/**
 * leaui album image manager plugin
 * copyright leaui
 * leaui.com
 */
var LEAUI_DATAS = [];
tinymce.PluginManager.add('leaui_image', function(editor, url) {
	//当url改变时, 得到图片的大小
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

	function showDialog() {
		var dom = editor.dom;

		var content = editor.selection.getContent();
		// get images and attrs
		var p = /<img.*?\/>/g;
		var images = content.match(p);
		var newNode = document.createElement("p");
		var datas = [];
		for(var i in images) {
			newNode.innerHTML = images[i];
			var imgElm = newNode.firstChild;
			if(imgElm && imgElm.nodeName == "IMG") {
				var data = {};
				data.src = dom.getAttrib(imgElm, 'data-src') || dom.getAttrib(imgElm, 'src');
				data.width = dom.getAttrib(imgElm, 'width');
				data.height = dom.getAttrib(imgElm, 'height');
				data.title = dom.getAttrib(imgElm, 'title');
				datas.push(data);
			}
		}
		LEAUI_DATAS = datas;

		function GetTheHtml(){
			var html = '<iframe id="leauiIfr" src="/album/index'+ '?' + new Date().getTime() + '" frameborder="0"></iframe>';
			return html;
		}

		var w = $(document).width() - 10;
		if(w > 805) {
			w = 805;
		}
		var h = $(document).height() - 100;
		if(h > 365) {
			h = 365;
		}
		win = editor.windowManager.open({
			title: "Image",
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
				text: 'Insert Image',
				subtype: 'primary',
				onclick: function(e) {
					var _iframe = document.getElementById('leauiIfr').contentWindow;
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
											data.src = urlPrefix + "/api/file/getImage?fileId=" + re.Id;
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

	editor.addButton('leaui_image', {
		icon: 'image',
		tooltip: 'Insert/edit image',
		onclick: showDialog,
		stateSelector: 'img:not([data-mind-json])'
	});

	editor.addMenuItem('leaui_image', {
		icon: 'image',
		text: 'Insert image',
		onclick: showDialog,
		context: 'insert',
		prependToContext: true
	});

	// 为解决在editor里拖动图片问题
	// 2014/7/8 21:43 浮躁的一天终有收获
	// 2015/10/16
	// TODO 如果把编辑器内的图片拖到外面去, 还是会出现drop images to here
    var dragStart = false;
    editor.on("dragstart", function(e) {
    	// readonly时不让drag图片
    	if (LEA.readOnly) {
	    	e.preventDefault();
	    	e.stopPropagation();
    	}
    	dragStart = true;
    });
    editor.on("dragend", function(e) {
    	dragStart = false;
    });
	editor.on("dragover", function(e) {
	    if(dragStart) {
    		// 表示编辑器内在拖动图片, 则停止冒泡
    		e.preventDefault();
	    	e.stopPropagation();
    	}
    });
});
