
/**
 * http://leanote.com
 */

var tinymce = top.tinymce;
var editor = tinymce.EditorManager.activeEditor; 
var dom = editor.dom;
var imgElm = editor.selection.getNode();

$(function() {
	// 隐藏iframe border
//	top.hiddenIframeBorder();
	
	// bind event
	$("#uploader").change(function() {
		$("#upl").submit(); 
		leanoteImagesDialog.inProgress();
	});
	
	$("#insertImageBtn").click(function() {
		insertImage();
	});
	$("#closeBtn").click(function() {
		closeWin();
	});
	
	// 是否选择了image
	
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
		top.closeDialog();
//		editor.windowManager.close();
//		editor.windowManager.close();	
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

// 使用input上传
var leanoteImagesDialog = {
	timeoutStore : false,
	
	// 上传时
	inProgress : function() {
		$("#uploadFormContainer").hide();
		$("#uploadInProgress").show();
		// 可能是返回的数据不对
		this.timeoutStore = window.setTimeout(function() {
			$("#uploadInProgress").html("服务器发生错误, 超时.");
		}, 20000);
	},
	
	// 上传结束后, 关闭窗口, 填充image
	uploadFinish : function(result) {
		if (result.resultCode == '0') {
			window.clearTimeout(this.timeoutStore);
			$("#uploadInProgress").html("服务器发生错误: " + result.result);
		} else {
			$("#uploadInProgress").html("上传成功!");
			
			// 这里, 如果图片宽度过大, 这里设置成500px
			var d = {};
			var imgElm;
			getImageSize(result.filename, function(wh) {
				// life 4/25
				if(wh && wh.width) {
					if(wh.width > 600) {
						wh.width = 600;
					}
					d.width = wh.width;
					dom.setAttrib(imgElm, 'width', d.width);
				}
				dom.setAttrib(imgElm, 'src', result.filename);
			});
			// 先显示loading...
			d.id = '__mcenew';
			d.src = "http://leanote.com/images/loading-24.gif";
			editor.insertContent(dom.createHTML('img', d));
			imgElm = dom.get('__mcenew');
			dom.setAttrib(imgElm, 'id', null);
			closeWin();
		}
	}
};