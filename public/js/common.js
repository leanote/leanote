// leanote 通用方法

//--------------
// 命名空间
//--------------

// 最上级变量
var LEA = {};
// 命名空间
var Notebook = {
	cache: {}, // notebookId => {Title, Seq}
}
var Note = {
	cache: {}, // noteId => {Title, Tags, Content, Desc}
};
// var UserInfo = {}; // 博客有问题, 会覆盖
var Tag = {};
var Notebook = {};
var Share = {};

// markdown
var Converter;
var MarkdownEditor;
var ScrollLink;

//---------------------
// 公用方法

function trimLeft(str, substr) {
	if(!substr || substr == " ") {
		return $.trim(str);
	}
	while(str.indexOf(substr) == 0) {
		str = str.substring(substr.length);
	}
	return str;
}

function json(str) {
	return eval("(" + str + ")")
}

// '<div id="?" class="?" onclick="?">'
function t() {
	var args = arguments;
	if(args.length <= 1) {
		return args[0];
	}
	var text = args[0];
	if(!text) {
		return text;
	}
	
	// 先把所有的?替换成, 很有可能替换的值有?会造成循环,没有替换想要的
	var pattern = "LEAAEL"
	text = text.replace(/\?/g, pattern);
	
	// args[1] 替换第一个?
	for(var i = 1; i <= args.length; ++i) {
		text = text.replace(pattern, args[i]);
	}
	return text;
}

// 判断数组是否相等
function arrayEqual(a, b) {
	a = a || [];
	b = b || [];
	return a.join(",") == b.join(",");
}

// 是否是数组
function isArray(obj) {  
	return Object.prototype.toString.call(obj) === '[object Array]';   
}

/**
 * 是否为空
 * 可判断任意类型，string array
 */
function isEmpty(obj) {
	if(!obj) {
		return true;
	}
	
	if(isArray(obj)) {
		if(obj.length == 0) {
			return true;
		}
	}
	
	return false;
}

//------------
//得到form的数据
//返回json
function getFormJsonData(formId) {
	var data = formArrDataToJson($('#' + formId).serializeArray());
	return data;
}

//$('#form').serializeArray()的数据[{name: a, value: b}, {name: "c[]", value: d}]
//转成{a:b}
function formArrDataToJson(arrData) {
	var datas = {};
	var arrObj= {}; // {a:[1, 2], b:[2, 3]};
	for(var i in arrData) {
		var attr = arrData[i].name;
		var value = arrData[i].value;
		// 判断是否是a[]形式
		if(attr.substring(attr.length-2, attr.length) == '[]') {
			attr = attr.substring(0, attr.length-2);
			if(arrObj[attr] == undefined) {
				arrObj[attr] = [value];
			} else {
				arrObj[attr].push(value);
			}
			continue;
		}
		
		datas[attr] = value;
	}
	
	return $.extend(datas, arrObj);
}

//将serialize的的form值转成json
function formSerializeDataToJson(formSerializeData) {
	var arr = formSerializeData.split("&");
	var datas = {};
	var arrObj= {}; // {a:[1, 2], b:[2, 3]};
	for(var i = 0; i < arr.length; ++i) {
		var each = arr[i].split("=");
		var attr = decodeURI(each[0]);
		var value = decodeURI(each[1]);
		// 判断是否是a[]形式
		if(attr.substring(attr.length-2, attr.length) == '[]') {
			attr = attr.substring(0, attr.length-2);
			if(arrObj[attr] == undefined) {
				arrObj[attr] = [value];
			} else {
				arrObj[attr].push(value);
			}
			continue;
		}
		datas[attr] = value;
	}
	
	return $.extend(datas, arrObj);
}


// ajax请求返回结果后的操作
// 用于ajaxGet(), ajaxPost()
function _ajaxCallback(ret, successFunc, failureFunc) {
	// 总会执行
	if(ret === true || ret == "true" || typeof ret == "object") {
		// 是否是NOTELOGIN
		if(ret && typeof ret == "object") {
			if(ret.Msg == "NOTLOGIN") {
				alert("你还没有登录, 请先登录!");
				return;
			}
		}
		if(typeof successFunc == "function") {
			successFunc(ret);
		}
	} else {
		if(typeof failureFunc == "function") {
			failureFunc(ret);
		} else {
			alert("error!")
		}
	}
}
function _ajax(type, url, param, successFunc, failureFunc, async) {
	log("-------------------ajax:");
	log(url);
	log(param);
	if(typeof async == "undefined") {
		async = true;
	} else {
		async = false;
	}
	$.ajax({
		type: type,
		url: url,
		data: param,
		async: async, // 是否异步
		success: function(ret) {
			_ajaxCallback(ret, successFunc, failureFunc);
		},
		error: function(ret) {
			_ajaxCallback(ret, successFunc, failureFunc);
		}
	});
}

/**
 * 发送ajax get请求
 * @param url
 * @param param
 * @param successFunc
 * @param failureFunc
 * @param hasProgress
 * @param async 是否异步
 * @returns
 */
function ajaxGet(url, param, successFunc, failureFunc, async) {
	_ajax("GET", url, param, successFunc, failureFunc, async);
}

/**
 * 发送post请求
 * @param url
 * @param param
 * @param successFunc
 * @param failureFunc
 * @param hasProgress
 * @param async 是否异步, 默认为true
 * @returns
 */
function ajaxPost(url, param, successFunc, failureFunc, async) {
	_ajax("POST", url, param, successFunc, failureFunc, async);
}
function ajaxPostJson(url, param, successFunc, failureFunc, async) {
	log("-------------------ajaxPostJson:");
	log(url);
	log(param);
	
	// 默认是异步的
	if(typeof async == "undefined") {
		async = true;
	} else {
		async = false;
	}
	$.ajax({
	    url : url,
	    type : "POST",
	    contentType: "application/json; charset=utf-8",
	    datatype: "json",
	    async: async,
	    data : JSON.stringify(param),
	    success : function(ret, stats) {
			_ajaxCallback(ret, successFunc, failureFunc);
	    },
		error: function(ret) {
			_ajaxCallback(ret, successFunc, failureFunc);
		}
	});
}

function findParents(target, selector) {
	if($(target).is(selector)) {
		return $(target);
	}
	var parents = $(target).parents();
	for(var i = 0; i < parents.length; ++i) {
		log(parents.eq(i))
		if(parents.eq(i).is(selector)) {
			return parents.seq(i);
		}
	}
	return null;
}

/*
ajaxPostJson(
	"http://localhost:9000/notebook/index?i=100&name=life", 
	{Title: "you can",  UserId:"52a9e409f4ea49d6576fdbca", Subs:[{title: "xxxxx", Seq:11}, {title:"life..."}]}, 
	function(e) {
		log(e);
	});
*/

//-----------------

// 切换编辑器时要修改tabIndex
function editorIframeTabindex(index) {
	var $i = $("#editorContent_ifr");
	if($i.size() == 0) {
		setTimeout(function() {
			editorIframeTabindex(index);
		}, 100);
	} else {
		$i.attr("tabindex", index);
	}
}
//切换编辑器
function switchEditor(isMarkdown) {
	// 富文本永远是2
	if(!isMarkdown) {
		$("#editor").show();
		$("#mdEditor").css("z-index", 1);
		
		// 刚开始没有
		editorIframeTabindex(2);
		$("#wmd-input").attr("tabindex", 3);
		$("#leanoteNav").show();
	} else {
		$("#mdEditor").css("z-index", 3).show();
		
		editorIframeTabindex(3);
		$("#wmd-input").attr("tabindex", 2);
		$("#leanoteNav").hide();
	}
}

// editor 设置内容
// 可能是tinymce还没有渲染成功
var previewToken = "<div style='display: none'>FORTOKEN</div>"
function setEditorContent(content, isMarkdown, preview) {
	if(!content) {
		content = "";
	}
	if(!isMarkdown) {
		$("#editorContent").html(content);
		var editor = tinymce.activeEditor;
		if(editor) {
			editor.setContent(content);
			editor.undoManager.clear(); // 4-7修复BUG
		} else {
			// 等下再设置
			setTimeout(function() {
				setEditorContent(content, false);
			}, 100);
		}
	} else {
		$("#wmd-input").val(content);
		$("#wmd-preview").html(""); // 防止先点有的, 再点tinymce再点没内容的
		if(!content || preview) { // 没有内容就不要解析了
			$("#wmd-preview").html(preview).css("height", "auto");
			if(ScrollLink) {
				ScrollLink.onPreviewFinished(); // 告诉scroll preview结束了
			}
		} else {
			// 还要清空preview
			if(MarkdownEditor) {
				$("#wmd-preview").html(previewToken + "<div style='text-align:center; padding: 10px 0;'><img src='http://leanote.com/images/loading-24.gif' /> 正在转换...</div>");
				MarkdownEditor.refreshPreview();
			} else {
				// 等下再设置
				setTimeout(function() {
					setEditorContent(content, true, preview);
				}, 200);
			}
		}
	}
}

// preview是否为空
function previewIsEmpty(preview) {
	if(!preview || preview.substr(0, previewToken.length) == previewToken) {
		return true;
	}
	return false;
}

// 有tinymce得到的content有<html>包围
function getEditorContent(isMarkdown) {
	if(!isMarkdown) {
		var editor = tinymce.activeEditor;
		if(editor) {
			var content = $(editor.getBody());
			// 去掉恶心的花瓣注入
			// <pinit></pinit>
			// 把最后的<script>..</script>全去掉
			content.find("pinit").remove();
			content.find(".thunderpin").remove();
			content.find(".pin").parent().remove();
			content = $(content).html();
			if(content) {
				while(true) {
					var lastEndScriptPos = content.lastIndexOf("</script>");
					if (lastEndScriptPos == -1) {
						return content;
					}
					var length = content.length;
					// 证明</script>在最后, 去除之
					if(length - 9 == lastEndScriptPos) {
						var lastScriptPos = content.lastIndexOf("<script ");
						if(lastScriptPos == -1) {
							lastScriptPos = content.lastIndexOf("<script>");
						}
						if(lastScriptPos != -1) {
							content = content.substring(0, lastScriptPos);
						} else {
							return content;
						}
					} else {
						// 不在最后, 返回
						return content;
					}
				}
			}
			return content;
		}
	} else {
		return [$("#wmd-input").val(), $("#wmd-preview").html()]
	}
}

// 禁用编辑
LEA.editorStatus = true;
function disableEditor() {
	var editor = tinymce.activeEditor;
	if(editor) {
		editor.hide();
		LEA.editorStatus = false;
		$("#mceTollbarMark").show().css("z-index", 1000);
	}
	
	// toolbar 来个遮着...
}
function enableEditor() {
	if(LEA.editorStatus) {
		return;
	}
	$("#mceTollbarMark").css("z-index", -1).hide();
	var editor = tinymce.activeEditor;
	if(editor) {
		editor.show();
	}
}

//-----------
// dialog
//-----------
function showDialog(id, options) {
	$("#leanoteDialog #modalTitle").html(options.title);
	$("#leanoteDialog .modal-body").html($("#" + id + " .modal-body").html());
	$("#leanoteDialog .modal-footer").html($("#" + id + " .modal-footer").html());
	delete options.title;
	options.show = true;
	$("#leanoteDialog").modal(options);
}
function hideDialog(timeout) {
	if(!timeout) {
		timeout = 0;
	}
	setTimeout(function() {
		$("#leanoteDialog").modal('hide');
	}, timeout);
}

// 更通用
function closeDialog() {
	$(".modal").modal('hide');
}

// 原生的
function showDialog2(id, options) {
	options = options || {};
	options.show = true;
	$(id).modal(options);
}
function hideDialog2(id, timeout) {
	if(!timeout) {
		timeout = 0;
	}
	setTimeout(function() {
		$(id).modal('hide');
	}, timeout);
}

// 远程
function showDialogRemote(url, data) {
	data = data || {};
	url += "?";
	for(var i in data) {
		url += i + "=" + data[i] + "&";
	}
	$("#leanoteDialogRemote").modal({remote: url});
}

function hideDialogRemote() {
	$("#leanoteDialogRemote").modal('hide');
}
//---------------
// notify
// 没用
$(function() {
	if($.pnotify) {
		$.pnotify.defaults.delay = 1000;
	}
})

function notifyInfo(text) {
	$.pnotify({
	    title: '通知',
	    text: text,
	    type: 'info',
	    styling: 'bootstrap'
	});
}
function notifyError(text) {
	$.pnotify.defaults.delay = 2000
	$.pnotify({
	    title: '通知',
	    text: text,
	    type: 'error',
	    styling: 'bootstrap'
	});
}
function notifySuccess(text) {
	$.pnotify({
	    title: '通知',
	    text: text,
	    type: 'success',
	    styling: 'bootstrap'
	});
}

// 对Date的扩展，将 Date 转化为指定格式的String   
//月(M)、日(d)、小时(h)、分(m)、秒(s)、季度(q) 可以用 1-2 个占位符，   
//年(y)可以用 1-4 个占位符，毫秒(S)只能用 1 个占位符(是 1-3 位的数字)   
//例子：   
//(new Date()).format("yyyy-MM-dd hh:mm:ss.S") ==> 2006-07-02 08:09:04.423   
//(new Date()).format("yyyy-M-d h:m:s.S")      ==> 2006-7-2 8:9:4.18
Date.prototype.format = function(fmt) { //author: meizz   
  var o = {   
    "M+" : this.getMonth()+1,                 //月份   
    "d+" : this.getDate(),                    //日   
    "h+" : this.getHours(),                   //小时   
    "m+" : this.getMinutes(),                 //分   
    "s+" : this.getSeconds(),                 //秒   
    "q+" : Math.floor((this.getMonth()+3)/3), //季度   
    "S"  : this.getMilliseconds()             //毫秒   
  };   
  if(/(y+)/.test(fmt))   
    fmt=fmt.replace(RegExp.$1, (this.getFullYear()+"").substr(4 - RegExp.$1.length));   
  for(var k in o)   
    if(new RegExp("("+ k +")").test(fmt))   
  fmt = fmt.replace(RegExp.$1, (RegExp.$1.length==1) ? (o[k]) : (("00"+ o[k]).substr((""+ o[k]).length)));   
  return fmt; 
}

//2014-01-06T18:29:48.802+08:00
function goNowToDatetime(goNow) {
	if(!goNow) {
		return "";
	}
	return goNow.substr(0, 10) + " " + goNow.substr(11, 8);
}
function getCurDate() {
	return (new Date()).format("yyyy-M-d");
}

// 回车键的动作
function enter(parent, children, func) {
	if(!parent) {
		parent = "body";
	}
	$(parent).on("keydown", children, function(e) {
		if (e.keyCode == 13) {
			func.call(this);
		}
	});
}

// 回车则blue
function enterBlur(parent, children) {
	if(!parent) {
		parent = "body";
	}
	if(!children) {
		children = parent;
		parent = "body";
	}
	$(parent).on("keydown", children, function(e) {
		if (e.keyCode == 13) {
			$(this).trigger("blur");
		}
	});
}

// 生成mongodb ObjectId
function getObjectId() {
	return ObjectId();
}

//-----------------------------------------
function resizeEditor(second) {
	var ifrParent = $("#editorContent_ifr").parent();
    ifrParent.css("overflow", "auto");
    var height = $("#editorContent").height();
    ifrParent.height(height);
    // log(height + '---------------------------------------')
    $("#editorContent_ifr").height(height);
    
    /*
    // 第一次时可能会被改变
    if(!second) {
		setTimeout(function() {
			resizeEditorHeight(true);
		}, 1000);
    }
    */
}

//----------
// msg位置固定
function showMsg(msg, timeout) {
	$("#msg").html(msg);
	if(timeout) {
		setTimeout(function() {
			$("#msg").html("");
		}, timeout)
	}
}
function showMsg2(id, msg, timeout) {
	$(id).html(msg);
	if(timeout) {
		setTimeout(function() {
			$(id).html("");
		}, timeout)
	}
}

//--------------
// type == danger, success, warning
function showAlert(id, msg, type, id2Focus) {
	$(id).html(msg).removeClass("alert-danger").removeClass("alert-success").removeClass("alert-warning").addClass("alert-" + type).show();
	if(id2Focus) {
		$(id2Focus).focus();
	}
}
function hideAlert(id, timeout) {
	if(timeout) {
		setTimeout(function() {
			$(id).hide();
		}, timeout);
	} else {
		$(id).hide();
	}
}

//-------------------
// for leanote ajax

// post
// return {Ok, Msg, Data}
// btnId 是按钮包括#
function post(url, param, func, btnId) {
	var btnPreText;
	if(btnId) {
		btnPreText = $(btnId).html();
		$(btnId).html("正在处理").addClass("disabled");
	}
	ajaxPost(url, param, function(ret) {
		if(btnPreText) {
			$(btnId).html(btnPreText).removeClass("disabled");
		}
		if (typeof ret == "object") {
			if(typeof func == "function") {
				func(ret);
			}
		} else {
			alert("leanote出现了错误!");
		}
	});
}

// 是否是正确的email
function isEmail(email) {
	var myreg = /^([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+@([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+\.[0-9a-zA-Z]{2,3}$/;
	return myreg.test(email);
}

// 正确返回该email
function isEmailFromInput(inputId, msgId, selfBlankMsg, selfInvalidMsg) {
	var val = $(inputId).val();
	var msg = function() {};
	if(msgId) {
		msg = function(msgId, msg) {
			showAlert(msgId, msg, "danger", inputId);
		}
	}
	if(!val) {
		msg(msgId, selfBlankMsg || "请输入邮箱");
	} else if(!isEmail(val)) {
		msg(msgId, selfInvalidMsg || "请输入正确的邮箱");
	} else {
		return val;
	}
}

// 复制文本
function initCopy(aId, postFunc) {
	// 定义一个新的复制对象
	var clip = new ZeroClipboard(document.getElementById(aId), {
	  moviePath: "/js/ZeroClipboard/ZeroClipboard.swf"
	});

	// 复制内容到剪贴板成功后的操作
	clip.on('complete', function(client, args) {
		postFunc(args);
	});   
}

function showLoading() {
	$("#loading").css("visibility", "visible");
}

function hideLoading() {
	$("#loading").css("visibility", "hidden");
}

// 注销, 先清空cookie
function logout() {
	$.removeCookie("LEANOTE_SESSION");
	location.href = "/logout?id=1";
}

// 得到图片width, height, callback(ret); ret = {width:11, height:33}
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

// 插件中使用
function hiddenIframeBorder() {
	$('.mce-window iframe').attr("frameborder", "no").attr("scrolling", "no");
}

var email2LoginAddress = {
    'qq.com': 'http://mail.qq.com',
    'gmail.com': 'http://mail.google.com',
    'sina.com': 'http://mail.sina.com.cn',
    '163.com': 'http://mail.163.com',
    '126.com': 'http://mail.126.com',
    'yeah.net': 'http://www.yeah.net/',
    'sohu.com': 'http://mail.sohu.com/',
    'tom.com': 'http://mail.tom.com/',
    'sogou.com': 'http://mail.sogou.com/',
    '139.com': 'http://mail.10086.cn/',
    'hotmail.com': 'http://www.hotmail.com',
    'live.com': 'http://login.live.com/',
    'live.cn': 'http://login.live.cn/',
    'live.com.cn': 'http://login.live.com.cn',
    '189.com': 'http://webmail16.189.cn/webmail/',
    'yahoo.com.cn': 'http://mail.cn.yahoo.com/',
    'yahoo.cn': 'http://mail.cn.yahoo.com/',
    'eyou.com': 'http://www.eyou.com/',
    '21cn.com': 'http://mail.21cn.com/',
    '188.com': 'http://www.188.com/',
    'foxmail.coom': 'http://www.foxmail.com'
};

function getEmailLoginAddress(email) {
	if(!email) {
		return;
	}
	var arr = email.split('@');
	if(!arr || arr.length < 2) {
		return;
	}
    var addr = arr[1];
    return email2LoginAddress[addr] || "http://mail." + addr;
}

// 返回是否是re.Ok == true
function reIsOk(re) {
	return re && typeof re == "object" && re.Ok;
}

// marker
// 下拉扩展工具栏用, 点击文档导航用, 切换编辑模式时用
LEA.bookmark = null;
LEA.hasBookmark = false;
function saveBookmark() {
	try {
		LEA.bookmark = tinymce.activeEditor.selection.getBookmark(); // 光标, 为了处理后重新定位到那个位置
		// 如果之前没有focus, 则会在文档开头设置bookmark, 添加一行, 不行.
		// $p不是<p>, 很诡异
		// 6-5
		if(LEA.bookmark && LEA.bookmark.id) {
			var $ic = $($("#editorContent_ifr").contents());
			var $body = $ic.find("body");
			var $p = $body.children().eq(0);
			// 找到
			if($p.is("span")) {
				var $children = $p;
				var $c = $children.eq(0);
				if($c.attr("id") == LEA.bookmark.id + "_start") {
					LEA.hasBookmark = false;
					$c.remove();
				} else {
					LEA.hasBookmark = true;
				}
			} else if($p.is("p")) {
				var $children = $p.children();
				if($children.length == 1 && $.trim($p.text()) == "") {
					var $c = $children.eq(0);
					if($c.attr("id") == LEA.bookmark.id + "_start") {
						LEA.hasBookmark = false;
						$p.remove();
					} else {
						LEA.hasBookmark = true;
					}
				} else {
					LEA.hasBookmark = true;
				}
			}
		}
		
	} catch(e) {
	}
}
function restoreBookmark() {
	try {
		if(LEA.hasBookmark) {
			// 必须要focus()!!!
			var editor = tinymce.activeEditor;
			editor.focus();
			editor.selection.moveToBookmark(LEA.bookmark);
		}
	} catch(e) {
	}
}

// 是否是手机浏览器
var u = navigator.userAgent;
LEA.isMobile = /Mobile|Android|iPhone/i.test(u);
// LEA.isMobile = u.indexOf('Android')>-1 || u.indexOf('Linux')>-1;
// LEA.isMobile = false;
//if($("body").width() < 600) {
//	location.href = "/mobile/index";
//}

// 国际化 i18n
function getMsg(key) {
	return MSG[key] || key;
}