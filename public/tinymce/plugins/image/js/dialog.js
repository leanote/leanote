/**
 * life
 */

var tinymce = top.tinymce;
var editor = tinymce.EditorManager.activeEditor; 
var dom = editor.dom;
var imgElm = editor.selection.getNode();

var jbImagesDialog = {
	resized : false,
	iframeOpened : false,
	timeoutStore : false,
	
	inProgress : function() {
		document.getElementById("upload_infobar").style.display = 'none';
		document.getElementById("upload_additional_info").innerHTML = '';
		document.getElementById("upload_form_container").style.display = 'none';
		document.getElementById("upload_in_progress").style.display = 'block';
		this.timeoutStore = window.setTimeout(function(){
			document.getElementById("upload_additional_info").innerHTML = 'This is taking longer than usual.' + '<br />' + 'An error may have occurred.' + '<br /><a href="#" onClick="jbImagesDialog.showIframe()">' + 'View script\'s output' + '</a>';
			// tinyMCEPopup.editor.windowManager.resizeBy(0, 30, tinyMCEPopup.id);
		}, 20000);
	},
	
	showIframe : function() {
		if (this.iframeOpened == false)
		{
			document.getElementById("upload_target").className = 'upload_target_visible';
			// tinyMCEPopup.editor.windowManager.resizeBy(0, 190, tinyMCEPopup.id);
			this.iframeOpened = true;
		}
	},
	
	uploadFinish : function(result) {
		if (result.resultCode == '0')
		{
			window.clearTimeout(this.timeoutStore);
			document.getElementById("upload_in_progress").style.display = 'none';
			document.getElementById("upload_infobar").style.display = 'block';
			$("#upload_infobar").html(result.result).show();
			document.getElementById("upload_form_container").style.display = 'block';
			
			if (this.resized == false)
			{
				// tinyMCEPopup.editor.windowManager.resizeBy(0, 30, tinyMCEPopup.id);
				this.resized = true;
			}
		}
		else
		{
			document.getElementById("upload_in_progress").style.display = 'none';
			document.getElementById("upload_infobar").style.display = 'block';
			document.getElementById("upload_infobar").innerHTML = 'Upload Complete';
			
			var w = this.getWin();
			tinymce = w.tinymce;
		
			tinymce.EditorManager.activeEditor.insertContent('<img src="' + result.filename +'">');
			
			this.close();

			// 添加undo
			// var editor = tinymce.activeEditor;
			// editor.undoManager.transact(function() {
			// });
		}
	},
	
	getWin : function() {
		return (!window.frameElement && window.dialogArguments) || opener || parent || top;
	},
	
	close : function() {
		var t = this;

		// To avoid domain relaxing issue in Opera
		function close() {
			tinymce.EditorManager.activeEditor.windowManager.close(window);
			tinymce = tinyMCE = t.editor = t.params = t.dom = t.dom.doc = null; // Cleanup
		};

		if (tinymce.isOpera)
			this.getWin().setTimeout(close, 0);
		else
			close();
	}

};

//----------------------- url

$(function() {
	top.hiddenIframeBorder();
	
	var oldWidth, oldHeihgt;
	// 是否选择的是image
	if(imgElm.nodeName == "IMG") {
		var $node = $(imgElm);
		$("#imageSrc").val($node.attr("src"));
		oldWidth = $node.width();
		oldHeight = $node.height();
		$("#imageWidth").val(oldWidth);
		$("#imageHeight").val(oldHeight);
		
		$('#myTab a:last').tab('show');
	} else {
		imgElm = null;
	}
	
	$("#imageSrc").blur(function(){
		getImageSize($(this).val(), function(ret) {
			if(ret.width) {
				oldWidth = ret.width;
				oldHeight = ret.Height;
			
				$("#imageWidth").val(ret.width);
				$("#imageHeight").val(ret.height);
			}
		});
	});
	
	// 按比例缩放
	function scale(isWidth) {
		var autoScale = $("#autoScale").is(":checked");
		var width = $("#imageWidth").val();
		var height = $("#imageHeight").val();
		
		if(autoScale && oldWidth && oldHeight) {
			if(isWidth) {
				height = parseInt((width/oldWidth) * oldHeight);
				$("#imageHegiht").val(height);
			} else {
				width = parseInt((height/oldHeight) * oldWidth);
				$("#imageWidth").val(width);
			}
		}
		
		oldWidth = width;
		oldHeight = height;	
	}
	
	$("#imageWidth").blur(function() {
		scale(true);
	});
	$("#imageHeight").blur(function() {
		scale(false);
	});
});

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

function closeWin() {
	try {
		editor.windowManager.close();
		editor.windowManager.close();	
	} catch(e) {
		
	}
}

// 插入之
var insertImage = function() {
	// 判断是否是第二个tab active
	if($("#myTab li:last").attr("class") != "active") {
		closeWin();
		return;
	}
	
	// 加载图片并插入之
	function waitLoad(imgElm) {
		function selectImage() {
			imgElm.onload = imgElm.onerror = null;
			editor.selection.select(imgElm);
			editor.nodeChanged();
		}

		// 如果没有设置width, height, 就用图片本身的大小
		imgElm.onload = function() {
			if (!data.width && !data.height) {
				dom.setAttribs(imgElm, {
					width: imgElm.clientWidth,
					height: imgElm.clientHeight
				});
			}

			selectImage();
		};

		imgElm.onerror = selectImage;
	}

	// 这是通过url插入图片
	// iframe里得到...
	var data = {width:null, height:null, src:null, style:null};
	data.width = $("#imageWidth").val();
	data.height = $("#imageHeight").val();
	data.src = $("#imageSrc").val();
	data.alt = "";

	editor.undoManager.transact(function() {
		// 删除图片
		if (!data.src) {
			if (imgElm) {
				dom.remove(imgElm);
				editor.nodeChanged();
			}

			return;
		}
		if (!imgElm) {
			data.id = '__mcenew';
			editor.selection.setContent(dom.createHTML('img', data));
			imgElm = dom.get('__mcenew');
			dom.setAttrib(imgElm, 'id', null);
		} else {
			dom.setAttribs(imgElm, data);
		}

		waitLoad(imgElm);
	});
	
	closeWin();
}