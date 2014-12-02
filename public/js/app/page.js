// 主页渲染
//-------------

//----------------------
// 编辑器模式
function editorMode() {
	this.writingHash = "#writing";
	this.normalHash = "#normal";
	this.isWritingMode = location.hash == this.writingHash;
	
	this.toggleA = null;
}

editorMode.prototype.toggleAText = function(isWriting) {
	var self = this;
	setTimeout(function() {
		var toggleA = $(".toggle-editor-mode a");
		var toggleSpan = $(".toggle-editor-mode span");
		if(isWriting) {
			toggleA.attr("href", self.normalHash);
			toggleSpan.text(getMsg("normalMode"));
		} else {
			toggleA.attr("href", self.writingHash);
			toggleSpan.text(getMsg("writingMode"));
		}	
	}, 0);
}
editorMode.prototype.isWriting = function(hash) {
	return hash == this.writingHash;
}
editorMode.prototype.init = function() {
	this.changeMode(this.isWritingMode);
	var self = this;
	$(".toggle-editor-mode").click(function() {
		// 
		saveBookmark();
		var $a = $(this).find("a");
		var isWriting = self.isWriting($a.attr("href"));
		self.changeMode(isWriting);
		// 
		restoreBookmark();
	});
}
// 改变模式
editorMode.prototype.changeMode = function(isWritingMode) {
	this.toggleAText(isWritingMode);
	if(isWritingMode) {
		this.writtingMode();
	} else {
		this.normalMode();
	}
	
	$("#moreBtn i").removeClass("fa-angle-up").addClass("fa-angle-down");
}

editorMode.prototype.resizeEditor = function() {
	// css还没渲染完
	setTimeout(function() {
		resizeEditor();
	}, 10);
	setTimeout(function() {
		resizeEditor();
	}, 20);
	setTimeout(function() {
		resizeEditor();
	}, 500);
}
editorMode.prototype.normalMode = function() {
	/*
	var w = $(document).width();
	var h = $(document).height();
	$("#lock").css({right:0, bottom:0});
	*/
	
	var $c = $("#editorContent_ifr").contents();
	
	$c.contents().find("#writtingMode").remove();
	$c.contents().find('link[href$="editor-writting-mode.css"]').remove();
			
	$("#noteItemListWrap, #notesAndSort").show();
	$("#noteList").unbind("mouseenter").unbind("mouseleave"); 
	
	var theme = UserInfo.Theme || "default";
	theme += ".css";
	$("#themeLink").attr("href", "/css/theme/" + theme);
	
	$("#mceToolbar").css("height", "30px");
	
//	$("#lock").animate({right:w},1000);
	
	this.resizeEditor();
	
	$("#noteList").width(UserInfo.NoteListWidth);
	$("#note").css("left", UserInfo.NoteListWidth);
}
editorMode.prototype.writtingMode = function() {
	// $("#pageInner").removeClass("animated fadeInUp");
	
	$("#themeLink").attr("href", "/css/theme/writting-overwrite.css");
	
	setTimeout(function() {
		var $c = $("#editorContent_ifr").contents();
		$c.contents().find("head").append('<link type="text/css" rel="stylesheet" href="/css/editor/editor-writting-mode.css" id="writtingMode">');
	}, 0);
		
	$("#noteItemListWrap, #notesAndSort").fadeOut();
	$("#noteList").hover(function() {
		$("#noteItemListWrap, #notesAndSort").fadeIn();
	}, function() {
		$("#noteItemListWrap, #notesAndSort").fadeOut();
	});
	
	// 点击扩展会使html的height生成, 切换后会覆盖css文件的
	$("#mceToolbar").css("height", "40px");
	
	//$("#pageInner").addClass("animated fadeInUp");

	this.resizeEditor();
	
	$("#noteList").width(250);
	$("#note").css("left", 0);
}

editorMode.prototype.getWritingCss = function() {
	if(this.isWritingMode) {
		return ["css/editor/editor-writting-mode.css"];
	}
	return [];
}
var em = new editorMode();

//----------------
// 拖拉改变变宽度
var Resize = {
	lineMove: false,
	mdLineMove: false,
	target: null,
	
	leftNotebook: $("#leftNotebook"),
	notebookSplitter: $("#notebookSplitter"),
	noteList: $("#noteList"),
	noteAndEditor: $("#noteAndEditor"),
	noteSplitter: $("#noteSplitter"),
	note: $("#note"),
	body: $("body"),
	leftColumn: $("#left-column"),
	rightColumn: $("#right-column"),
	mdSplitter: $("#mdSplitter"),
	
	init: function() {
		var self = this;
		self.initEvent();
	},
	
	initEvent: function() {
		var self = this;
		
		// 鼠标点下
		$(".noteSplit").bind("mousedown", function(event) {
			event.preventDefault(); // 防止选择文本
			self.lineMove = true;
			$(this).css("background-color", "#ccc");
			self.target = $(this).attr("id");
			// 防止iframe捕获不了事件
			$("#noteMask").css("z-index", 99999); // .css("background-color", // "#ccc");
		});
		
		// 鼠标点下
		self.mdSplitter.bind("mousedown", function(event) {
			event.preventDefault(); // 防止选择文本
			self.mdLineMove = true;
			$(this).css("background-color", "#ccc");
		});
		
		// 鼠标移动时
		self.body.bind("mousemove", function(event) {
			if(self.lineMove) { // 如果没有这个if会导致不能选择文本
				event.preventDefault();
				self.resize3Columns(event);
			} else if(self.mdLineMove) {
				event.preventDefault();
				self.resizeMdColumns(event);
			}
		});	

		// 鼠标放开, 结束
		self.body.bind("mouseup", function(event) {
			self.stopResize();
			// 取消遮罩
			$("#noteMask").css("z-index", -1);
		});
	},
	// 停止, 保存数据
	stopResize: function() {
		var self = this;
		if(self.lineMove || self.mdLineMove) {
			// ajax保存
			ajaxGet("/user/updateColumnWidth", {mdEditorWidth: UserInfo.MdEditorWidth, notebookWidth: UserInfo.NotebookWidth, noteListWidth: UserInfo.NoteListWidth}, function() {
			});
		}
		self.lineMove = false;
		self.mdLineMove = false;
		$(".noteSplit").css("background", "none");
		self.mdSplitter.css("background", "none");
	},
	
	// 最终调用该方法
	set3ColumnsWidth: function(notebookWidth, noteListWidth) {
		var self = this;
		if(notebookWidth < 150 || noteListWidth < 100) {
			return;
		}
		var noteWidth = self.body.width() - notebookWidth - noteListWidth;
		if(noteWidth < 400) {
			return;
		}
		
		self.leftNotebook.width(notebookWidth);
		self.notebookSplitter.css("left", notebookWidth);
		
		self.noteAndEditor.css("left", notebookWidth);
		self.noteList.width(noteListWidth);
		self.noteSplitter.css("left", noteListWidth);
		self.note.css("left", noteListWidth);
		
		UserInfo.NotebookWidth = notebookWidth;
		UserInfo.NoteListWidth = noteListWidth;
	},
	resize3Columns: function(event, isFromeIfr) {
		var self = this;
		if (isFromeIfr) {
			event.clientX += self.body.width() - self.note.width();
		}
		
		var notebookWidth, noteListWidth;
		if(self.lineMove) {
			if (self.target == "notebookSplitter") {
				notebookWidth = event.clientX;
				noteListWidth = self.noteList.width();
				self.set3ColumnsWidth(notebookWidth, noteListWidth);
			} else {
				notebookWidth = self.leftNotebook.width();
				noteListWidth = event.clientX - notebookWidth;
				self.set3ColumnsWidth(notebookWidth, noteListWidth);
			}
	
			resizeEditor();
		}
	},
	
	// mdeditor
	resizeMdColumns: function(event) {
		var self = this;
		if (self.mdLineMove) {
			var mdEditorWidth = event.clientX - self.leftNotebook.width() - self.noteList.width();
			self.setMdColumnWidth(mdEditorWidth);
		}
	},
	// 设置宽度
	setMdColumnWidth: function(mdEditorWidth) { 
		var self = this;
		if(mdEditorWidth > 100) {
			UserInfo.MdEditorWidth = mdEditorWidth;
			self.leftColumn.width(mdEditorWidth);
			self.rightColumn.css("left", mdEditorWidth);
			self.mdSplitter.css("left", mdEditorWidth);
		}
	}
}

//--------------------------
// 手机端访问之
Mobile = {
	// 点击之笔记
	// 切换到编辑器模式
	noteO: $("#note"),
	bodyO: $("body"),
	setMenuO: $("#setMenu"),
	hashChange: function() {
		var self = Mobile;
		var hash = location.hash;
		// noteId
		if(hash.indexOf("noteId") != -1) {
			self.toEditor(false);
			var noteId = hash.substr(8);
			Note.changeNote(noteId, false, false);
		} else {
			// 笔记本和笔记列表
			self.toNormal(false);
		}
	},
	init: function() {
		var self = this;
		self.isMobile();
		$(window).on("hashchange", self.hashChange);
		self.hashChange();
		/*
		$("#noteItemList").on("tap", ".item", function(event) {
			$(this).click();
		});
		$(document).on("swipeleft",function(e){
			e.stopPropagation();
			e.preventDefault();
			self.toEditor();
		});
		$(document).on("swiperight",function(e){
			e.stopPropagation();
			e.preventDefault();
			self.toNormal();
		});
		*/
	},
	isMobile: function() {
		var u = navigator.userAgent;
		LEA.isMobile = false;
		LEA.isMobile = /Mobile|Android|iPhone|iPad/i.test(u);
		LEA.isIpad =  /iPhone|iPad/i.test(u);
		if(!LEA.isMobile && $(document).width() <= 700){ 
			LEA.isMobile = true
		}
		return LEA.isMobile;
	},
	changeNote: function(noteId) {
		var self = this;
		if(!LEA.isMobile) {return true;}
		self.toEditor(true, noteId);
		return false;
	},
	
	toEditor: function(changeHash, noteId) {
		var self = this;
		self.bodyO.addClass("full-editor");
		self.noteO.addClass("editor-show");
		if(changeHash) {
			if(!noteId) {
				noteId = Note.curNoteId;
			}
			location.hash = "noteId=" + noteId;
		}
	},
	toNormal: function(changeHash) {
		var self = this;
		self.bodyO.removeClass("full-editor");
		self.noteO.removeClass("editor-show");
	
		if(changeHash) {
			location.hash = "notebookAndNote";
		}
	},
	switchPage: function() {
		var self = this;
		if(!LEA.isMobile || LEA.isIpad) {return true;}
		if(self.bodyO.hasClass("full-editor")) {
			self.toNormal(true);
		} else {
			self.toEditor(true);
		}
		return false;
	}
} 


function initSlimScroll() {
	if(Mobile.isMobile()) {
		return;
	}
	$("#notebook").slimScroll({
	    height: "100%", // $("#leftNotebook").height()+"px"
	});
	$("#noteItemList").slimScroll({
	    height: "100%", // ($("#leftNotebook").height()-42)+"px"
	});
	$("#wmd-input").slimScroll({
	    height: "100%", // $("#wmd-input").height()+"px"
	});
	$("#wmd-input").css("width", "100%");
	
	$("#wmd-panel-preview").slimScroll({
	    height: "100%", // $("#wmd-panel-preview").height()+"px"
	});
	
	$("#wmd-panel-preview").css("width", "100%");
}

//-----------
// 初始化编辑器
function initEditor() {
	// editor
	// toolbar 下拉扩展, 也要resizeEditor
	var mceToobarEverHeight = 0;
	$("#moreBtn").click(function() {
		saveBookmark();
		
		var height = $("#mceToolbar").height();

		// 现在是折叠的
		if (height < $("#popularToolbar").height()) {
			$("#mceToolbar").height($("#popularToolbar").height());
			$(this).find("i").removeClass("fa-angle-down").addClass("fa-angle-up");
			mceToobarEverHeight = height;
		} else {
			$("#mceToolbar").height(mceToobarEverHeight);
			$(this).find("i").removeClass("fa-angle-up").addClass("fa-angle-down");
		}
		
		resizeEditor();
		
		restoreBookmark();
	});

	// 初始化编辑器
	tinymce.init({
		setup: function(ed) {
			ed.on('keydown', Note.saveNote);
			// indent outdent
			ed.on('keydown', function(e) {
				var num = e.which ? e.which : e.keyCode;
		    	if (num == 9) { // tab pressed
				
		    		if(!e.shiftKey) {
//		                ed.execCommand('Indent');
		    			// TODO 如果当前在li, ul, ol下不执行!!
		    			// 如果在pre下就加tab
			    		var node = ed.selection.getNode();
						if(node.nodeName == "PRE") {
		                    ed.execCommand('mceInsertRawHTML', false, '\x09'); // inserts tab
						} else {
		                    ed.execCommand('mceInsertRawHTML', false, "&nbsp;&nbsp;&nbsp;&nbsp;"); // inserts 空格
						}
		    		} else {
		    			// delete 4 个空格
//		                ed.execCommand('Outdent');
		    		}
		    		
		            e.preventDefault();
		            e.stopPropagation();   			
		            return false;
		       }
			});
			
			// 为了把下拉菜单关闭
	        ed.on("click", function(e) {
	          $("body").trigger("click");
	        });
	        
	        // 鼠标移上时
	        ed.on("click", function() {
	        	log(ed.selection.getNode())
	        });
		},
		
		// fix TinyMCE Removes site base url
		// http://stackoverflow.com/questions/3360084/tinymce-removes-site-base-urls
		convert_urls:true,
		relative_urls:false,
		remove_script_host:false,
		
		selector : "#editorContent",
		// height: 100,//这个应该是文档的高度, 而其上层的高度是$("#content").height(),
		// parentHeight: $("#content").height(),
		content_css : ["css/bootstrap.css", "css/editor/editor.css"].concat(em.getWritingCss()),
		skin : "custom",
		language: LEA.locale, // 语言
		plugins : [
				"autolink link leaui_image lists charmap hr", "paste",
				"searchreplace leanote_nav leanote_code tabfocus",
				"table directionality textcolor codemirror" ], // nonbreaking
				
		toolbar1 : "formatselect | forecolor backcolor | bold italic underline strikethrough | leaui_image | leanote_code | bullist numlist | alignleft aligncenter alignright alignjustify",
		toolbar2 : "outdent indent blockquote | link unlink | table | hr removeformat | subscript superscript |searchreplace | code | pastetext pasteCopyImage | fontselect fontsizeselect",

		// 使用tab键: http://www.tinymce.com/wiki.php/Plugin3x:nonbreaking
		// http://stackoverflow.com/questions/13543220/tiny-mce-how-to-allow-people-to-indent
		// nonbreaking_force_tab : true,
		
		menubar : false,
		toolbar_items_size : 'small',
		statusbar : false,
		url_converter: false,
		font_formats : "Arial=arial,helvetica,sans-serif;"
				+ "Arial Black=arial black,avant garde;"
				+ "Times New Roman=times new roman,times;"
				+ "Courier New=courier new,courier;"
				+ "Tahoma=tahoma,arial,helvetica,sans-serif;"
				+ "Verdana=verdana,geneva;" + "宋体=SimSun;"
				+ "新宋体=NSimSun;" + "黑体=SimHei;"
				+ "微软雅黑=Microsoft YaHei",
		block_formats : "Header 1=h1;Header 2=h2;Header 3=h3; Header 4=h4;Pre=pre;Paragraph=p",
		codemirror: {
		    indentOnInit: true, // Whether or not to indent code on init. 
		    path: 'CodeMirror', // Path to CodeMirror distribution
		    config: {           // CodeMirror config object
		       //mode: 'application/x-httpd-php',
		       lineNumbers: true
		    },
		    jsFiles: [          // Additional JS files to load
		       // 'mode/clike/clike.js',
		       //'mode/php/php.js'
		    ]
		  },
		  // This option specifies whether data:url images (inline images) should be removed or not from the pasted contents. 
		  // Setting this to "true" will allow the pasted images, and setting this to "false" will disallow pasted images.  
		  // For example, Firefox enables you to paste images directly into any contentEditable field. This is normally not something people want, so this option is "false" by default.
		  paste_data_images: true
	});
	
	// 刷新时保存 参考autosave插件
	window.onbeforeunload = function(e) {
    	Note.curChangedSaveIt();
	}
	
	// 全局ctrl + s
	$("body").on('keydown', Note.saveNote);
}

//-----------------------
// 导航
var random = 1;
function scrollTo(self, tagName, text) {
	var iframe = $("#editorContent_ifr").contents();
	var target = iframe.find(tagName + ":contains(" + text + ")");
	random++;
	
	// 找到是第几个
	// 在nav是第几个
	var navs = $('#leanoteNavContent [data-a="' + tagName + '-' + encodeURI(text) + '"]');
//	alert('#leanoteNavContent [data-a="' + tagName + '-' + encodeURI(text) + '"]')
	var len = navs.size();
	for(var i = 0; i < len; ++i) {
		if(navs[i] == self) {
			break;
		}
	}
	
	if (target.size() >= i+1) {
		target = target.eq(i);
		// 之前插入, 防止多行定位不准
		var top = target.offset().top;
		var nowTop = iframe.scrollTop();
		
		// iframe.scrollTop(top);
		// $(iframe).animate({scrollTop: top}, 300); // 有问题
		
		var d = 200; // 时间间隔
		for(var i = 0; i < d; i++) {
			setTimeout(
			(function(top) {
				return function() {
					iframe.scrollTop(top);
				}
			})(nowTop + 1.0*i*(top-nowTop)/d), i);
		}
		// 最后必然执行
		setTimeout(function() {
			iframe.scrollTop(top);
		}, d+5);
		return;
	}
}

//--------------
// 调用之
$(function() {
	// 窗口缩放时
	$(window).resize(function() {
		Mobile.isMobile();
		resizeEditor();
	});
	
	// 初始化编辑器
	initEditor();

	// 左侧, folder 展开与关闭
	$(".folderHeader").click(function() {
		var body = $(this).next();
		var p = $(this).parent();
		if (!body.is(":hidden")) {
			$(".folderNote").removeClass("opened").addClass("closed");
//					body.hide();
			p.removeClass("opened").addClass("closed");
			$(this).find(".fa-angle-down").removeClass("fa-angle-down").addClass("fa-angle-right");
		} else {
			$(".folderNote").removeClass("opened").addClass("closed");
//					body.show();
			p.removeClass("closed").addClass("opened");
			$(this).find(".fa-angle-right").removeClass("fa-angle-right").addClass("fa-angle-down");
		}
	});
	
	// 导航隐藏与显示
	$("#leanoteNav h1").on("click", function(e) {
		if (!$("#leanoteNav").hasClass("unfolder")) {
			$("#leanoteNav").addClass("unfolder");
		} else {
			$("#leanoteNav").removeClass("unfolder");
		}
	});
	
	// 打开设置
	function openSetInfoDialog(whichTab) {
		showDialogRemote("/user/account", {tab: whichTab});
	}
	// 帐号设置
	$("#setInfo").click(function() {
		openSetInfoDialog(0);
	});
	// 邮箱验证
	$("#wrongEmail").click(function() {
		openSetInfoDialog(1);
	});
	
	$("#setAvatarMenu").click(function() {
		showDialog2("#avatarDialog", {title: "头像设置", postShow: function() {
		}});
	});
	$("#setTheme").click(function() {
		showDialog2("#setThemeDialog", {title: "主题设置", postShow: function() {
			if (!UserInfo.Theme) {
				UserInfo.Theme = "default";
			}
			$("#themeForm input[value='" + UserInfo.Theme + "']").attr("checked", true);
		}});
	});
	
	//---------
	// 主题
	$("#themeForm").on("click", "input", function(e) {
		var val = $(this).val();
		$("#themeLink").attr("href", "/css/theme/" + val + ".css");
		
		ajaxPost("/user/updateTheme", {theme: val}, function(re) {
			if(reIsOk(re)) {
				UserInfo.Theme = val
			}
		});
	});
	
	//-------------
	// 邮箱验证
	if(!UserInfo.Verified) {
//		$("#leanoteMsg").hide();
//		$("#verifyMsg").show();
	}
	
	// 禁止双击选中文字
	$("#notebook, #newMyNote, #myProfile, #topNav, #notesAndSort", "#leanoteNavTrigger").bind("selectstart", function(e) {
		e.preventDefault();
		return false;
	});
	
	// 左侧隐藏或展示
	function updateLeftIsMin(is) {
		ajaxGet("/user/updateLeftIsMin", {leftIsMin: is})
	}
	function minLeft(save) {
		$("#leftNotebook").width(30);
		$("#notebook").hide();
		// 左侧
		$("#noteAndEditor").css("left", 30)	
		$("#notebookSplitter").hide();
		
//		$("#leftSwitcher").removeClass("fa-angle-left").addClass("fa-angle-right");
		
		// logo
		$("#logo").hide();
		$("#leftSwitcher").hide();
		$("#leftSwitcher2").show();
		$("#leftNotebook .slimScrollDiv").hide();
		
		if(save) {
			updateLeftIsMin(true);
		}
	}
	
	function maxLeft(save) {
		$("#noteAndEditor").css("left", UserInfo.NotebookWidth);
		$("#leftNotebook").width(UserInfo.NotebookWidth);
		$("#notebook").show();
		$("#notebookSplitter").show();
		
//		$("#leftSwitcher").removeClass("fa-angle-right").addClass("fa-angle-left");
		
		$("#leftSwitcher2").hide();
		$("#logo").show();
		$("#leftSwitcher").show();
		$("#leftNotebook .slimScrollDiv").show();
		
		if(save) {
			updateLeftIsMin(false);
		}
	}
	
	$("#leftSwitcher2").click(function() {
		maxLeft(true);
	});
	$("#leftSwitcher").click(function() {
		if(Mobile.switchPage()) {
			minLeft(true);
		}
	});
	
	// 得到最大dropdown高度
	// 废弃
	function getMaxDropdownHeight(obj) {
		var offset = $(obj).offset();
		var maxHeight = $(document).height()-offset.top;
		maxHeight -= 70;
		if(maxHeight < 0) {
			maxHeight = 0;
		}	
		
		var preHeight = $(obj).find("ul").height();
		return preHeight < maxHeight ? preHeight : maxHeight;
	}
	
	// mini版
	// 点击展开
	$("#notebookMin div.minContainer").click(function() {
		var target = $(this).attr("target");
		maxLeft(true);
		if(target == "#notebookList") {
			if($("#myNotebooks").hasClass("closed")) {
				$("#myNotebooks .folderHeader").trigger("click");
			}
		} else if(target == "#tagNav") {
			if($("#myTag").hasClass("closed")) {
				$("#myTag .folderHeader").trigger("click");
			}
		} else {
			if($("#myShareNotebooks").hasClass("closed")) {
				$("#myShareNotebooks .folderHeader").trigger("click");
			}
		}
	});
	
	//------------------------
	// 界面设置, 左侧是否是隐藏的
	UserInfo.NotebookWidth = UserInfo.NotebookWidth || $("#notebook").width();
	UserInfo.NoteListWidth = UserInfo.NoteListWidth || $("#noteList").width();
	
	Resize.init();
	Resize.set3ColumnsWidth(UserInfo.NotebookWidth, UserInfo.NoteListWidth);
	Resize.setMdColumnWidth(UserInfo.MdEditorWidth);
	
	if (UserInfo.LeftIsMin) {
		minLeft(false);
	}
	
	// end
	// 开始时显示loading......
	// 隐藏mask
	$("#mainMask").html("");
	$("#mainMask").hide(100);
	
	// 4/25 防止dropdown太高
	// dropdown
	$('.dropdown').on('shown.bs.dropdown', function () {
		var $ul = $(this).find("ul");
		// $ul.css("max-height", getMaxDropdownHeight(this));
	});
	
	//--------
	// 编辑器帮助
	$("#tipsBtn").click(function() {
		showDialog2("#tipsDialog");
	});
	
	//--------
	// 建议
	$("#yourSuggestions").click(function() {
		showDialog2("#suggestionsDialog");
	});
	$("#suggestionBtn").click(function(e) {
		e.preventDefault();
		var suggestion = $.trim($("#suggestionTextarea").val());
		if(!suggestion) {
			$("#suggestionMsg").html("请输入您的建议, 谢谢!").show().addClass("alert-warning").removeClass("alert-success");
			$("#suggestionTextarea").focus();
			return;
		}
		$("#suggestionBtn").html("正在处理...").addClass("disabled");
		$("#suggestionMsg").html("正在处理...");
		$.post("/suggestion", {suggestion: suggestion}, function(ret) {
			$("#suggestionBtn").html("提交").removeClass("disabled");
			if(ret.Ok) {
				$("#suggestionMsg").html("谢谢反馈, 我们会第一时间处理, 祝您愉快!").addClass("alert-success").removeClass("alert-warning").show();
			} else {
				$("#suggestionMsg").html("出错了").show().addClass("alert-warning").removeClass("alert-success");
			}
		});
	});
	
	// 编辑器模式
	em.init();
	
	// 手机端?
	Mobile.init();
});

