//------------------------------------
// 共享, notbeook, note
//------------------------------------

// 默认共享notebook id
Share.defaultNotebookId = "share0";
Share.defaultNotebookTitle = "默认共享";
Share.sharedUserInfos = {}; // userId => {}

// 在render时就创建, 以后复用之
Share.userNavs = {}; // userId => {"forList":html, "forNew":html}

// 缓存都不要, 统一放在Note.cache中
// 放在这里, 只是为了debug, 分离
Share.notebookCache = {}; // notebooks 的cache
Share.cache = {}; // note的cache

// 分享的弹出框是note的
Share.dialogIsNote = true;

// 设置缓存 note
Share.setCache = function(note) {
	if(!note || !note.NoteId) {
		return;
	}
	Share.cache[note.NoteId] = note;
}

/**
 * 我的共享notebooks	    
 <div id="shareNotebooks">
	 <div class="folderNote closed">
      <div class="folderHeader">
        <a>
          <h1>
            <i class="fa fa-angle-right"></i> 
            Life's</h1>
        </a>
      </div>
      <ul class="folderBody">
        <li><a>Hadoop</a></li>
        <li><a>Node webkit</a></li>
        <li><a>Hadoop</a></li>
        <li><a>Node webkit</a></li>
      </ul>
    </div>
 */
// TODO 层级
// shareNotebooks = {userId => {}}
Share.renderShareNotebooks = function(sharedUserInfos, shareNotebooks) {
	if(isEmpty(sharedUserInfos)) {
		return;
	}
	
	if(!shareNotebooks || typeof shareNotebooks != "object" || shareNotebooks.length < 0) {
		return;
	}
	
	var $shareNotebooks = $("#shareNotebooks");

	// render每一个用户的share给我的笔记本, 之前先建一个默认共享
	for(var i in sharedUserInfos) {
		var userInfo = sharedUserInfos[i];
		var userNotebooks = shareNotebooks[userInfo.UserId] || [];
		
		userNotebooks = [{NotebookId: Share.defaultNotebookId, Title: Share.defaultNotebookTitle}].concat(userNotebooks)

		var username = userInfo.Username || userInfo.Email;
		userInfo.Username = username;
		Share.sharedUserInfos[userInfo.UserId] = userInfo;
		var userId = userInfo.UserId;
		var header = t('<li class="each-user"><div class="" fromUserId="?"><i class="fa fa-angle-down"></i><span>?</span></div>', userInfo.UserId, username);
		var friendId = "friendContainer" + i;
		var body = '<ul class="" id="' + friendId + '">';
		
		var forList = ""; // 全部
		var forNew = ""; // 必须要有权限的
		
		for(var j in userNotebooks) {
			var notebook = userNotebooks[j];
			
			// 缓存起来, 像Note
			notebook.IsShared = true;
			notebook.UserId = userId;
			Share.notebookCache[notebook.NotebookId] = notebook;
			// notebook的cache也缓存一份, 为了显示标题
			Notebook.cache[notebook.NotebookId] = notebook;
			
			body += t('<li><a notebookId="?" fromUserId="?">?</a></li>', notebook.NotebookId, userId, notebook.Title)
			
			// 
			var each = t('<li role="presentation"><a role="menuitem" tabindex="-1" href="#" userId="?" notebookId="?">?</a></li>', userId, notebook.NotebookId, notebook.Title);
			forList += each;
			if(j != 0 && notebook.Perm) {
				forNew += t('<li role="presentation" class="clearfix" userId="?" notebookId="?"><div class="new-note-left pull-left">?</div><div class="new-note-right pull-left">Markdown</div></li>', userId, notebook.NotebookId, notebook.Title);
			}
		}
		
		body += "</ul>";
		Share.userNavs[userId] = {"forList": forList, "forNew": forNew};
		
		$shareNotebooks.append(header + body + "</div>")
		
		// mainShare
		// $("#minShareNotebooks").append('<div class="minContainer" target="#' + friendId + '" title="' + username + ' 的分享"><i class="fa fa-user"></i><ul class="dropdown-menu"></ul></li>')
	}
};

Share.isDefaultNotebookId = function(notebookId) {
	return Share.defaultNotebookId == notebookId;
}

// 转成共享的nav
// for list和for new
// 如果forNew没有, 那么还是保持我的nav
Share.toggleToSharedNav = function(userId, notebookId) {
	// for list
	$("#sharedNotebookNavForListNote").html(Share.userNavs[userId].forList);
	$("#sharedNotebookNavForListNav").show();
	$("#curSharedNotebookForListNote").html(Share.notebookCache[notebookId].Title + '(' + Share.sharedUserInfos[userId].Username + ")");
	$("#myNotebookNavForListNav").hide();
	
	// for new
	var forNew = Share.userNavs[userId].forNew;
	if(forNew) {
		$("#notebookNavForNewSharedNote").html(forNew);
		// 新建之, 可能当前选择的没有权限新建. 此时需要得到第一个
		var curNotebookId = "";
		var curNotebookTitle = "";
		if(Share.notebookCache[notebookId].Perm) {
			curNotebookId = notebookId;
			curNotebookTitle = Share.notebookCache[notebookId].Title;
		} else {
			// 得到第一个
			var $f = $("#notebookNavForNewSharedNote li").eq(0);
			curNotebookId = $f.attr("notebookId");
			curNotebookTitle = $f.text();
		}
		
		$("#curNotebookForNewSharedNote").html(curNotebookTitle + '(' + Share.sharedUserInfos[userId].Username + ')');
		$("#curNotebookForNewSharedNote").attr("notebookId", curNotebookId);
		$("#curNotebookForNewSharedNote").attr("userId", userId);
		
		$("#newSharedNote").show();
		$("#newMyNote").hide();
	}
	
	// 隐藏tag
	$("#tagSearch").hide();
}

//改变笔记本
//0. 改变样式
//1. 改变note, 此时需要先保存
//2. ajax得到该notebook下的所有note
//3. 使用Note.RederNotes()
Share.changeNotebook = function(userId, notebookId) {
	// 选中
	Notebook.selectNotebook($(t('#shareNotebooks a[notebookId="?"]', notebookId)));
	
	// 改变nav!!!! TODO
	Share.toggleToSharedNav(userId, notebookId);
	
	// 1
	Note.curChangedSaveIt();
	
	// 2 先清空所有
	Note.clearAll();
	
	var url = "/share/ListShareNotes/";
	var param = {userId: userId};
	if(!Share.isDefaultNotebookId(notebookId)) {
		param.notebookId = notebookId;
	}
	
	// 2 得到笔记本
	// 这里可以缓存起来, note按notebookId缓存
	ajaxGet(url, param, function(ret) {
		// note 导航
		// 
		// 如果是特定笔记本下的notes, 那么传过来的没有权限信息, 此时权限由notebookId决定
		if(param.notebookId) {
			
		}
		Note.renderNotes(ret, false, true);
		// 渲染第一个
		// 这里, 有点小复杂, 还要判断权限...
		if(!isEmpty(ret)) {
			// 定位
			Note.changeNote(ret[0].NoteId, true);
		} else {
		}
	});
}

// 是否有更新权限
// called by Note
Share.hasUpdatePerm = function(notebookId) {
	var note = Share.cache[notebookId];
	if(!note || !note.Perm) {
		return false;
	}
	return true;
}

//---------------------------
// 我删除别人共享给我的笔记本
Share.deleteShareNotebook = function(target) {
	var notebookId = $(target).attr("notebookId");
	var fromUserId = $(target).attr("fromUserId"); // 谁共享给了我 from
	ajaxGet("/share/DeleteShareNotebookBySharedUser", {notebookId: notebookId, fromUserId: fromUserId}, function(ret) {
		if(ret) {
			$(target).parent().remove();
		}
	});
}
Share.deleteShareNote = function(target) {
	var noteId = $(target).attr("noteId");
	var fromUserId = $(target).attr("fromUserId"); // 谁共享给了我 from
	ajaxGet("/share/DeleteShareNoteBySharedUser", {noteId: noteId, fromUserId: fromUserId}, function(ret) {
		if(ret) {
			$(target).remove();
		}
	});
}
Share.deleteUserShareNoteAndNotebook = function(target) {
	var fromUserId = $(target).attr("fromUserId"); // 谁共享给了我 from
	ajaxGet("/share/deleteUserShareNoteAndNotebook", {fromUserId: fromUserId}, function(ret) {
		if(ret) {
			$(target).parent().remove();
		}
	});
}

// 新建shared note
Share.changeNotebookForNewNote = function(notebookId) {
	// 改变nav for list, for new
	Notebook.selectNotebook($(t('#shareNotebooks [notebookId="?"]', notebookId)));
	var userId = Share.notebookCache[notebookId].UserId;
	Share.toggleToSharedNav(userId, notebookId);	
	
	// 得到笔记本
	var url = "/share/ListShareNotes/";
	var param = {userId: userId, notebookId: notebookId};
		
	// 2 得到笔记本
	// 这里可以缓存起来, note按notebookId缓存
	ajaxGet(url, param, function(ret) {
		// note 导航
		Note.renderNotes(ret, true, true);
	});
}

// 删除笔记, 我有权限, 且是我创建的笔记
Share.deleteSharedNote = function(target, contextmenuItem) {
	Note.deleteNote(target, contextmenuItem, true);
}
Share.copySharedNote = function(target, contextmenuItem) {
	Note.copyNote(target, contextmenuItem, true);
}

Share.contextmenu = null;
Share.initContextmenu = function() {
	if(Share.contextmenu) {
		Share.contextmenu.unbind("contextmenu");
	}
	// 得到可移动的notebook
	var notebooksCopy = [];
	
	// 到时这个可以缓存起来
	$("#notebookNavForNewNote li .new-note-left").each(function() {
		var notebookId = $(this).attr("notebookId");
		var title = $(this).text();
		var copy = {text: title, notebookId: notebookId, action: Share.copySharedNote}
		notebooksCopy.push(copy);
	});
	
	//---------------------
	// context menu
	//---------------------
	var noteListMenu = {
		width: 170, 
		items: [
			{ text: "复制到我的笔记本", alias: "copy", icon: "",
				type: "group", 
				width: 150, 
				items: notebooksCopy
			},
			{ type: "splitLine" },
			{ text: "删除", alias: "delete", icon: "", faIcon: "fa-trash-o", action: Share.deleteSharedNote }
		], 
		onShow: applyrule,
		parent: "#noteItemList",
		children: ".item-shared",
	}
	function applyrule(menu) {
		var noteId = $(this).attr("noteId");
		var note = Share.cache[noteId];
		if(!note) {
			return;
		}
		var items = [];
		if(!(note.Perm && note.CreatedUserId == UserInfo.UserId)) {
			items.push("delete");
		}
		// 不是自己的创建的不能删除
        menu.applyrule({
        	name: "target...",
            disable: true,
            items: items
        });		
	   
	}
	
	Share.contextmenu = $("#noteItemList .item-shared").contextmenu(noteListMenu);
}

$(function() {
	// 点击notebook
	$("#shareNotebooks").on("click", "ul li a", function() {
		var notebookId = $(this).attr("notebookId");
		var userId = $(this).attr("fromUserId");
		Share.changeNotebook(userId, notebookId);
	});
	// min
	$("#minShareNotebooks").on("click", "li", function() {
		var self = $(this).find("a");
		var notebookId = $(self).attr("notebookId");
		var userId = $(self).attr("fromUserId");
		Share.changeNotebook(userId, notebookId);
	});
	
	//-----------------------------
	// contextmenu shareNotebooks
	// 删除共享笔记本
	var shareNotebookMenu = {
			width: 150, 
			items: [
				{ text: "删除共享笔记本", icon: "", faIcon: "fa-trash-o", action: Share.deleteShareNotebook }
			], 
			onShow: applyrule,
			onContextMenu: beforeContextMenu,
			
			parent: "#shareNotebooks .folderBody",
			children: "li a",
	};
	function applyrule(menu) {
		return;
	}
	// 默认共享不能删除
	function beforeContextMenu() {
		var notebookId = $(this).attr("notebookId");
		return !Share.isDefaultNotebookId(notebookId);
	}
	
	$("#shareNotebooks").contextmenu(shareNotebookMenu);
	
	//---------------------------
	// contextmenu shareNotebooks
	// 删除某用户所有的
	var shareUserMenu = {
			width: 150, 
			items: [
				{ text: "删除所有共享", icon: "", faIcon: "fa-trash-o", action: Share.deleteUserShareNoteAndNotebook }
			],
			parent: "#shareNotebooks",
			children: ".folderHeader",
	};
	
	$("#shareNotebooks").contextmenu(shareUserMenu);
	
	//---------------------------
	// 新建笔记
	// 1. 直接点击新建 OR
	// 2. 点击nav for new note
	$("#newSharedNoteBtn").click(function() {
		var notebookId = $("#curNotebookForNewSharedNote").attr('notebookId');
		var userId = $("#curNotebookForNewSharedNote").attr('userId');
		Note.newNote(notebookId, true, userId);
	});
	$("#newShareNoteMarkdownBtn").click(function() {
		var notebookId = $("#curNotebookForNewSharedNote").attr('notebookId');
		var userId = $("#curNotebookForNewSharedNote").attr('userId');
		Note.newNote(notebookId, true, userId, true);
	});
	$("#notebookNavForNewSharedNote").on("click", "li div", function() {
		var notebookId = $(this).parent().attr("notebookId");
		var userId = $(this).parent().attr("userId");
		
		if($(this).text() == "Markdown") {
			Note.newNote(notebookId, true, userId, true);
		} else {
			Note.newNote(notebookId, true, userId);
		}
	});
	
	//------------------
	Share.initContextmenu();
	
	//------------------
	// 添加共享
	$("#leanoteDialogRemote").on("click", ".change-perm", function() {
		var self = this;
		var perm = $(this).attr("perm");
		var noteOrNotebookId = $(this).attr("noteOrNotebookId");
		var toUserId = $(this).attr("toUserId");
		var toHtml = "可编辑";
		var toPerm = "1";
		if(perm == "1") {
			toHtml = "只读";
			toPerm = "0";
		}
		var url = "/share/UpdateShareNotebookPerm";
		var param = {perm: toPerm, toUserId: toUserId};
		if(Share.dialogIsNote) {
			url = "/share/UpdateShareNotePerm";
			param.noteId = noteOrNotebookId;
		} else {
			param.notebookId = noteOrNotebookId;
		}
		ajaxGet(url, param, function(ret) {
			if(ret) {
				$(self).html(toHtml);
				$(self).attr("perm", toPerm);
			}
		});
	});
	
	$("#leanoteDialogRemote").on("click", ".delete-share", function() {
		var self = this;
		var noteOrNotebookId = $(this).attr("noteOrNotebookId");
		var toUserId = $(this).attr("toUserId");
		
		var url = "/share/DeleteShareNotebook";
		var param = {toUserId: toUserId};
		if(Share.dialogIsNote) {
			url = "/share/DeleteShareNote";
			param.noteId = noteOrNotebookId;
		} else {
			param.notebookId = noteOrNotebookId;
		}	
		
		ajaxGet(url, param, function(ret) {
			if(ret) {
				$(self).parent().parent().remove();
			}
		});
	});
	
	// 添加共享
	var seq = 1;
	$("#leanoteDialogRemote").on("click", "#addShareNotebookBtn", function() {
		seq++;
		var tpl = '<tr id="tr' + seq + '"><td>#</td><td><input id="friendsEmail" type="text" class="form-control" style="width: 200px" placeholder="好友邮箱"/></td>';
		tpl += '<td><label for="readPerm' + seq + '"><input type="radio" name="perm' + seq + '" checked="checked" value="0" id="readPerm' + seq + '"> 只读</label>';
		tpl += ' <label for="writePerm' + seq + '"><input type="radio" name="perm' + seq + '" value="1" id="writePerm' + seq + '"> 可编辑</label></td>';
		tpl += '<td><button class="btn btn-success" onclick="addShareNoteOrNotebook(' + seq + ')">分享</button>';
		tpl += ' <button class="btn btn-warning" onclick="deleteShareNoteOrNotebook(' + seq + ')">删除</button>';
		tpl += "</td></tr>";
		$("#shareNotebookTable tbody").prepend(tpl);
		
		$("#tr" + seq + " #friendsEmail").focus();
	});
	
	//-------------------
	// 发送邀请邮件
	$("#registerEmailBtn").click(function() {
		var content = $("#emailContent").val();
		var toEmail = $("#toEmail").val();
		if(!content) {
			showAlert("#registerEmailMsg", "邮件内容不能为空", "danger");
			return;
		}
		post("/user/sendRegisterEmail", {content: content, toEmail: toEmail}, function(ret) {
			showAlert("#registerEmailMsg", "发送成功!", "success");
			hideDialog2("#sendRegisterEmailDialog", 1000);
		}, this);
	});
});

// trSeq 1,2,3...
function addShareNoteOrNotebook(trSeq) {
	var trId = "#tr" + trSeq;
	var id = Share.dialogNoteOrNotebookId;
	
	var emails = isEmailFromInput(trId + " #friendsEmail", "#shareMsg", "请输入好友邮箱");
	if(!emails) {
		return;
	}
	var shareNotePerm = $(trId + ' input[name="perm' + trSeq + '"]:checked').val() || 0;
	var perm = shareNotePerm;
	// emails = emails.split(";");
	var url = "share/addShareNote";
	var data = {noteId: id, emails: [emails], perm: shareNotePerm};
	if(!Share.dialogIsNote) {
		url = "share/addShareNotebook";
		data = {notebookId: id, emails: [emails], perm: shareNotePerm};
	}
	hideAlert("#shareMsg");
	post(url, data, function(ret) {
		var ret = ret[emails];
		if(ret) {
			// 成功
			// 成功了则去掉输入框
			if(ret.Ok) {
				var tpl = t('<td>?</td>', '#');
				tpl += t('<td>?</td>', emails);
				tpl += t('<td><a href="#" noteOrNotebookId="?" perm="?" toUserId="?" title="点击改变权限" class="btn btn-default change-perm">?</a></td>', id, perm, ret.Id, !perm || perm == '0' ? "只读" : "可编辑");
				tpl += t('<td><a href="#" noteOrNotebookId="?" toUserId="?" class="btn btn-warning delete-share">删除</a></td>', id, ret.Id);
				$(trId).html(tpl);
			} else {
				var shareUrl = 'http://leanote/register?from=' + UserInfo.Username;
				showAlert("#shareMsg", "该用户还没有注册, 复制邀请链接发送给Ta一起来体验leanote, 邀请链接: " + shareUrl + ' <a id="shareCopy"  data-clipboard-target="copyDiv">点击复制</a> <span id="copyStatus"></span> <br /> 或者发送邀请邮件给Ta, <a href="#" onclick="sendRegisterEmail(\'' + emails + '\')">点击发送', "warning");
				$("#copyDiv").text(shareUrl);
				initCopy("shareCopy", function(args) {
					if(args.text) {
						showMsg2("#copyStatus", "复制成功", 1000);
					} else {
						showMsg2("#copyStatus", "对不起, 复制失败, 请自行复制", 1000);
					}
				});
			}
		}
	}, trId + " .btn-success");
}

// 发送邀请邮件
function sendRegisterEmail(email) {
	showDialog2("#sendRegisterEmailDialog", {postShow: function() {
		$("#emailContent").val("Hi, 我是" + UserInfo.Username + ", leanote非常好用, 快来注册吧!");
		setTimeout(function() {
			$("#emailContent").focus();
		}, 500);
		$("#toEmail").val(email);
	}});
}

function deleteShareNoteOrNotebook(trSeq) {
	$("#tr" + trSeq).remove();	
}
