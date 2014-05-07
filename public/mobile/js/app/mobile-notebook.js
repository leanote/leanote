Notebook.curNotebookId = "";
Notebook.cache = {}; // notebookId => {};
Notebook.notebooks = []; // 按次序
// <li role="presentation"><a role="menuitem" tabindex="-1" href="#">CSS</a></li>	
Notebook.notebookNavForListNote = ""; // html 为了note list上面和新建时的ul
Notebook.notebookNavForNewNote = ""; // html 为了note list上面和新建时的ul

// 设置缓存
Notebook.setCache = function(notebook) {
	var notebookId = notebook.NotebookId;
	if(!notebookId) {
		return;
	}
	if(!Notebook.cache[notebookId]) {
		Notebook.cache[notebookId] = {};
	}
	$.extend(Notebook.cache[notebookId], notebook);
}

Notebook.GetCurNotebookId = function() {
	return Notebook.curNotebookId;
};

// 得到notebook标题, 给note显示其notebook标题用
Notebook.GetNotebook = function(notebookId) {
	return Notebook.cache[notebookId];
}
Notebook.GetNotebookTitle = function(notebookId) {
	var notebook = Notebook.cache[notebookId];
	if(notebook) {
		return notebook.Title;
	} else {
		return "未知";
	}
}

/**
 * 我的notebooks
<ul class="folderBody" id="notebookList">
	<li><a class="active">所有</a></li>
	<li><a class="active">Hadoop</a></li>
	<li><a>August 13, 2013</a></li>
</ul>
 */
// TODO 层级
Notebook.allNotebookId = "0";
Notebook.trashNotebookId = "-1";
Notebook.RenderNotebooks = function(notebooks) {
	if(!notebooks || typeof notebooks != "object" || notebooks.length < 0) {
		notebooks = [];
	}
	
	notebooks = [{NotebookId: Notebook.allNotebookId, Title: getMsg("all")}].concat(notebooks);
	notebooks.push({NotebookId: Notebook.trashNotebookId, Title: getMsg("trash")});
	Notebook.notebooks = notebooks; // 缓存之
	
	var $notebookList = $("#notebookList");
	var nav = "";
	for(var i in notebooks) {
		var notebook = notebooks[i];
		Notebook.cache[notebook.NotebookId] = notebook;
		var classes = "";
		if(i == 0) {
			classes = "active";
			Notebook.curNotebookId = notebook.NotebookId;
		}
		$notebookList.append(t('<li><a class="?" notebookId="?">?</a></li>', classes, notebook.NotebookId, notebook.Title))
	}
	
	// 渲染nav
	// Notebook.renderNav();
	
	// 渲染第一个notebook作为当前
	// Notebook.changeNotebookNavForNewNote(notebooks[0].NotebookId);
}

// RenderNotebooks调用, 
// nav 为了新建, 快速选择, 移动笔记
// 这些在添加,修改,删除notebooks都要变动!!!
Notebook.renderNav = function(nav) {
	var navForListNote = "";
	var navForNewNote = "";
	var navForMoveNote = "";
	var len = Notebook.notebooks.length-1;
	var contextmenu = [];
	for(var i in Notebook.notebooks) {
		var notebook = Notebook.notebooks[i];
		var each = t('<li role="presentation"><a role="menuitem" tabindex="-1" href="#" notebookId="?">?</a></li>', notebook.NotebookId, notebook.Title);
		var eachForNew = t('<li role="presentation" class="clearfix"><div class="new-note-left pull-left" title="为该笔记本新建笔记" href="#" notebookId="?">?</div><div title="为该笔记本新建markdown笔记" class="new-note-right pull-left" notebookId="?">Markdown</div></li>', notebook.NotebookId, notebook.Title, notebook.NotebookId);
		navForListNote  += each;
		if(i != 0 && i != len) {
			navForMoveNote += each;
			navForNewNote += eachForNew;
		}
	}
	
	$("#notebookNavForListNote").html(navForListNote);
	$("#notebookNavForNewNote").html(navForNewNote);
	$("#notebookNavForMoveNote").html(navForMoveNote);
}

// 修改,添加,删除notebook后调用
// 改变nav
// 直接从html中取!
Notebook.changeNav = function() {
	var navForListNote = "";
	var navForNewNote = "";
	
	var i = 0;
	var $list = $("#notebookList li a");
	var len = $list.length - 1;
	$list.each(function() {
		var notebookId = $(this).attr("notebookId");
		var notebook = Notebook.cache[notebookId];
		var each = t('<li role="presentation"><a role="menuitem" tabindex="-1" href="#" notebookId="?">?</a></li>', notebook.NotebookId, notebook.Title);
		var eachForNew = t('<li role="presentation" class="clearfix"><div class="new-note-left pull-left" title="为该笔记本新建笔记" href="#" notebookId="?">?</div><div title="为该笔记本新建markdown笔记" class="new-note-right pull-left" notebookId="?">Markdown</div></li>', notebook.NotebookId, notebook.Title, notebook.NotebookId);
			
		navForListNote  += each;
		var isActive = $(this).hasClass('active'); // 万一修改的是已选择的, 那么...
		if(isActive) {
			$("#curNotebookForListNote").html(notebook.Title);
		}
		if(i != 0 && i != len) {
			navForNewNote  += eachForNew;
			if(isActive) {
				$("#curNotebookForNewNote").html(notebook.Title);
			}
		}
		i++;
	});
	
	$("#notebookNavForListNote").html(navForListNote);
	$("#notebookNavForNewNote").html(navForNewNote);
	$("#notebookNavForMoveNote").html(navForNewNote);
	
	// 移动, 复制重新来, 因为nav变了, 移动至-----的notebook导航也变了
	Note.InitContextmenu();
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
Notebook.RenderShareNotebooks = function(sharedUserInfos, shareNotebooks) {
	if(isEmpty(sharedUserInfos)) {
		return;
	}
	
	if(!shareNotebooks || typeof shareNotebooks != "object" || shareNotebooks.length < 0) {
		return;
	}
	
	var $shareNotebooks = $("#shareNotebooks");
	var user2ShareNotebooks = {};
	for(var i in shareNotebooks) {
		var userNotebooks = shareNotebooks[i];
		user2ShareNotebooks[userNotebooks.UserId] = userNotebooks;
	}
	for(var i in sharedUserInfos) {
		var userInfo = sharedUserInfos[i];
		var userNotebooks = user2ShareNotebooks[userInfo.UserId] || {ShareNotebooks:[]};
		
		userNotebooks.ShareNotebooks = [{NotebookId: "-2", Title: "默认共享"}].concat(userNotebooks.ShareNotebooks)

		var username = userInfo.Username || userInfo.Email;
		var header = t('<div class="folderNote closed"><div class="folderHeader"><a><h1 title="? 的共享"><i class="fa fa-angle-right"></i>?</h1></a></div>', username, username);
		var body = '<ul class="folderBody">';
		for(var j in userNotebooks.ShareNotebooks) {
			var notebook = userNotebooks.ShareNotebooks[j];
			body += t('<li><a notebookId="?">?</a></li>', notebook.NotebookId, notebook.Title)
		}
		body += "</ul>";
		
		$shareNotebooks.append(header + body + "</div>")
	}
}

// 左侧导航, 选中某个notebook
Notebook.SelectNotebook = function(target) {
	$("#notebookList li a").removeClass("active");
	$(target).addClass("active");
};

// 新建笔记导航
Notebook.changeNotebookNavForNewNote = function(notebookId, title) {
	// 没有notebookId, 则选择第1个notebook
	// 第一个是全部笔记
	if(!notebookId) {
		var notebook = Notebook.notebooks[0];
		notebookId = notebook.NotebookId;
		title = notebook.Title;
	}
	if(!title) {
		var notebook = Notebook.cache[0];
		title = notebook.Title;
	}
	
	if(!Notebook.IsAllNotebookId(notebookId) && !Notebook.IsTrashNotebookId(notebookId)) {
		$("#curNotebookForNewNote").html(title).attr("notebookId", notebookId);
	} else if(!$("#curNotebookForNewNote").attr("notebookId")) {
		// 但又没有一个笔记, 默认选第一个吧
		// 这里很可能会死循环, 万一用户没有其它笔记呢?
		// 服务端肯定要在新建一个用户时给他创建一个默认笔记本的
		if(Notebook.notebooks.length > 2) {
			var notebook = Notebook.notebooks[1];
			notebookId = notebook.NotebookId;
			title = notebook.Title;
			Notebook.changeNotebookNavForNewNote(notebookId, title);
		}
	}
}

// 改变导航, 两处
// 单击左侧, 单击新建下拉时调用
// 1 选中左侧导航, 
// 2 notelist上面 >
// 3 新建笔记 - js >
// 转成我的nav <-> 共享
Notebook.toggleToMyNav = function(userId, notebookId) {
	$("#sharedNotebookNavForListNav").hide();
	$("#myNotebookNavForListNav").show();
	
	$("#newMyNote").show();
	$("#newShareNote").hide();
	
	// 搜索tag隐藏
	$("#tagSearch").hide();
}
Notebook.changeNotebookNav = function(notebookId) {
	Notebook.toggleToMyNav();
	
	// 1
	Notebook.SelectNotebook($(t('#notebookList [notebookId="?"]', notebookId)));
	
	var notebook = Notebook.cache[notebookId];
	
	if(!notebook) {
		return;
	}
	
	// 2
	$("#curNotebookForListNote").html(notebook.Title);
	
	// 3
	Notebook.changeNotebookNavForNewNote(notebookId, notebook.Title);
}

Notebook.IsAllNotebookId = function(notebookId) {
	return notebookId == Notebook.allNotebookId;
}
Notebook.IsTrashNotebookId = function(notebookId) {
	return notebookId == Notebook.trashNotebookId;
}
// 当前选中的笔记本是否是"所有"
Notebook.CurActiveNotebookIsAll = function() {
	return Notebook.IsAllNotebookId($("#notebookList .active").attr("notebookId"));
}

// 改变笔记本
// 0. 改变样式
// 1. 改变note, 此时需要先保存
// 2. ajax得到该notebook下的所有note
// 3. 使用Note.RederNotes()
Notebook.ChangeNotebook = function(notebookId) {
	Notebook.changeNotebookNav(notebookId);
		
	Notebook.curNotebookId = notebookId;
		
	// 1
	// Note.CurChangedSaveIt();
	
	// 2 先清空所有
	// Note.ClearAll();
	
	var url = "/note/ListNotes/";
	var param = {notebookId: notebookId};
	
	var notebook = Notebook.cache[notebookId]
	toggle("notes", notebook.Title + " 的笔记")
	
	// 废纸篓
	if(Notebook.IsTrashNotebookId(notebookId)) {
		url = "/note/listTrashNotes";
		param = {};
	} else if(Notebook.IsAllNotebookId(notebookId)) {
		param = {};
		// 得到全部的...
		cacheNotes = Note.GetNotesByNotebookId();
		if(!isEmpty(cacheNotes)) { // 万一真的是没有呢?
			Note.RenderNotes(cacheNotes);
			return;
		}
	} else {
		cacheNotes = Note.GetNotesByNotebookId(notebookId);
		if(!isEmpty(cacheNotes)) { // 万一真的是没有呢?
			Note.RenderNotes(cacheNotes);
			return;
		}
	}
	
	// 2 得到笔记本
	// 这里可以缓存起来, note按notebookId缓存
	ajaxGet(url, param, Note.RenderNotes);
	
}

// 是否是当前选中的notebookId
// 还包括共享
Notebook.IsCurNotebook = function(notebookId) {
	return $(t('#notebookList [notebookId="?"], #shareNotebooks [notebookId="?"]', notebookId, notebookId)).attr("class") == "active";
}

// 改变nav, 为了新建note
Notebook.ChangeNotebookForNewNote = function(notebookId) {
	// 废纸篓
	if(Notebook.IsTrashNotebookId(notebookId) || Notebook.IsAllNotebookId(notebookId)) {
		return;
	}
	
	Notebook.changeNotebookNav(notebookId);
	Notebook.curNotebookId = notebookId;
	
	var url = "/note/ListNotes/";
	var param = {notebookId: notebookId};
		
	// 2 得到笔记本
	// 这里可以缓存起来, note按notebookId缓存
	ajaxGet(url, param, function(ret) {
		// note 导航
		Note.RenderNotes(ret, true);
	});
};

//---------------------------
// 显示共享信息
Notebook.ListNotebookShareUserInfo = function(target) {
	var notebookId = $(target).attr("notebookId");
	showDialogRemote("share/listNotebookShareUserInfo", {notebookId: notebookId});
}
// 共享笔记本
Notebook.ShareNotebook = function(target) {
	var title = $(target).text();
	showDialog("dialogShareNote", {title: "分享笔记本给好友-" + title});
	setTimeout(function() {
		$("#friendsEmail").focus();
	}, 500);
	var notebookId = $(target).attr("notebookId");
	
	shareNoteOrNotebook(notebookId, false);
}

//-----------------------------
// 设为blog/unset
Notebook.SetNotebook2Blog = function(target) {
	var notebookId = $(target).attr("notebookId");
	var notebook = Notebook.cache[notebookId];
	var isBlog = true;
	if(notebook.IsBlog != undefined) {
		isBlog = !notebook.IsBlog;
	}
	
	// 那么, 如果当前是该notebook下, 重新渲染之
	if(Notebook.curNotebookId == notebookId) {
		if(isBlog) {
			$("#noteList .item-blog").show();
		} else {
			$("#noteList .item-blog").hide();
		}
		
	// 如果当前在所有笔记本下
	} else if(Notebook.curNotebookId == Notebook.allNotebookId){
		$("#noteItemList .item").each(function(){
			var noteId = $(this).attr("noteId");
			var note = Note.cache[noteId];
			if(note.NotebookId == notebookId) {
				if(isBlog) $(this).find(".item-blog").show();
				else $(this).find(".item-blog").hide();
			}
		});
	}
	ajaxPost("blog/setNotebook2Blog", {notebookId: notebookId, isBlog: isBlog}, function(ret) {
		if(ret) {
			// 这里要设置notebook下的note的blog状态
			Note.setAllNoteBlogStatus(notebookId, isBlog);
			Notebook.setCache({NotebookId: notebookId, IsBlog: isBlog});
		}
	});
}

// 添加, 修改完后都要对notebook的列表重新计算 TODO

// 修改笔记本标题
Notebook.UpdateNotebookTitle = function(target) {
	var notebookTitle = $(target).text();
	var id = "editNotebookTitle";
	$(target).html(t('<input type="text" value="?" everValue="?" id="?" notebookId="?"/>', notebookTitle, notebookTitle, id, $(target).attr("notebookId")));
	$("#" + id).focus();
}
Notebook.DoUpdateNotebookTitle = function() {
	var title = $(this).val();
	var everTitle = $(this).attr("everTitle");
	var notebookId = $(this).attr("notebookId");
	
	if(!title) {
		title = everTitle;
	}
	$(this).parent().html(title);
	
	if(title != everTitle) {
		ajaxPost("/notebook/updateNotebookTitle", {notebookId: notebookId, title: title}, function(ret) {
			// 修改缓存
			Notebook.cache[notebookId].Title = title;
			// 改变nav
			Notebook.changeNav();
		});
	}
}

//-----------
// 添加笔记本
// 1 确保是展开的
// 2 在所有后面添加<li></li>
Notebook.addNotebookSeq = 1; // inputId
Notebook.AddNotebook = function() {
	if($("#myNotebooks").hasClass("closed")) {
		$("#myNotebooks .folderHeader").trigger("click");
	}
	var inputId = "newNotebookInput" + Notebook.addNotebookSeq;
	Notebook.addNotebookSeq++;
	$("#notebookList li").eq(0).after(t('<li><a><input id="?"/></a></li>', inputId));
	
	$("#" + inputId).focus();
	
	// 回车调用blur
	enterBlur("#" + inputId);
	$("#" + inputId).blur(function() {
		// 为防多次发生blur
		$(this).unbind("blur");
		
		var title = $(this).val();
		if(!title) {
			$(this).parent().parent().remove();
		} else {
			// 添加之
			var notebookId = getObjectId();
			var $a = $(this).parent();
			ajaxPost("/notebook/addNotebook", {notebookId: notebookId, title: title}, function(ret) {
				if(ret.NotebookId) {
					Notebook.cache[ret.NotebookId] = ret;
					$a.attr("notebookId", notebookId);
					$a.html(title);
					// 选中之
					Notebook.ChangeNotebook(notebookId);
					
					// 改变nav
					Notebook.changeNav();
				}
			});
		}
	});
}

//-------------
// 删除
Notebook.DeleteNotebook = function(target) {
	var notebookId = $(target).attr("notebookId");
	if(!notebookId) {
		return;
	}
	
	ajaxGet("/notebook/deleteNotebook", {notebookId: notebookId}, function(ret) {
		if(ret.Ok) {
			$(target).parent().remove();
			delete Notebook.cache[notebookId];
			// 改变nav
			Notebook.changeNav();
		} else {
			alert(ret.Msg);
		}
	});
}

$(function() {
	//-------------------
	// 点击notebook
	$("#notebookList").on("click", "li a", function() {
		var notebookId = $(this).attr("notebookId");
		Notebook.ChangeNotebook(notebookId);
	});
	
	// 添加笔记本
	$("#addNotebookPlus").click(function(e) {
		e.stopPropagation();
		Notebook.AddNotebook();
	});
	
});

