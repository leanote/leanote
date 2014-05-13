// 主页渲染
//-------------
// 编辑器模式
var em = new editorMode();

// ifr 的高度, 默认是小20px, 启动1s后运行resizeEditor()调整之

// 鼠标拖动改变宽度
var lineMove = false;
var target = null;
function stopResize3Columns() {
	if (lineMove) {
		// ajax保存
		ajaxGet("/user/updateColumnWidth", {notebookWidth: UserInfo.NotebookWidth, noteListWidth: UserInfo.NoteListWidth}, function() {
		});
	}
	
	lineMove = false;
	$(".noteSplit").css("background", "none");
}

// 最终调用该方法
function resize3ColumnsEnd(notebookWidth, noteListWidth) {
	if(notebookWidth < 150 || noteListWidth < 100) {
//		return;
	}
	var noteWidth = $("body").width() - notebookWidth - noteListWidth;
	if(noteWidth < 400) {
//		return;
	}
	
	$("#leftNotebook").width(notebookWidth);
	$("#notebookSplitter").css("left", notebookWidth);
	
	$("#noteAndEditor").css("left", notebookWidth);
	$("#noteList").width(noteListWidth);
	$("#noteSplitter").css("left", noteListWidth);
	$("#note").css("left", noteListWidth);
	
	UserInfo.NotebookWidth = notebookWidth;
	UserInfo.NoteListWidth = noteListWidth;
}

function resize3Columns(event, isFromeIfr) {
	if (isFromeIfr) {
		event.clientX += $("body").width() - $("#note").width();
	}
	
	var notebookWidth, noteListWidth;

	if (lineMove == true) {
		if (target == "notebookSplitter") {
			notebookWidth = event.clientX;
			noteListWidth = $("#noteList").width();
			resize3ColumnsEnd(notebookWidth, noteListWidth);
		} else {
			notebookWidth = $("#leftNotebook").width();
			noteListWidth = event.clientX - notebookWidth;
			resize3ColumnsEnd(notebookWidth, noteListWidth);
		}

		resizeEditor();
	}
}

// editor

$(function() {
	// 高度设置
//	$("#editor").css("top", $("#noteTop").height());

	$(".noteSplit").bind("mousedown", function(event) {
		event.preventDefault(); // 防止选择文本
		lineMove = true;
		$(this).css("background-color", "#ccc");
		target = $(this).attr("id");

		// 防止iframe捕获不了事件
		$("#noteMask").css("z-index", 99999); // .css("background-color",
											// "#ccc");
	});

	$("body").bind("mouseup", function(event) {
		stopResize3Columns();
		// 取消遮罩
		$("#noteMask").css("z-index", -1);
	});

	$("body").bind("mousemove", function(event) {
		if(lineMove) { // 如果没有这个if会导致不能选择文本
			event.preventDefault();
			resize3Columns(event);
		}
	});

	// toolbar 下拉扩展, 也要resizeEditor
	$("#moreBtn").click(function() {
		saveBookmark();
		
		var height = $("#mceToolbar").height();

		// 现在是折叠的
		if (height < $("#popularToolbar").height()) {
			$("#mceToolbar").height($("#popularToolbar").height());
			$(this).find("i").removeClass("fa-angle-down").addClass("fa-angle-up");
			
		} else {
			$("#mceToolbar").height(height/2);
			$(this).find("i").removeClass("fa-angle-up").addClass("fa-angle-down");
		}
		
		/*
		// 新加 3.12
		var mceToolbarHeight = $("#mceToolbar").height();
		$("#editorContent").css("top", mceToolbarHeight);
		
		// 新加3/22
		$("#leanoteNav").css("top", mceToolbarHeight + 2);

		$("#editor").css("top", $("#noteTop").height());
		*/
		
		resizeEditor();
		
		restoreBookmark();
	});

	// 窗口缩放时
	$(window).resize(function() {
		resizeEditor();
	});

	// 左侧, folder 展开与关闭
	$(".folderHeader").click(
			function() {
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
		selector : "#editorContent",
		// height: 100,//这个应该是文档的高度, 而其上层的高度是$("#content").height(),
		// parentHeight: $("#content").height(),
		content_css : ["css/bootstrap.css", "css/editor/editor.css"].concat(em.getWritingCss()),
		skin : "custom",
		language: LEA.locale, // 语言
		plugins : [
				"autolink link leanote_image lists charmap hr", "paste",
				"searchreplace leanote_nav leanote_code tabfocus",
				"table directionality textcolor codemirror" ], // nonbreaking
				
		toolbar1 : "formatselect | forecolor backcolor | bold italic underline strikethrough | leanote_image | leanote_code | bullist numlist | alignleft aligncenter alignright alignjustify",
		toolbar2 : "outdent indent blockquote | link unlink | table | hr removeformat | subscript superscript |searchreplace | code | pastetext | fontselect fontsizeselect",

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
});

// ie下拒绝访问
// 有兼容性问题
// 不能设置iframe src
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
		/*
		$(target).prepend(
				'<a class="r-' + random + '" name="' + random + '"></a>')
		$("#editorContent_ifr").attr("src", "#" + random);
		iframe.find(".r-" + random).remove();
		*/
	}
}

$(function() {
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
		showDialog("dialogSetInfo", {title: "帐户设置", postShow: function() {
			$('#myTabs a').eq(whichTab).tab('show');
			$("#username").val(UserInfo.Username);
		}});
	}
	
	// 帐号设置
	$("#setInfo").click(function() {
		if(UserInfo.Email) {
			openSetInfoDialog(0);
		} else {
			showDialog("thirdDialogSetInfo", {title: "帐户设置", postShow: function() {
				$('#thirdMyTabs a').eq(0).tab('show');
			}});
		}
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
	
	//--------------
	// 第三方账号设置
	$("#leanoteDialog").on("click", "#accountBtn", function(e) {
		e.preventDefault();
		var email = $("#thirdEmail").val();
		var pwd = $("#thirdPwd").val();
		var pwd2 = $("#thirdPwd2").val();
		if(!email) {
			showAlert("#thirdAccountMsg", "请输入邮箱", "danger", "#thirdEmail");
			return;
		} else {
			var myreg = /^([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+@([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+\.[a-zA-Z]{2,3}$/;
			if(!myreg.test(email)) {
				showAlert("#thirdAccountMsg", "请输入正确的邮箱", "danger", "#thirdEmail");
				return;
			}
		}
		if(!pwd) {
			showAlert("#thirdAccountMsg", "请输入密码", "danger", "#thirdPwd");
			return;
		} else {
			if(pwd.length < 6) {
				showAlert("#thirdAccountMsg", "密码长度至少6位", "danger", "#thirdPwd");
				return;
			}
		}
		if(!pwd2) {
			showAlert("#thirdAccountMsg", "请重复输入密码", "danger", "#thirdPwd2");
			return;
		} else {
			if(pwd != pwd2) {
				showAlert("#thirdAccountMsg", "两次密码输入不一致", "danger", "#thirdPwd2");
				return;
			}
		}
		
		hideAlert("#thirdAccountMsg");
		post("/user/addAccount", {email: email, pwd: pwd}, function(ret) {
			if(ret.Ok) {
				showAlert("#thirdAccountMsg", "添加成功!", "success");
				UserInfo.Email = email;
				$("#curEmail").html(email);
				hideDialog(1000);
			} else {
				showAlert("#thirdAccountMsg", ret.Msg || "添加失败!", "danger");
			}
		}, this);
	});
	
	//-------------
	$("#leanoteDialog").on("click", "#usernameBtn", function(e) {
		e.preventDefault();
		var username = $("#leanoteDialog #username").val();
		if(!username) {
			showAlert('#usernameMsg', "请输入用户名", "danger");
			return;
		} else if(username.length < 4) {
			showAlert('#usernameMsg', "用户名长度至少4位", "danger");
			return;
		} else if(/[^0-9a-zzA-Z_\-]/.test(username)) {
			// 是否含特殊字段?
			showAlert('#usernameMsg', "用户名不能含除数字,字母之外的字符", "danger");
			return;
		}
		hideAlert("#usernameMsg");
		post("/user/updateUsername", {username: username}, function(ret) {
			if(ret.Ok) {
				UserInfo.UsernameRaw = username;
				UserInfo.Username = username.toLowerCase();
				$(".username").html(username);
				showAlert('#usernameMsg', "用户名修改成功!", "success");
			} else {
				showAlert('#usernameMsg', re.Msg || '该用户名已存在', "danger");
			}
		}, "#usernameBtn");
		
	});
	
	// 修改邮箱
	$("#leanoteDialog").on("click", "#emailBtn", function(e) {
		e.preventDefault();
		var email = isEmailFromInput("#email", "#emailMsg");
		if(!email) {
			return;
		}
		
		hideAlert("#emailMsg");
		post("/user/updateEmailSendActiveEmail", {email: email}, function(e) {
			if(e.Ok) {
				var url = getEmailLoginAddress(email);
				showAlert("#emailMsg", "验证邮件已发送, 请及时查阅邮件并验证. <a href='" + url + "' target='_blank'>立即验证</a>", "success");
			} else {
				showAlert("#emailMsg", e.Msg || "邮件发送失败", "danger");
			}
		}, "#emailBtn");
	});
	
	// 修改密码
	$("#leanoteDialog").on("click", "#pwdBtn", function(e) {
		e.preventDefault();
		var oldPwd = $("#oldPwd").val();
		var pwd = $("#pwd").val();
		var pwd2 = $("#pwd2").val();
		
		if(!oldPwd) {
			showAlert("#pwdMsg", "请输入旧密码", "danger", "#oldPwd");
			return;
		} else {
			if(oldPwd.length < 6) {
				showAlert("#pwdMsg", "密码长度至少6位", "danger", "#oldPwd");
				return;
			}
		}
		if(!pwd) {
			showAlert("#pwdMsg", "请输入新密码", "danger", "#pwd");
			return;
		} else {
			if(pwd.length < 6) {
				showAlert("#pwdMsg", "密码长度至少6位", "danger", "#pwd");
				return;
			}
		}
		if(!pwd2) {
			showAlert("#pwdMsg", "请重复输入新密码", "danger", "#pwd2");
			return;
		} else {
			if(pwd != pwd2) {
				showAlert("#pwdMsg", "两次密码输入不一致", "danger", "#pwd2");
				return;
			}
		}
		
		hideAlert("#pwdMsg");
		post("/user/updatePwd", {oldPwd: oldPwd, pwd: pwd}, function(e) {
			if(e.Ok) {
				showAlert("#pwdMsg", "修改密码成功", "success");
			} else {
				showAlert("#pwdMsg", e.Msg, "danger");
			}
		}, "#pwdBtn");
	});
	
	//-------------
	//-------------
	// 邮箱验证
	if(!UserInfo.Verified) {
//		$("#leanoteMsg").hide();
//		$("#verifyMsg").show();
	}
	
	// 帐号设置
	$("#wrongEmail").click(function() {
		openSetInfoDialog(1);
	});
	
	// 重新发送
	$("#leanoteDialog").on("click", ".reSendActiveEmail", function() {
		// 弹框出来
		showDialog("reSendActiveEmailDialog", {title: "发送验证邮件", postShow: function() {
			ajaxGet("/user/reSendActiveEmail", {}, function(ret) {
				if (typeof ret == "object" && ret.Ok) {
					$("#leanoteDialog .text").html("发送成功!")
					$("#leanoteDialog .viewEmailBtn").removeClass("disabled");
					$("#leanoteDialog .viewEmailBtn").click(function() {
						hideDialog();
						var url = getEmailLoginAddress(UserInfo.Email);
						window.open(url, "_blank");
					});
				} else {
					$("#leanoteDialog .text").html("发送失败")
				}
			});
		}});
	});
	
	// 现在去验证
	$("#leanoteDialog").on("click", ".nowToActive", function() {
		var url = getEmailLoginAddress(UserInfo.Email);
		window.open(url, "_blank");
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
		
		if(save) {
			updateLeftIsMin(false);
		}
	}
	
	$("#leftSwitcher2").click(function() {
		maxLeft(true);
	});
	$("#leftSwitcher").click(function() {
		minLeft(true);
		/*
		if(!$("#notebook").is(":hidden")) {
		} else {
			maxLeft(true);
		}
		*/
	});
	
	// 得到最大dropdown高度
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
	$("#notebookMin div.minContainer").hover(function() {
			var target = $(this).attr("target");
			// show的时候要计算高度, 防止过高
			// 先show再计算, 不然高度有偏差
			$(this).find("ul").html($(target).html()).show().height(getMaxDropdownHeight(this));
		}, function() {
			$(this).find("ul").hide();
		}
	);
	
	//------------------------
	// 界面设置, 左侧是否是隐藏的
	UserInfo.NotebookWidth = UserInfo.NotebookWidth || $("#notebook").width();
	UserInfo.NoteListWidth = UserInfo.NoteListWidth || $("#noteList").width();
	if(LEA.isMobile) {
		UserInfo.NoteListWidth = 101;
	}
//	resize3ColumnsEnd(UserInfo.NotebookWidth, UserInfo.NoteListWidth);
	if (UserInfo.LeftIsMin) {
		minLeft(false);
	}
	
	// end
	$("#mainMask").html("");
	$("#mainMask").hide(100);
	
	// 4/25 防止dropdown太高
	// dropdown
	$('.dropdown').on('shown.bs.dropdown', function () {
		var $ul = $(this).find("ul");
		$ul.height(getMaxDropdownHeight(this));
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
	
	// slimScroll
	//---
	setTimeout(function() {
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
	}, 10);
	
	// 编辑器模式
	em.init();
});

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
		toggleA = $("#toggleEditorMode a");
		if(isWriting) {
			toggleA.attr("href", self.normalHash).text("普通模式");
		} else {
			toggleA.attr("href", self.writingHash).text("写作模式");
		}	
	}, 0);
}
editorMode.prototype.isWriting = function(hash) {
	return hash == this.writingHash;
}
editorMode.prototype.init = function() {
	this.changeMode(this.isWritingMode);
	var self = this;
	$("#toggleEditorMode").click(function() {
		var $a = $(this).find("a");
		var isWriting = self.isWriting($a.attr("href"));
		self.changeMode(isWriting);
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
}
editorMode.prototype.normalMode = function() {
	var $c = $("#editorContent_ifr").contents();
	
	$c.contents().find("#writtingMode").remove();
	$c.contents().find('link[href$="editor-writting-mode.css"]').remove();
			
		$("#noteItemListWrap, #notesAndSort").show();
	$("#noteList").unbind("mouseenter").unbind("mouseleave"); 
	
	var theme = UserInfo.Theme || "default";
	theme += ".css";
	$("#themeLink").attr("href", "/css/theme/" + theme);
	
	
}
editorMode.prototype.writtingMode = function() {
//	$("body").fadeOut();
	
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
	
//	$("body").fadeIn();
}

editorMode.prototype.getWritingCss = function() {
	if(this.isWritingMode) {
		return ["css/editor/editor-writting-mode.css"];
	}
	return [];
}