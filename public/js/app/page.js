// 主页渲染
//-------------

function sendLog (key, value) {
	if (!key) {
		return;
	}
	if (!value) {
		value = '';
	}
	ajaxGet('/index/log', {key: key, value: value});
}

//----------------------
// 编辑器模式
function editorMode() {
	this.writingHash = "writing";
	this.normalHash = "normal";
	this.isWritingMode = location.hash.indexOf(this.writingHash) >= 0;
	this.toggleA = null;
}

editorMode.prototype.toggleAText = function(isWriting) {
	var self = this;
	setTimeout(function() {
		var toggleA = $(".toggle-editor-mode a");
		var toggleSpan = $(".toggle-editor-mode span");
		if(isWriting) {
			toggleA.attr("href", "#" + self.normalHash);
			toggleSpan.text(getMsg("normalMode"));
		} else {
			toggleA.attr("href", "#" + self.writingHash);
			toggleSpan.text(getMsg("writingMode"));
		}	
	}, 0);
}
editorMode.prototype.isWriting = function(hash) {
	if(!hash) {
		hash = location.hash;
	}
	return hash.indexOf(this.writingHash) >= 0
}
editorMode.prototype.init = function() {
	this.$themeLink = $("#themeLink");
	this.changeMode(this.isWritingMode);
	var self = this;
	$(".toggle-editor-mode").click(function(e) {
		e.preventDefault();
		saveBookmark();
		var $a = $(this).find("a");
		var isWriting = self.isWriting($a.attr("href"));
		self.changeMode(isWriting);
		// 
		if(isWriting) {
			setHash("m", self.writingHash);
		} else {
			setHash("m", self.normalHash);
		}
		
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
};

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
	// 最开始的时候就调用?
	/*
	var $c = $("#editorContent_ifr").contents();
	$c.contents().find("#writtingMode").remove();
	$c.contents().find('link[href$="editor-writting-mode.css"]').remove();
	*/

	$("#noteItemListWrap, #notesAndSort").show();
	$("#noteList").unbind("mouseenter").unbind("mouseleave"); 
	
	var theme = UserInfo.Theme || "default";
	theme += ".css";
	var $themeLink = $("#themeLink");
	// 如果之前不是normal才换
	if(this.$themeLink.attr('href').indexOf('writting-overwrite.css') != -1) {
		this.$themeLink.attr("href", LEA.sPath + "/css/theme/" + theme);
	}
	
	$("#noteList").width(UserInfo.NoteListWidth);
	$("#note").css("left", UserInfo.NoteListWidth);

	this.isWritingMode = false;
	this.resizeEditor();
};

editorMode.prototype.writtingMode = function() {
	if (Note.inBatch) {
		return;
	}
	if(this.$themeLink.attr('href').indexOf('writting-overwrite.css') == -1) {
		this.$themeLink.attr("href", LEA.sPath + "/css/theme/writting-overwrite.css");
	}

	/*
	setTimeout(function() {
		var $c = $("#editorContent_ifr").contents();
		$c.contents().find("head").append('<link type="text/css" rel="stylesheet" href="/css/editor/editor-writting-mode.css" id="writtingMode">');
	}, 0);
	*/
		
	$("#noteItemListWrap, #notesAndSort").fadeOut();
	$("#noteList").hover(function() {
		$("#noteItemListWrap, #notesAndSort").fadeIn();
	}, function() {
		$("#noteItemListWrap, #notesAndSort").fadeOut();
	});
	
	// 点击扩展会使html的height生成, 切换后会覆盖css文件的
	// $("#mceToolbar").css("height", "40px");
	
	//$("#pageInner").addClass("animated fadeInUp");

	this.resizeEditor();
	
	$("#noteList").width(250);
	$("#note").css("left", 0);
	
	// 切换到写模式
	Note.toggleWriteable();

	this.isWritingMode = true;
};

editorMode.prototype.getWritingCss = function() {
	if(this.isWritingMode) {
		return ["/css/editor/editor-writting-mode.css"];
	}
	return [];
}
var em = new editorMode();
LEA.em = em;

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
	rightColumn: $("#right-column"), // $("#preview-panel"), // 
	mdSplitter: $("#mdSplitter2"),
	
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
			if($(this).hasClass('open')) {
				self.mdLineMove = true;
			}
			// $(this).css("background-color", "#ccc");
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
		
		// 瞬间
		var everLeftWidth;
		$('.layout-toggler-preview').click(function() {
			var $t = $(this);
			var $p = self.leftColumn.parent();
			// 是开的
			if($t.hasClass('open')) {
				var totalWidth = $p.width();
				var minRightWidth = 22;
				var leftWidth = totalWidth - minRightWidth;
				everLeftWidth = self.leftColumn.width();
				self.leftColumn.width(leftWidth);
				self.rightColumn.css('left', 'auto').width(minRightWidth);
				
				// 禁止split
				$t.removeClass('open');//.addClass('close');
				self.rightColumn.find('.layout-resizer').removeClass('open');
				$('.preview-container').hide();

				if(MD) {
					MD.resize();
				}
			} else {
				$t.addClass('open');
				self.rightColumn.find('.layout-resizer').addClass('open');
				self.leftColumn.width(everLeftWidth);
				$('.preview-container').show();
				self.rightColumn.css('left', everLeftWidth).width('auto');
				
				if(MD) { 
					MD.resize();
				}
			}
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

	resizeMDInterval: null,
	// mdeditor
	resizeMdColumns: function(event) {
		var self = this;
		if (self.mdLineMove) {
			var mdEditorWidth = event.clientX - self.leftColumn.offset().left; // self.leftNotebook.width() - self.noteList.width();
			self.setMdColumnWidth(mdEditorWidth);

			clearInterval(self.resizeMDInterval);

			self.resizeMDInterval = setTimeout(function () {
				MD.resize && MD.resize();
			}, 50);
		}
	},
	// 设置宽度
	setMdColumnWidth: function(mdEditorWidth) { 
		var self = this;
		var allWidth = $('#note').width();
		if(mdEditorWidth > 100 && mdEditorWidth < allWidth - 80) {
			UserInfo.MdEditorWidth = mdEditorWidth;
			self.leftColumn.width(mdEditorWidth);
			self.rightColumn.css("left", mdEditorWidth);
			// self.mdSplitter.css("left", mdEditorWidth);
		}

		// 这样, scrollPreview 才会到正确的位置
		if(MD) {
			MD.onResize();
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
	// 弃用, 统一使用Pjax
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
		// $(window).on("hashchange", self.hashChange);
		// self.hashChange();
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
		LEA.isIpad =  /iPad/i.test(u);
		LEA.isIphone = /iPhone/i.test(u);
		if(!LEA.isMobile && $(document).width() <= 700){ 
			LEA.isMobile = true
		}
		return LEA.isMobile;
	},
	// 改变笔记, 此时切换到编辑器模式下
	// note.js click事件处理, 先切换到纯编辑器下, 再调用Note.changeNote()
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
		/*
		if(changeHash) {
			if(!noteId) {
				noteId = Note.curNoteId;
			}
			location.hash = "noteId=" + noteId;
		}
		*/
	},
	toNormal: function(changeHash) {
		var self = this;
		self.bodyO.removeClass("full-editor");
		self.noteO.removeClass("editor-show");
	
		/*
		if(changeHash) {
			location.hash = "notebookAndNote";
		}
		*/
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
	/*
	$("#wmd-input").slimScroll({
	    height: "100%", // $("#wmd-input").height()+"px"
	});
	$("#wmd-input").css("width", "100%");
	*/
	
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
		var $editor = $('#editor');
		if($editor.hasClass('all-tool')) {
			$editor.removeClass('all-tool');
		} else {
			$editor.addClass('all-tool');
		}

		restoreBookmark();
	});

	// 初始化编辑器
	tinymce.init({
		inline: true,
		theme: 'leanote',
		valid_children: "+pre[div|#text|p|span|textarea|i|b|strong]", // ace
		/*
		protect: [
	        /\<\/?(if|endif)\>/g, // Protect <if> & </endif>
	        /\<xsl\:[^>]+\>/g, // Protect <xsl:...>
	        // /<pre.*?>.*?<\/pre>/g, // Protect <pre ></pre>
	        // /<p.*?>.*?<\/p>/g, // Protect <pre ></pre>
	        // /<\?php.*?\?>/g // Protect php code
	    ],
	    */
		setup: function(ed) {
			ed.on('keydown', function(e) {
				// 如果是readony, 则不能做任何操作
				var num = e.which ? e.which : e.keyCode;
				// 如果是readony, 则不能做任何操作, 除了复制
				if(Note.readOnly && !((e.ctrlKey || e.metaKey) && num == 67)) {
					e.preventDefault();
					return;
				}

				// 当输入的时候, 把当前raw删除掉
				LeaAce.removeCurToggleRaw();
			});
			
			// 为了把下拉菜单关闭
			/*
	        ed.on("click", function(e) {
	          // $("body").trigger("click");
	          // console.log(tinymce.activeEditor.selection.getNode());
	        });
	        */
	        
	        // electron下有问题, Ace剪切导致行数减少, #16
			ed.on('cut', function(e) {
				if($(e.target).hasClass('ace_text-input')) {
					e.preventDefault();
					return;
				}
			});
		},
		
		// fix TinyMCE Removes site base url
		// http://stackoverflow.com/questions/3360084/tinymce-removes-site-base-urls
		convert_urls: false, // true会将url变成../api/
		relative_urls: true,
		remove_script_host:false,
		
		selector : "#editorContent",
		
		// content_css 不再需要
		// content_css : [LEA.sPath + "/css/editor/editor.css"], // .concat(em.getWritingCss()),
		skin : "custom",
		language: LEA.locale, // 语言
		plugins : [
				"autolink link leaui_image leaui_mindmap lists hr", "paste",
				"searchreplace leanote_nav leanote_code tabfocus",
				"table textcolor" ], // nonbreaking directionality charmap
		toolbar1 : "formatselect | forecolor backcolor | bold italic underline strikethrough | leaui_image leaui_mindmap | leanote_code leanote_inline_code | bullist numlist | alignleft aligncenter alignright alignjustify",
		toolbar2 : "outdent indent blockquote | link unlink | table | hr removeformat | subscript superscript |searchreplace | pastetext | leanote_ace_pre | fontselect fontsizeselect",

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
		block_formats : "Header 1=h1;Header 2=h2;Header 3=h3;Header 4=h4;Paragraph=p",
		/*
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
		  */
		  // This option specifies whether data:url images (inline images) should be removed or not from the pasted contents. 
		  // Setting this to "true" will allow the pasted images, and setting this to "false" will disallow pasted images.  
		  // For example, Firefox enables you to paste images directly into any contentEditable field. This is normally not something people want, so this option is "false" by default.
		  paste_data_images: true
	});
	
	// 刷新时保存 参考autosave插件
	window.onbeforeunload = function(e) {
		if (LEA.isLogout) {
			return;
		}
    	Note.curChangedSaveIt(true, null, {refresh: true});
	}

	// 全局快捷键
	// ctrl + s 保存
	// ctrl+e 切换只读与可写
	$('body').on('keydown', function (e) {
		var num = e.which ? e.which : e.keyCode;
		var ctrlOrMetaKey = e.ctrlKey || e.metaKey;
	    if(ctrlOrMetaKey) {
			// 保存
		    if (num == 83 ) { // ctrl + s or command + s
		    	Note.curChangedSaveIt(true, null, {ctrls: true});
		    	e.preventDefault();
		    	return false;
		    }
		    else if (num == 69) { // e
		    	Note.toggleWriteableAndReadOnly();
		    	e.preventDefault();
		    	return false;
		    }
	    }
	});
}

//-----------------------
// 导航
var random = 1;
function scrollTo(self, tagName, text) {
	var iframe = $("#editorContent"); // .contents();
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
		// log(target.scrollTop());
		var top = iframe.scrollTop() - iframe.offset().top + target.offset().top; // 相对于iframe的位置
		// var nowTop = iframe.scrollTop();
		// log(nowTop);
		// log(top);
		// iframe.scrollTop(top);
		iframe.animate({scrollTop: top}, 300); // 有问题
		
		/*
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
		*/
		return;
	}
}

function hideMask () {
	$("#mainMask").html("");
	$("#mainMask").hide(100);
}

//--------------
// 调用之
// $(function() {
	LEA.s3 = new Date();
	console.log('initing...');
	
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
	$(".leanoteNav h1").on("click", function(e) {
		var $leanoteNav = $(this).closest('.leanoteNav');
		if (!$leanoteNav.hasClass("unfolder")) {
			$leanoteNav.addClass("unfolder");
		} else {
			$leanoteNav.removeClass("unfolder");
		}
	});
	
	// 邮箱验证
	$("#wrongEmail").click(function() {
		openSetInfoDialog(1);
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
		var preHref = $("#themeLink").attr("href"); // default.css?id=7
		var arr = preHref.split('=');
		var id = 1;
		if (arr.length == 2) {
			id = arr[1];
		}
		$("#themeLink").attr("href", LEA.sPath + "/css/theme/" + val + ".css?id=" + id);
		ajaxPost("/user/updateTheme", {theme: val}, function(re) {
			if(reIsOk(re)) {
				UserInfo.Theme = val
			}
		});
	});

	// 禁止双击选中文字
	$("#notebook, #newMyNote, #myProfile, #topNav, #notesAndSort", "#leanoteNavTrigger").bind("selectstart", function(e) {
		e.preventDefault();
		return false;
	});
	
	// 左侧隐藏或展示
	function updateLeftIsMin(is) {
		ajaxGet("/user/updateLeftIsMin", {leftIsMin: is})
	}

	// 最小化左侧
	var $page = $('#page');
	function minLeft(save) {
		$page.addClass('mini-left');
		if(save) {
			updateLeftIsMin(true);
		}
	}

	// 展开右侧
	function maxLeft(save) {
		$page.removeClass('mini-left');
		$("#noteAndEditor").css("left", UserInfo.NotebookWidth);
		$("#leftNotebook").width(UserInfo.NotebookWidth);
		if(save) {
			updateLeftIsMin(false);
		}
	}
	
	$("#leftSwitcher2").on('click', function() {
		maxLeft(true);
	});
	$("#leftSwitcher").click('click', function() {
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
	
	if (!Mobile.isMobile()) {
		if (UserInfo.LeftIsMin) {
			minLeft(false);
		}
		else {
			maxLeft(false);
		}
	}
	else {
		maxLeft(false);
	}
	
	// end
	// 开始时显示loading......
	// 隐藏mask
	// hideMask();
	
	// 4/25 防止dropdown太高
	// dropdown
	$('.dropdown').on('shown.bs.dropdown', function () {
		var $ul = $(this).find("ul");
		// $ul.css("max-height", getMaxDropdownHeight(this));
	});
	
	/*
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
	*/
	
	// 编辑器模式
	em.init();
	
	// 手机端?
	Mobile.init();
//});

//------------
// pjax
//------------
var Pjax = {
	init: function() {
		var me = this;
		// 当history改变时
		window.addEventListener('popstate', function(evt){
			var state = evt.state;
			if(!state) {
				return;
			}
			document.title = state.title || "Untitled";
			log("pop");
			me.changeNotebookAndNote(state.noteId);
		}, false);
		
		// ie9
		if(!history.pushState) {
			$(window).on("hashchange", function() {
				var noteId = getHash("noteId");;
				if(noteId) {
					me.changeNotebookAndNote(noteId);
				}
			});
		}
	},
	// pjax调用
	// popstate事件发生时, 转换到noteId下, 此时要转换notebookId
	changeNotebookAndNote: function(noteId) {
		var note = Note.getNote(noteId);
		if(!note) {
			return;
		}
		var isShare = note.Perm != undefined;
		
		var notebookId = note.NotebookId;
		// 如果是在当前notebook下, 就不要转换notebook了
		if(Notebook.curNotebookId == notebookId) {
			// 不push state
			Note.changeNoteForPjax(noteId, false);
			return;
		}
		
		// 自己的
		if(!isShare) {
			// 先切换到notebook下, 得到notes列表, 再changeNote
			Notebook.changeNotebook(notebookId, function(notes) {
				Note.renderNotes(notes);
				// 不push state
				Note.changeNoteForPjax(noteId, false, true);
			});
		// 共享笔记
		} else {
			Share.changeNotebook(note.UserId, notebookId, function(notes) {
				Note.renderNotes(notes);
				// 不push state
				Note.changeNoteForPjax(noteId, false, true);
			});
		}
	},
		
	// ajax后调用
	changeNote: function(noteInfo) {
		var me = this;
		var noteId = noteInfo.NoteId;
		var title = noteInfo.Title;
		var url = '/note/' + noteId;
		if (location.href.indexOf('?online') > 0) {
			url += '?online=' + /online=([0-9])/.exec(location.href)[1];
		}
		if(location.hash) {
			url += location.hash;
		}
		// 如果支持pushState
		if(history.pushState) {
			var state=({
				url: url,
				noteId: noteId,
				title: title,
			});
			history.pushState(state, title, url);
			document.title = title || 'Untitled';
		// 不支持, 则用hash
		} else {
			setHash("noteId", noteId);
		}
	}
};
$(function() {
	Pjax.init();
});

//----------
// aceEditor
LeaAce = {
	// aceEditorID
	_aceId: 0,
	// {id=>ace}
	_aceEditors: {},
	_isInit: false,
	_canAce: false,
	isAce: true, // 切换pre, 默认是true
	disableAddHistory: function() {
		tinymce.activeEditor.undoManager.setCanAdd(false);
	},
	resetAddHistory: function() {
		tinymce.activeEditor.undoManager.setCanAdd(true);
	},
	canAce: function() {
		if(this._isInit) {
			return this._canAce;
		}
		if(getVendorPrefix() == "webkit" && !Mobile.isMobile()) {
			this._canAce = true;
		} else {
			this._canAce = false;
		}
		this._isInit = true;
		return this._canAce;
	},
	canAndIsAce: function() {
		return this.canAce() && this.isAce;
	},
	getAceId: function () {
		this.aceId++;
		return "leanote_ace_" + (new Date()).getTime() + "_" + this._aceId;
	},
	initAce: function(id, val, force) {
		var me = this;
		if(!force && !me.canAndIsAce()) {
			return;
		}
		var $pre = $('#' + id);
		if($pre.length == 0) {
			return;
		}
		var rawCode = $pre.html(); // 原生code
		try {
			me.disableAddHistory();
			
			// 本身就有格式的, 防止之前有格式的显示为<span>(ace下)
			var classes = $pre.attr('class') || '';
			var isHtml = classes.indexOf('brush:html') != -1;
			if($pre.attr('style') || 
				(!isHtml && $pre.html().indexOf('style') != -1)) { // 如果是html就不用考虑了, 因为html格式的支持有style
				$pre.html($pre.text());
			}
			$pre.find('.toggle-raw').remove();
			var preHtml = $pre.html();

			$pre.removeClass('ace-to-pre');
			$pre.attr("contenteditable", false); // ? 避免tinymce编辑
			var aceEditor = ace.edit(id);

			aceEditor.container.style.lineHeight = 1.5;
			aceEditor.setTheme("ace/theme/tomorrow");

			var brush = me.getPreBrush($pre);
			var b = "";
			if(brush) {
				try {
					b = brush.split(':')[1];
				} catch(e) {}
			}
			if (!b || b === 'false') {
				b = 'javascript';
			}
			
			aceEditor.session.setMode("ace/mode/" + b);
			aceEditor.session.setOption("useWorker", false); // 不用语法检查
			// retina
			if(window.devicePixelRatio == 2) {
				aceEditor.setFontSize("12px");
			}
			else {
				aceEditor.setFontSize("14px");
			}
			aceEditor.getSession().setUseWorker(false); // 不用语法检查
			aceEditor.setOption("showInvisibles", false); // 不显示空格, 没用
			aceEditor.setShowInvisibles(false); // OK 不显示空格
			aceEditor.setOption("wrap", "free");
			aceEditor.setShowInvisibles(false);
			
			aceEditor.setReadOnly(Note.readOnly);
			
			aceEditor.setAutoScrollEditorIntoView(true);
			aceEditor.setOption("maxLines", 10000);
			aceEditor.commands.addCommand({
			    name: "undo",
			    bindKey: {win: "Ctrl-z", mac: "Command-z"},
			    exec: function(editor) {
			    	var undoManager = editor.getSession().getUndoManager();
			    	if(undoManager.hasUndo()){ 
			    		undoManager.undo();
			    	} else {
			    		undoManager.reset();
			    		tinymce.activeEditor.undoManager.undo();
			    	}
			    }
			});
			this._aceEditors[id] = aceEditor;
			if(val) {
				aceEditor.setValue(val);
				// 不要选择代码
				// TODO
			} else {
				// 防止 <pre><div>xx</div></pre> 这里的<div>消失
				// preHtml = preHtml.replace('/&nbsp;/g', ' '); // 以前是把' ' 全换成了&nbsp;
				// aceEditor.setValue(preHtml);
				// 全不选
				// aceEditor.selection.clearSelection();
			}

			// XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
			// "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
			me.resetAddHistory();
			return aceEditor;
		} catch(e) {
			// 当有错误时, 会有XXXXX的形式, 此时不要ace, 直接原生的!!!
			console.error('ace error!!!!');
			console.error(e);
			$pre.attr("contenteditable", true);
			$pre.removeClass('ace-tomorrow ace_editor ace-tm');
			$pre.html(rawCode);
			me.resetAddHistory();
		}
	},
	clearIntervalForInitAce: null,
	initAceFromContent: function(editor) {
		if(!this.canAndIsAce()) {
			var content = $(editor.getBody());
			content.find('pre').removeClass('ace_editor');
			return;
		}
		var me = this;
		// 延迟
		if(this.clearIntervalForInitAce) {
			clearInterval(this.clearIntervalForInitAce);
		}
		this.clearIntervalForInitAce = setTimeout(function() {
			var content = $(editor.getBody());
			var pres = content.find('pre');
			for(var i = 0 ; i < pres.length; ++i) {
				var pre = pres.eq(i);
				
				var aceAndNode = me.isInAce(pre);
				if(aceAndNode) {
					if(isAceError(aceAndNode[0].getValue())) {
						console.error('之前有些没有destroy掉');
					}
					else {
						break;
					}
				}
				
				setTimeout((function(pre) {
					return function() {
						pre.find('.toggle-raw').remove();
						var value = pre.html();
						value = value.replace(/ /g, "&nbsp;").replace(/\<br *\/*\>/gi,"\n").replace(/</g, '&lt;').replace(/>/g, '&gt;');
						pre.html(value);
						var id = pre.attr('id');
						if(!id) {
							id = me.getAceId();
							pre.attr('id', id);
						}
						me.initAce(id);
					}
				})(pre));
			}
		}, 10);
	},

	allToPre: function(editor) {
		if(!this.canAndIsAce()) {
			return;
		}
		var me = this;
		// 延迟
		if(me.clearIntervalForInitAce) {
			clearInterval(me.clearIntervalForInitAce);
		}
		me.clearIntervalForInitAce = setTimeout(function() {
			var content = $(editor.getBody());
			var pres = content.find('pre');
			for(var i = 0 ; i < pres.length; ++i) {
				var pre = pres.eq(i);
				setTimeout((function(pre) {
					return function() {
						me.aceToPre(pre);
					}
				})(pre));
			}
		}, 10);
	},

	undo: function(editor) {
		if(!this.canAndIsAce()) {
			return;
		}
		var me = this;
		// 延迟
		if(this.clearIntervalForInitAce) {
			clearInterval(this.clearIntervalForInitAce);
		}
		this.clearIntervalForInitAce = setTimeout(function() {
			var content = $(editor.getBody());
			var pres = content.find('pre');
			for(var i = 0 ; i < pres.length; ++i) {
				var pre = pres.eq(i);
				setTimeout((function(pre) {
					return function() {
						var value = pre.html();
						var id = pre.attr('id');
						var aceEditor = me.getAce(id);
						if(aceEditor) {
							var value = aceEditor.getValue();
							aceEditor.destroy();
							var aceEditor = me.initAce(id, value);
							// 全不选
							aceEditor.selection.clearSelection();
						} else {
							value = value.replace(/ /g, "&nbsp;").replace(/\<br *\/*\>/gi,"\n");
							pre.html(value);
							var id = pre.attr('id');
							if(!id) {
								id = me.getAceId();
								pre.attr('id', id);
							}
							me.initAce(id);
						}
					}
				})(pre));
			}
		}, 10);
	},
	destroyAceFromContent: function(everContent) {
		if(!this.canAce()) {
			return;
		}
		var pres = everContent.find('pre');
		for(var i = 0 ; i < pres.length; ++i) {
			var id = pres.eq(i).attr('id');
			var aceEditorAndPre = this.getAce(id);
			if(aceEditorAndPre) {
				aceEditorAndPre.destroy();
				this._aceEditors[id] = null;
			}
		}
	},
	getAce: function(id) {
		if(!this.canAce()) {
			return;
		}
		return this._aceEditors[id];
	},
	setAceReadOnly: function(pre, readOnly) {
		var me = this;
		if(typeof pre == 'object') {
			var id = pre.attr('id');
		}
		else {
			var id = pre;
		}
		var ace = me.getAce(id);
		if(ace) {
			ace.setReadOnly(readOnly);
		}
	},
	// 当前焦点是否在aceEditor中
	nowIsInAce: function () {
		if(!this.canAce()) {
			return;
		}
		
		var node = tinymce.activeEditor.selection.getNode();
		// log("now...");
		// log(node);
		return this.isInAce(node);

	},
	nowIsInPre: function(){
		var node = tinymce.activeEditor.selection.getNode();
		// log("now...");
		// log(node);
		return this.isInPre(node);
	},
	isInPre: function(node) {
		var $node = $(node);
		var node = $node.get(0);
		if(node.nodeName == "PRE") {
			return true;
		} else {
			// 找到父是pre
			$pre = $node.closest("pre");
			if($pre.length == 0) {
				return false;
			}
			return true;
		}
	},
	// 是否在node内
	isInAce: function(node) {
		if(!this.canAce()) {
			return;
		}
		var $node = $(node);
		var node = $node.get(0);
		if(node.nodeName == "PRE") {
			// $node.data('brush', brush);
			var id = $node.attr('id');
			var aceEditor = this.getAce(id);
			if(aceEditor) {
				return [aceEditor, $node];
			}
			return false;
		} else {
			// 找到父是pre
			$pre = $node.closest("pre");
			if($pre.length == 0) {
				return false;
			}
			return this.isInAce($pre);
		}
		return false;
	},
	getPreBrush: function (node) {
		var $pre = $(node);
		var classes = $pre.attr('class');
		if(!classes) {
			return '';
		}
		var m = classes.match(/brush:[^ ]*/);
		var everBrush = "";
		if(m && m.length > 0) {
			everBrush = m[0];
		}	
		return everBrush;
	},
	// pre转换成ace
	preToAce: function (pre, force) {
		if(!force && !this.canAce()) {
			return;
		}
		var $pre = $(pre);
		var id = this.getAceId();
		$pre.attr('id', id);
		var editor = this.initAce(id, "", true);
		if(editor) {
			editor.focus();
		}
	},
	aceToPre: function(pre, isFocus) {
		var me = this;
		var $pre = $(pre);
		// 转成pre
		var aceEditorAndPre = me.isInAce($pre);
		if(aceEditorAndPre) {
			var aceEditor = aceEditorAndPre[0];
			var $pre = aceEditorAndPre[1];
			var value = aceEditor.getValue();
			// 表示有错
			if(isAceError(value)) {
				value = $pre.html();
			}
			value = value.replace(/</g, '&lt').replace(/>/g, '&gt');
			// var id = getAceId();
			var replacePre = $('<pre class="' + $pre.attr('class') + ' ace-to-pre">' + value + "</pre>");
			$pre.replaceWith(replacePre);
			aceEditor.destroy();
			me._aceEditors[$pre.attr('id')] = null;
			// log($replacePre);
			if(isFocus) {
				setTimeout(function() {
					var tinymceEditor = tinymce.activeEditor;
					var selection = tinymceEditor.selection;
					var rng = selection.getRng();
					// rng.setStart(replacePre.get(0), 1);
					// rng.setEnd(replacePre.get(0), 9);
					rng.selectNode(replacePre.get(0));
					// selection.setRng(rng);
					// replacePre.focus();
					tinymceEditor.focus();
					replacePre.trigger("click");
					replacePre.html(value + " ");
					// log(">>>>>>>>>>>>>>")
				}, 0);
			}
		}
	},
	// 当删除了pre时, 也要删除toggle raw
	removeAllToggleRaw: function () {
		$('#editorContent .toggle-raw').remove();
	},
	removeCurToggleRaw: function() {
		if(this.curToggleRaw) {
			try {
				this.curToggleRaw.remove();
			}
			catch(e){}
		}
	},
	curToggleRaw: null,
	// 转换raw <-> code
	handleEvent: function () {
		if(!this.canAce()) {
			return;
		}
		var me = this;
		$("#editorContent").on('mouseenter', 'pre', function(e) {
			// log('in');
			// log($(this));
			var $t = $(this);
			$raw = $t.find('.toggle-raw');
			if($raw.length == 0) {
				var curToggleRaw = $('<div class="toggle-raw" title="Toggle code with raw html"><input type="checkbox" /></div>');
				$t.append(curToggleRaw);
				me.curToggleRaw = curToggleRaw;
			}
			$input = $t.find('.toggle-raw input');
			if(LeaAce.isInAce($t)) {
				$input.prop('checked', true);
			} else {
				$input.prop('checked', false);
			}
		});
		$("#editorContent").on('mouseleave', 'pre', function(){
			var $raw = $(this).find('.toggle-raw');
			$raw.remove();
		});
		$("#editorContent").on('change', '.toggle-raw input', function(){
			var checked = $(this).prop('checked');
			var $pre = $(this).closest('pre');
			if (checked) {
				// 转成ace
				me.preToAce($pre, true);
			} else {
				me.aceToPre($pre, true);
			}
		});

		// 当ace里没有内容时, 连续删除则把ace remove掉
		// keydown的delete事件没有
		var lastDeleteTime;
		$("#editorContent").on('keyup', 'pre',  function(e) {
			var keyCode = e.keyCode;
			// console.log('keyup');
			if(keyCode == 8 || keyCode == 46) { // BackSpace || Delete
				// console.log('delete');
				if(!lastDeleteTime) {
					lastDeleteTime = (new Date()).getTime();
				}
				else {
					var now = (new Date()).getTime();
					if(now - lastDeleteTime < 300) { // 间隔时间很短
						var inAce = me.isInAce($(this))
						if(inAce && !inAce[0].getValue()) {
							// console.log('destroy');
							inAce[0].destroy();
							$(this).remove();
							return;
						}
					}
					lastDeleteTime = now;
				}
				// console.log($(this));
			}
		});
	}
};

function initLeanoteIfrPlugin () {
	// 如果在iframe下, 很可能是嵌入了leanote
	if (self != window.parent) {
		LEA.topInfo = {};
		// 收到消息
		window.addEventListener('message', function(e) {
			console.log('child 收到消息: ')
			console.log(e.data);
			LEA.topInfo = e.data || {};
			LEA.topInfo.got = true;
		}, false);
		if (window.parent.postMessage) {
			window.parent.postMessage('leanote', '*');
		}
	}
}

// 通过src得到note
function getNoteBySrc(src, callback) {
	ajaxGet('/note/getNoteAndContentBySrc', {src: src}, function (ret) {
		if (ret && ret.Ok) {
			var data = ret.Item;
			if (data) {
				var noteInfo = data.NoteInfo;
				var contentInfo = data.NoteContentInfo;
				for (var i in contentInfo) {
					noteInfo[i] = contentInfo[i];
				}
				callback(noteInfo);
			}
			else {
				callback();
			}
		}
		else {
			callback();
		}
	});
}

// 得到top的info's src
var _topInfoStart = (new Date()).getTime();
function getTopInfoSrc (callback) {
	if (LEA.topInfo.got) {
		return callback(LEA.topInfo.src);
	}
	else {
		// 超过1000ms, 不行
		if ((new Date()).getTime() - _topInfoStart > 2000) {
			return callback();
		}
		setTimeout(function () {
			getTopInfoSrc(callback);
		}, 10);
	}
}

// note.html调用
// 实始化页面
function initPage() {
	initLeanoteIfrPlugin();
	if (LEA.topInfo) {
		getTopInfoSrc(function (src) {
			if (src) {
				getNoteBySrc (src, function (srcNote) {
					_initPage(srcNote, true);
				});
			} else {
				_initPage(false, true);
			}
		});
	}
	else {
		_initPage();
	}
}

function _initPage(srcNote, isTop) {
	if (srcNote) {
		curNoteId = srcNote.NoteId;
		curNotebookId = srcNote.NotebookId;
		noteContentJson = srcNote; // 当前笔记变成我的
	}
	else if(isTop) {
		curNoteId = null;
	}

	Notebook.renderNotebooks(notebooks);
	Share.renderShareNotebooks(sharedUserInfos, shareNotebooks);
	
	// 如果初始打开的是共享的笔记
	// 那么定位到我的笔记
	if(curSharedNoteNotebookId) {
		Share.firstRenderShareNote(curSharedUserId, curSharedNoteNotebookId, curNoteId);
	// 初始打开的是我的笔记
	} else {
		Note.setNoteCache(noteContentJson);
		// 判断srcNote是否在notes中
		var isExists = false;
		if (isTop && srcNote && notes) {
			for (var i = 0; i < notes.length; ++i) {
				var note = notes[i];
				if (note.NoteId === srcNote.NoteId) {
					isExists = true;
					notes.splice(i, 1);
					notes.unshift(srcNote);
					break;
				}
			}
			if (!isExists) {
				notes.unshift(srcNote);
			}
		}

		Note.renderNotes(notes);

		if(curNoteId) {
			// 指定某个note时才target notebook, /note定位到最新
			// ie10&+要setTimeout
			setTimeout(function() {
				Note.changeNoteForPjax(curNoteId, true, curNotebookId);
				if (isTop) {
					Note.toggleWriteable();
					setTimeout(function () {
						Note.toggleWriteable();
					}, 100);
					// 如果是markdown
					setTimeout(function () {
						Note.toggleWriteable();
					}, 1000);
				}
			});
			if(!curNotebookId) {
				Notebook.selectNotebook($(tt('#notebook [notebookId="?"]', Notebook.allNotebookId)));
			}
		}
	}

	// 指定笔记, 也要保存最新笔记
	if(latestNotes.length > 0) {
		for(var i = 0; i < latestNotes.length; ++i) {
			Note.addNoteCache(latestNotes[i]);
		}
	}
	
	Tag.renderTagNav(tagsJson);
	// init notebook后才调用
	initSlimScroll();

	LeaAce.handleEvent();

	// 如果是插件, 则切换到编辑页, 并切换到写作模式
	if (isTop) {
		Mobile.toEditor();

		// 如果没有, 则新建之
		if (!srcNote) {
			Note.newNote();
			Note.toggleWriteable(true);
			setTimeout(function () {
				Note.toggleWriteable(true);
			}, 100);
			// 如果是markdown
			setTimeout(function () {
				Note.toggleWriteable(true);
			}, 1000);
		}
	}

	hideMask();
}