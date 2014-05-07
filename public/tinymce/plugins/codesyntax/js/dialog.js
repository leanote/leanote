var editor = top.tinymce.EditorManager.activeEditor; 

// 是否是选择的文本
var text = "";
var isBrush = null;
function isBrushF() {
	var node = editor.selection.getNode();
	var selectedContent = editor.selection.getContent(); // 包含了html标签
	if(node.nodeName != null && node.nodeName == 'PRE' /*&& node.className.indexOf('brush:') != -1*/) {
		isBrush = $(node);
	} else if(selectedContent) {
		try {
			text = $(selectedContent).text();
		} catch(e) {
		}
		// 可能不是一个完整的html, 可能是一个文本此时.html()无
		if(!text) {
			text = selectedContent;
		}
	}
	return false;
}
isBrushF();

$(function() {
	if(isBrush) {
		var brush = "";
		try {
			var brush = $.trim(isBrush.attr("class").split(":")[1]);
		} catch(e) {
		}
		
		$("#lang").val(brush);
		$("#code").val(isBrush.html());
		$("#insertCodeBtn").html("替换");
		
	} else if(text) {
		$("#insertCodeBtn").html("替换");
		$("#code").val(text);
		
		var lang = $.cookie("syntaxLang");
		if(lang) {
			$("#lang").val(lang);
		}
	} else {
		var lang = $.cookie("syntaxLang");
		if(lang) {
			$("#lang").val(lang);
		}
	}
	
	// 代码格式化
	$("#format").click(function() {
		var code = $("#code").val();
		$("#oldCode").html(code);
		$("#code").val(js_beautify(code, 4, " "));
	});
	
	// 撤销
	$("#unFormat").click(function() {
		var oldCode = $("#oldCode").html();
		if(oldCode) {
			$("#code").val(oldCode);
		}
	});
	
	top.hiddenIframeBorder();
});

var syntax = {
	insert: function() {
		var code = $("#code").val();
		if(!code) {
			syntax.close();
			return;
		}
		var lang = $("#lang").val();
		var classes = "";
		if(lang) {
			classes = "brush: " + lang;
		}
		var html = '<pre class="' + classes + '">' + code + "</pre>";
		// 覆盖
		if(isBrush && !text) {
			isBrush.attr("class", classes);
			isBrush.html(code);
		} else {
			editor.insertContent(html);
		}
		
		// lang 保存到cookie中
		$.cookie("syntaxLang", lang);
	
		syntax.close();
	},
	close: function() {
		editor.windowManager.close();
		editor.windowManager.close();
	},
};

