// leanote 通用方法

//--------------
// 命名空间
//--------------

// 最上级变量
if(typeof LEA === 'undefined') {
	var LEA = {};
}
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
var Mobile = {}; // 手机端处理
var LeaAce = {};

// markdown
var Converter;
var MarkdownEditor;
var ScrollLink;
var MD;

//-------------
// 全局事件机制

$.extend(LEA, {
	_eventCallbacks: {},
	_listen: function(type, callback) {
        var callbacks = this._eventCallbacks[type] || (this._eventCallbacks[type] = []);
        callbacks.push(callback);
    },
    // on('a b', function(params) {})
    on: function(name, callback) {
        var names = name.split(/\s+/);
        for (var i = 0; i < names.length; ++i) {
        	this._listen(names[i], callback);
        }
        return this;
    },
    // off('a b', function(params) {})
    off: function(name, callback) {
        var types = name.split(/\s+/);
        var i, j, callbacks, removeIndex;
        for (i = 0; i < types.length; i++) {
            callbacks = this._eventCallbacks[types[i].toLowerCase()];
            if (callbacks) {
                removeIndex = null;
                for (j = 0; j < callbacks.length; j++) {
                    if (callbacks[j] == callback) {
                        removeIndex = j;
                    }
                }
                if (removeIndex !== null) {
                    callbacks.splice(removeIndex, 1);
                }
            }
        }
    },
    // LEA.trigger('a', {});
    trigger: function(type, params) {
        var callbacks = this._eventCallbacks[type] || [];
        if (callbacks.length === 0) {
            return;
        }
        for (var i = 0; i < callbacks.length; i++) {
            callbacks[i].call(this, params);
        }
    }
});

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
var tt = t; // 当slimscroll滑动时t被重新赋值了

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
				alert(getMsg("Please sign in firstly!"));
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
	// log("-------------------ajax:");
	// log(url);
	// log(param);
	if(typeof async == "undefined") {
		async = true;
	} else {
		async = false;
	}
	return $.ajax({
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
	return _ajax("GET", url, param, successFunc, failureFunc, async);
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
	// log("-------------------ajaxPostJson:");
	// log(url);
	// log(param);
	
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

function getVendorPrefix() {
  // 使用body是为了避免在还需要传入元素
  var body = document.body || document.documentElement,
    style = body.style,
    vendor = ['webkit', 'khtml', 'moz', 'ms', 'o'],
    i = 0;
 
  while (i < vendor.length) {
    // 此处进行判断是否有对应的内核前缀
    if (typeof style[vendor[i] + 'Transition'] === 'string') {
      return vendor[i];
    }
    i++;
  }
}

//-----------------

//切换编辑器
LEA.isM = false;
LEA.isMarkdownEditor = function() {
	return LEA.isM;
}
function switchEditor(isMarkdown) {
	LEA.isM = isMarkdown;
	// 富文本永远是2
	if(!isMarkdown) {
		$("#editor").show();
		$("#mdEditor").css("z-index", 1).hide();
		
		// 刚开始没有
		$("#leanoteNav").show();
	} else {
		$("#mdEditor").css("z-index", 3).show();
		
		$("#leanoteNav").hide();
	}
}

// editor 设置内容
// 可能是tinymce还没有渲染成功
var previewToken = "<div style='display: none'>FORTOKEN</div>"
var clearIntervalForSetContent;
function setEditorContent(content, isMarkdown, preview, callback) {
	if(!content) {
		content = "";
	}
	if(clearIntervalForSetContent) {
		clearInterval(clearIntervalForSetContent);
	}
	if(!isMarkdown) {
		// 先destroy之前的ace
		/*
		if(typeof tinymce != "undefined" && tinymce.activeEditor) {
			var editor = tinymce.activeEditor;
			var everContent = $(editor.getBody());
			if(everContent) {
				LeaAce.destroyAceFromContent(everContent);
			}
		}
		*/
		// $("#editorContent").html(content);
		// 不能先setHtml, 因为在tinymce的setContent前要获取之前的content, destory ACE
		if(typeof tinymce != "undefined" && tinymce.activeEditor) {
			var editor = tinymce.activeEditor;
			editor.setContent(content);
			callback && callback();
			editor.undoManager.clear(); // 4-7修复BUG
		} else {
			// 等下再设置
			clearIntervalForSetContent = setTimeout(function() {
				setEditorContent(content, false, false, callback);
			}, 100);
		}
	} else {
	/*
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
				clearIntervalForSetContent = setTimeout(function() {
					setEditorContent(content, true, preview);
				}, 200);
			}
		}
	*/
		if(MD) {
			MD.setContent(content);
			MD.clearUndo && MD.clearUndo();
			callback && callback();
		} else {
			clearIntervalForSetContent = setTimeout(function() {
				setEditorContent(content, true, false, callback);
			}, 100);
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

// ace的值有误?
function isAceError(val) {
	if(!val) {
		return false;
	}
	return val.indexOf('XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX') != -1;
}

// 有tinymce得到的content有<html>包围
// 总会出现<p>&nbsp;<br></p>, 原因, setContent('<p><br data-mce-bogus="1" /></p>') 会设置成 <p> <br></p>
// 所以, 要在getContent时, 当是<p><br data-mce-bogus="1"></p>, 返回 <p><br/></p>
function getEditorContent(isMarkdown) {
	var content = _getEditorContent(isMarkdown);
	if (content === '<p><br data-mce-bogus="1"></p>') {
		return '<p><br></p>';
	}
	return content;
}
function _getEditorContent(isMarkdown) {
	if(!isMarkdown) {
		var editor = tinymce.activeEditor;
		if(editor) {
			var content = $(editor.getBody()).clone();
			// 删除toggle raw 
			content.find('.toggle-raw').remove();

			// single页面没有LeaAce
			if(window.LeaAce && LeaAce.getAce) {
				// 替换掉ace editor
				var pres = content.find('pre');
				for(var i = 0 ; i < pres.length; ++i) {
					var pre = pres.eq(i);
					var id = pre.attr('id');
					var aceEditor = LeaAce.getAce(id);
					if(aceEditor) {
						var val = aceEditor.getValue();
						// 表示有错
						if(isAceError(val)) {
							val = pre.html();
						}
						val = val.replace(/</g, '&lt').replace(/>/g, '&gt');
						pre.removeAttr('style', '').removeAttr('contenteditable').removeClass('ace_editor');
						pre.html(val);
					}
				}
			}
			
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
		// return [$("#wmd-input").val(), $("#wmd-preview").html()]
		return [MD.getContent(), '<div>' + $("#preview-contents").html() + '</div>']
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

function hideDialogRemote(timeout) {
	if(timeout) {
		setTimeout(function() {
			$("#leanoteDialogRemote").modal('hide');
		}, timeout);
	} else {
		$("#leanoteDialogRemote").modal('hide');
	}
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
	if (typeof goNow == 'object') {
		try {
			return goNow.format("yyyy-M-d hh:mm:ss");
		} catch(e) {
			return getCurDate();
		}
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
	LEA.isM && MD && MD.resize && MD.resize();
	return;
	var ifrParent = $("#editorContent_ifr").parent();
    ifrParent.css("overflow", "auto");
    var height = $("#editorContent").height();
    ifrParent.height(height);
    // log(height + '---------------------------------------')
    $("#editorContent_ifr").height(height);

    // life 12.9
    // inline editor
    // $("#editorContent").css("top", $("#mceToolbar").height());
    
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
		$(btnId).button("loading"); // html("正在处理").addClass("disabled");
	}
	ajaxPost(url, param, function(ret) {
		if(btnId) {
			$(btnId).button("reset");
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
	var myreg = /^([a-zA-Z0-9]+[_|\_|\.|\-]?)*[a-zA-Z0-9\-]+@([a-zA-Z0-9\-]+[_|\_|\.|\-]?)*[a-zA-Z0-9\-]+\.[0-9a-zA-Z]{2,3}$/;
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
		msg(msgId, selfBlankMsg || getMsg("inputEmail"));
	} else if(!isEmail(val)) {
		msg(msgId, selfInvalidMsg || getMsg("errorEmail"));
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
function setCookie(c_name, value, expiredays) {
	var exdate = new Date();
	exdate.setDate(exdate.getDate() + expiredays);
	document.cookie = c_name+ "=" + escape(value) + ((expiredays==null) ? "" : ";expires="+exdate.toGMTString()) + 'path=/';
	document.cookie = c_name+ "=" + escape(value) + ((expiredays==null) ? "" : ";expires="+exdate.toGMTString()) + 'path=/note';
}
function logout() {
	Note.curChangedSaveIt(true);
	LEA.isLogout = true;

	setCookie("LEANOTE_SESSION", '', -1);
	location.href = UrlPrefix + "/logout?id=1";
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
    'foxmail.com': 'http://mail.foxmail.com'
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
// var u = navigator.userAgent;
// LEA.isMobile = /Mobile|Android|iPhone/i.test(u);
// LEA.isMobile = u.indexOf('Android')>-1 || u.indexOf('Linux')>-1;
// LEA.isMobile = false;
//if($("body").width() < 600) {
//	location.href = "/mobile/index";
//}

// 表单验证
var vd = {
	isInt: function(o) {
	    var intPattern=/^0$|^[1-9]\d*$/; //整数的正则表达式
	    result=intPattern.test(o);
	    return result;
	},
	isNumeric: function(o) {
		return $.isNumeric(o);
	},
	isFloat: function(floatValue){
	    var floatPattern=/^0(\.\d+)?$|^[1-9]\d*(\.\d+)?$/; //小数的正则表达式
	    result=floatPattern.test(floatValue);
	    return result;
	},
	isEmail: function(emailValue){
	    var emailPattern=/^([a-zA-Z0-9]+[_|\_|\.|\-]?)*[a-zA-Z0-9\-]+@([a-zA-Z0-9\-]+[_|\_|\.|\-]?)*[a-zA-Z0-9\-]+\.[0-9a-zA-Z]{2,3}$/; //邮箱的正则表达式
	    result=emailPattern.test(emailValue);
	   
	    return result;
	},
	isBlank: function(o) { 
		return !$.trim(o);
	},
	has_special_chars: function(o) {
		return /['"#$%&\^<>\?*]/.test(o);
	},
	
	// life
	// 动态验证
	// rules = {max: function() {}};
	// <input data-rules='[{rule: 'requried', msg:"请填写标题"}]' data-msg_target="#msg"/>
	init: function(form, rule_funcs) {
		var get_val = function(target) {
			if(target.is(":checkbox")) {
				var name = target.attr('name');
				var val = $('input[name="' + name + '"]:checked').length;
				return val;
			} else if(target.is(":radio")) {
			} else {
				return target.val();
			}
		}
		var default_rule_funcs = {
			// 必须输入
			required: function(target) {
				return get_val(target);
			},
			// 最少
			min: function(target, rule) {
				var val = get_val(target);
				if(val === "" && !is_required(target)) {
					return true;
				}
				if(val < rule.data) {
					return false;
				}
				return true;
			},
			minLength: function(target, rule) {
				var val = get_val(target);
				if(val === "" && !is_required(target)) {
					return true;
				}
				if(val.length < rule.data) {
					return false;
				}
				return true;
			},
			email: function(target, rule) {
				var val = get_val(target);
				if(val === "" && !is_required(target)) {
					return true;
				}
				return vd.isEmail(val);
			},
			noSpecialChars: function(target) {
				var val = get_val(target);
				if(!val) {
					return true;
				}
				if(/[^0-9a-zzA-Z_\-]/.test(val)) {
					return false;
				}
				return true;
			},
			password: function(target, rule) {
				var val = get_val(target);
				if(val === "" && !is_required(target)) {
					return true;
				}
				return val.length >= 6
			},
			equalTo: function(target, rule) {
				var val = get_val(target);
				if(val === "" && !is_required(target)) {
					return true;
				}
				return $(rule.data).val() == val;
			}
		}
		rule_funcs = rule_funcs || {};
		rule_funcs = $.extend(default_rule_funcs, rule_funcs);
		var rules = {}; // name对应的
		var msg_targets = {};
		// 是否是必须输入的
		function is_required(target) { 
			var name = get_name(target);
			var rules = get_rules(target, name);
			var required_rule = rules[0];
			if(required_rule['rule'] == "required")  {
				return true;
			}
			return false;
		}
		// 先根据msg_target_name, 再根据name
		function get_rules(target, name) {
			if(!rules[name]) {
				rules[name] = eval("(" + target.data("rules") + ")");
			}
			return rules[name];
		}
		
		// 以name为索引, 如果多个input name一样, 但希望有不同的msg怎么办?
		// 添加data-u_name=""
		function get_msg_target(target, name) {
			if(!msg_targets[name]) {
				var t = target.data("msg_target");
				if(!t) {
					// 在其父下append一个
					var msg_o = $('<div class="help-block alert alert-warning" style="display: block;"></div>');
					target.parent().append(msg_o);
					msg_targets[name] = msg_o;
				} else {
					msg_targets[name] = $(t);
				}
			}
			
			return msg_targets[name];
		}
		function hide_msg(target, name) {
			var msgT = get_msg_target(target, name);
			// 之前是正确信息, 那么不隐藏
			if(!msgT.hasClass("alert-success")) {
				msgT.hide();
			}
		}
		function show_msg(target, name, msg, msgData) {
			var t = get_msg_target(target, name);
			t.html(getMsg(msg, msgData)).removeClass("hide alert-success").addClass("alert-danger").show();
		}
		
		// 验证前修改
		function pre_fix(target) {
			var fix_name = target.data("pre_fix");
			if(!fix_name) {
				return;
			}
			switch(fix_name) {
				case 'int': int_fix(target);
				break;
				case 'price': price_fix(target);
				break;
				case 'decimal': decimal_fix(target);
				break;
			}
		}
		
		// 验证各个rule
		// 正确返回true
		function apply_rules(target, name) {
			var rules = get_rules(target, name);
			
			// 是否有前置fix data-pre_fix
			pre_fix(target);
			
			if(!rules) {
				return true;
			}
			for(var i = 0; i < rules.length; ++i) {
				var rule = rules[i];
				var rule_func_name = rule.rule;
				var msg = rule.msg;
				var msgData = rule.msgData;
				if(!rule_funcs[rule_func_name](target, rule)) {
					show_msg(target, name, msg, msgData);
					return false;
				}
			}
			
			hide_msg(target, name);
			
			// 这里, 如果都正确, 是否有sufix验证其它的
			var post_rule = target.data('post_rule');
			if(post_rule) {
				setTimeout(function() {
					var post_target = $(post_rule);
					apply_rules(post_target, get_name(post_target));
				},0);
			}
			
			return true;
		}
		
		function focus_func(e) {
			var target = $(e.target);
			var name = get_name(target);
			// 验证如果有错误, 先隐藏
			hide_msg(target, name);
			
			// key up的时候pre_fix
			pre_fix(target);
		}
		function unfocus_func(e) {
			var target = $(e.target);
			var name = get_name(target);
			// 验证各个rule
			apply_rules(target, name);
		}
		
		// u_name是唯一名, msg, rule的索引
		function get_name(target) {
			return target.data('u_name') || target.attr("name") || target.attr("id");
		}
		
		var $allElems = $(form).find('[data-rules]');
		var $form = $(form);
		$form.on({
			keyup: function(e) {
				if(e.keyCode != 13) { // 不是enter
					focus_func(e)
				}
			},
			blur: unfocus_func,
		}, 'input[type="text"], input[type="password"]');
		$form.on({
			change: function(e) {
				if($(this).val()) {
					focus_func(e);
				} else {
					unfocus_func(e);
				}
			}
		}, 'select');
		$form.on({
			change: function(e) {
				unfocus_func(e);
			}
		}, 'input[type="checkbox"]');
		
		// 验证所有的
		this.valid = function() {
			var $ts = $allElems;
			var is_valid = true;
			for(var i = 0; i < $ts.length; ++i) {
				var target = $ts.eq(i);
				var name = get_name(target);
				// 验证各个rule
				if(!apply_rules(target, name)) {
					is_valid = false;
					target.focus();
					return false
				} else {
				}
			}
			return is_valid;
		}
		
		// 验证某一元素(s)
		// .num-in, #life
		this.validElement = function(targets) {
			var targets = $(targets);
			var ok = true;
			for(var i = 0; i < targets.length; ++i) {
				var target = targets.eq(i);
				var name = get_name(target);
				// 验证各个rule
				if(!apply_rules(target, name)) {
					ok = false;
				}
			}
			return ok;
		}
	}
};

// 返回hash的#a=1&b=3 返回{a:1, b:3}
function getHashObject() {
	var hash = location.hash; // #life	
	if(!hash) {
		return {};
	}
	var hashKV = hash.substr(1);
	var kvs = hashKV.split("&");
	var kvsObj = {};
	for(var i = 0; i < kvs.length; ++i) {
		var kv = kvs[i].split('=');
		if(kv.length == 2) {
			kvsObj[kv[0]] = kv[1];
		}
	}
	return kvsObj;
}
function getHash(key, value) {
	var kvs = getHashObject();
	return kvs[key];
}
function setHash(key, value) {
	var hash = location.hash; // #life	
	if(!hash) {
		location.href = "#" + key + "=" + value;
		return;
	}
	var kvs = getHashObject();
	kvs[key] = value;
	var str = "";
	for(var i in kvs) {
		if(kvs[i]) {
			if(str) {
				str += "&";
			}
			str += i + '=' + kvs[i];
		}
	}
	location.href = "#" + str;
}

var trimTitle = function(title) {
	if(!title || typeof title != 'string') {
		return '';
	}
	return title.replace(/</g, "&lt;").replace(/>/g, "&gt;");
	// return title.replace(/<.*?script.*?>/g, '');
};
