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

Notebook.getCurNotebookId = function() {
	return Notebook.curNotebookId;
};

// 得到notebook标题, 给note显示其notebook标题用
// called by Note
Notebook.getNotebook = function(notebookId) {
	return Notebook.cache[notebookId];
}
// called by Note
Notebook.getNotebookTitle = function(notebookId) {
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
Notebook.curNotebookIsTrashOrAll = function() {
	return Notebook.curNotebookId == Notebook.trashNotebookId || Notebook.curNotebookId == Notebook.allNotebookId ;
}
Notebook.renderNotebooks = function(notebooks) {
	var self = this;

	if(!notebooks || typeof notebooks != "object" || notebooks.length < 0) {
		notebooks = [];
	}
	
	notebooks = [{NotebookId: Notebook.allNotebookId, Title: getMsg("all"), drop:false, drag: false}].concat(notebooks);
	notebooks.push({NotebookId: Notebook.trashNotebookId, Title: getMsg("trash"), drop:false, drag: false});
	Notebook.notebooks = notebooks; // 缓存之
	
	// 拖拽
	function beforeDrag(treeId, treeNodes) {
		for (var i=0,l=treeNodes.length; i<l; i++) {
			if (treeNodes[i].drag === false) {
				return false;
			}
		}
		return true;
	}
	function beforeDrop(treeId, treeNodes, targetNode, moveType) {
		return targetNode ? targetNode.drop !== false : true;
	}
	function onDrop(e, treeId, treeNodes, targetNode, moveType) {
		var treeNode = treeNodes[0];
		var parentNode;
		var treeObj = self.tree;
		var ajaxData = {curNotebookId: treeNode.NotebookId};
		
		// 成为子节点, 那么只需要得到targetNode下所有的子结点即可
		if(moveType == "inner") {
			parentNode = targetNode;
		// 在targetNode之前或之后, 
		// 那么: 1) 需要将该parentNode下所有的node重新排序即可; 2) treeNodes[0]为parentNode的子
		} else {
			parentNode = targetNode.getParentNode();
		}
		
		if(!parentNode) {
			var nodes = treeObj.getNodes();
		} else {
			ajaxData.parentNotebookId = parentNode.NotebookId;
			var nextLevel = parentNode.level+1;
			function filter(node) {
				return node.level == nextLevel;
			}
			var nodes = treeObj.getNodesByFilter(filter, false, parentNode);
		}
		
		ajaxData.siblings = [];
		for(var i in nodes) {
			var notebookId = nodes[i].NotebookId;
			if(!Notebook.isAllNotebookId(notebookId) && !Notebook.isTrashNotebookId(notebookId)) {
				ajaxData.siblings.push(notebookId);
			}
		}
		ajaxPost("/notebook/dragNotebooks", {data: JSON.stringify(ajaxData)});
	}
	// 添加自定义dom
	function addDiyDom(treeId, treeNode) {
		var spaceWidth = 5;
		var switchObj = $("#" + treeNode.tId + "_switch"),
		icoObj = $("#" + treeNode.tId + "_ico");
		switchObj.remove();
		icoObj.before(switchObj);
		if (treeNode.level > 1) {
			var spaceStr = "<span style='display: inline-block;width:" + (spaceWidth * treeNode.level)+ "px'></span>";
			switchObj.before(spaceStr);
		}
	}
	var setting = {
	
		view: {
			showLine: false,
			showIcon: false,
			selectedMulti: false,
			dblClickExpand: false,
			addDiyDom: addDiyDom
		},
		data: {
			key: {
				name: "Title",
				children: "Subs",
			}
		},
		edit: {
			enable: true,
			showRemoveBtn: false,
			showRenameBtn: false,
			drag: {
				isMove: true,
				prev: true,
				inner: true,
				next: true
			}
		},
		callback: {
			beforeDrag: beforeDrag,
			beforeDrop: beforeDrop,
			onDrop: onDrop,
			onClick: function(e, treeId, treeNode) {
				var notebookId = treeNode.NotebookId;
				Notebook.changeNotebook(notebookId);
			},
			onRightClick: function(event, treeId, treeNode) {
			},
			beforeRename: function(treeId, treeNode, newName, isCancel) {
				if(newName == "") {
					if(treeNode.IsNew) {
						// 删除之
						self.tree.removeNode(treeNode);
						return true;
					}
					return false;
				}
				if(treeNode.Title == newName) {
					return true;
				}
				
				// 如果是新添加的
				if(treeNode.IsNew) {
					var parentNode = treeNode.getParentNode();
					var parentNotebookId = parentNode ? parentNode.NotebookId : "";
					
					self.doAddNotebook(treeNode.NotebookId, newName, parentNotebookId);
				} else {
					self.doUpdateNotebookTitle(treeNode.NotebookId, newName);
				}
				return true;
			}
		}
	};
	
	self.tree = $.fn.zTree.init($("#notebookList"), setting, notebooks);
	
	// 展开/折叠图标
	var $notebookList = $("#notebookList");
	$notebookList.hover(function () {
		if (!$notebookList.hasClass("showIcon")) {
			$notebookList.addClass("showIcon");
		}
	}, function() {
		$notebookList.removeClass("showIcon");
	});
			
	// 缓存所有notebooks信息
	if(!isEmpty(notebooks)) {
		Notebook.curNotebookId = notebooks[0].NotebookId;
		self.cacheAllNotebooks(notebooks);
	}
	
	// 渲染nav
	Notebook.renderNav();
	
	// 渲染第一个notebook作为当前
	Notebook.changeNotebookNavForNewNote(notebooks[0].NotebookId);
}

Notebook.cacheAllNotebooks = function(notebooks) {
	var self = this;
	for(var i in notebooks) {
		var notebook = notebooks[i];
		Notebook.cache[notebook.NotebookId] = notebook;
		if(!isEmpty(notebook.Subs)) {
			self.cacheAllNotebooks(notebook.Subs);
		}
	}
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
		if(notebook) {
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
		}
	});
	
	$("#notebookNavForListNote").html(navForListNote);
	$("#notebookNavForNewNote").html(navForNewNote);
	$("#notebookNavForMoveNote").html(navForNewNote);
	
	// 移动, 复制重新来, 因为nav变了, 移动至-----的notebook导航也变了
	Note.initContextmenu();
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
Notebook.renderShareNotebooks = function(sharedUserInfos, shareNotebooks) {
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
Notebook.selectNotebook = function(target) {
	$("#notebook li a").removeClass("active");
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
	
	if(!Notebook.isAllNotebookId(notebookId) && !Notebook.isTrashNotebookId(notebookId)) {
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
	$("#newSharedNote").hide();
	
	// 搜索tag隐藏
	$("#tagSearch").hide();
}
Notebook.changeNotebookNav = function(notebookId) {
	Notebook.toggleToMyNav();
	
	// 1
	Notebook.selectNotebook($(t('#notebookList [notebookId="?"]', notebookId)));
	
	var notebook = Notebook.cache[notebookId];
	
	if(!notebook) {
		return;
	}
	
	// 2
	$("#curNotebookForListNote").html(notebook.Title);
	
	// 3
	Notebook.changeNotebookNavForNewNote(notebookId, notebook.Title);
}

Notebook.isAllNotebookId = function(notebookId) {
	return notebookId == Notebook.allNotebookId;
}
Notebook.isTrashNotebookId = function(notebookId) {
	return notebookId == Notebook.trashNotebookId;
}
// 当前选中的笔记本是否是"所有"
// called by Note
Notebook.curActiveNotebookIsAll = function() {
	return Notebook.isAllNotebookId($("#notebookList .active").attr("notebookId"));
}

// 改变笔记本
// 0. 改变样式
// 1. 改变note, 此时需要先保存
// 2. ajax得到该notebook下的所有note
// 3. 使用Note.RederNotes()
Notebook.changeNotebook = function(notebookId) {
	Notebook.changeNotebookNav(notebookId);
	
	Notebook.curNotebookId = notebookId;
		
	// 1
	Note.curChangedSaveIt();
	
	// 2 先清空所有
	Note.clearAll();
	
	var url = "/note/ListNotes/";
	var param = {notebookId: notebookId};
	
	// 废纸篓
	if(Notebook.isTrashNotebookId(notebookId)) {
		url = "/note/listTrashNotes";
		param = {};
	} else if(Notebook.isAllNotebookId(notebookId)) {
		param = {};
		// 得到全部的...
		cacheNotes = Note.getNotesByNotebookId();
		if(!isEmpty(cacheNotes)) { // 万一真的是没有呢?
			Note.renderNotesAndFirstOneContent(cacheNotes);
			return;
		}
	} else {
		cacheNotes = Note.getNotesByNotebookId(notebookId);
		if(!isEmpty(cacheNotes)) { // 万一真的是没有呢?
			Note.renderNotesAndFirstOneContent(cacheNotes);
			return;
		}
	}
	
	// 2 得到笔记本
	// 这里可以缓存起来, note按notebookId缓存
	ajaxGet(url, param, Note.renderNotesAndFirstOneContent);
}

// 是否是当前选中的notebookId
// 还包括共享
// called by Note
Notebook.isCurNotebook = function(notebookId) {
	return $(t('#notebookList [notebookId="?"], #shareNotebooks [notebookId="?"]', notebookId, notebookId)).attr("class") == "active";
}

// 改变nav, 为了新建note
// called by Note
Notebook.changeNotebookForNewNote = function(notebookId) {
	// 废纸篓
	if(Notebook.isTrashNotebookId(notebookId) || Notebook.isAllNotebookId(notebookId)) {
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
		Note.renderNotes(ret, true);
	});
};

//---------------------------
// 显示共享信息
Notebook.listNotebookShareUserInfo = function(target) {
	var notebookId = $(target).attr("notebookId");
	showDialogRemote("share/listNotebookShareUserInfo", {notebookId: notebookId});
}
// 共享笔记本
Notebook.shareNotebooks= function(target) {
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
Notebook.setNotebook2Blog = function(target) {
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
Notebook.updateNotebookTitle = function(target) {
	var self = Notebook;
	var notebookId = $(target).attr("notebookId");
	self.tree.editName(self.tree.getNodeByTId(notebookId));
	
	return;
	
	var notebookTitle = $(target).text();
	var id = "editNotebookTitle";
	$(target).html(t('<input type="text" value="?" everValue="?" id="?" notebookId="?"/>', notebookTitle, notebookTitle, id, $(target).attr("notebookId")));
	$("#" + id).focus();
}
Notebook.doUpdateNotebookTitle = function(notebookId, newTitle) {
	ajaxPost("/notebook/updateNotebookTitle", {notebookId: notebookId, title: newTitle}, function(ret) {
		// 修改缓存
		Notebook.cache[notebookId].Title = newTitle;
		// 改变nav
		Notebook.changeNav();
	});
}

//-----------
// 添加笔记本
// 1 确保是展开的
// 2 在所有后面添加<li></li>
Notebook.addNotebookSeq = 1; // inputId
Notebook.addNotebook = function() {
	var self = Notebook;
	if($("#myNotebooks").hasClass("closed")) {
		$("#myNotebooks .folderHeader").trigger("click");
	}
	
	// 添加并修改
	self.tree.addNodes(null, {Title: "", NotebookId: getObjectId(), IsNew: true}, true, true);
}

// rename 调用
Notebook.doAddNotebook = function(notebookId, title, parentNotebookId) {
	var self = Notebook;
	ajaxPost("/notebook/addNotebook", {notebookId: notebookId, title: title, parentNotebookId: parentNotebookId}, function(ret) {
		if(ret.NotebookId) {
			Notebook.cache[ret.NotebookId] = ret;
			var notebook = self.tree.getNodeByTId(notebookId);
			$.extend(notebook, ret);
			notebook.IsNew = false;
			
			// 选中之
			Notebook.changeNotebook(notebookId);
			
			// 改变nav
			Notebook.changeNav();
		}
	});
}

//-------------
// 添加子笔记本
Notebook.addChildNotebook = function(target) {
	var self = Notebook;
	if($("#myNotebooks").hasClass("closed")) {
		$("#myNotebooks .folderHeader").trigger("click");
	}
	
	var notebookId = $(target).attr("notebookId");
	
	// 添加并修改
	self.tree.addNodes(self.tree.getNodeByTId(notebookId), {Title: "", NotebookId: getObjectId(), IsNew: true}, false, true);
}

//-------------
// 删除
Notebook.deleteNotebook = function(target) {
	var self = Notebook;
	
	var notebookId = $(target).attr("notebookId");
	if(!notebookId) {
		return;
	}
	
	ajaxGet("/notebook/deleteNotebook", {notebookId: notebookId}, function(ret) {
		if(ret.Ok) {
			/*
			$(target).parent().remove();
			*/
			self.tree.removeNode(self.tree.getNodeByTId(notebookId));
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
	/*
	$("#myNotebooks").on("click", "ul.folderBody li a", function() {
		var notebookId = $(this).attr("notebookId");
		Notebook.changeNotebook(notebookId);
	});
	*/
	// min
	$("#minNotebookList").on("click", "li", function() {
		var notebookId = $(this).find("a").attr("notebookId");
		Notebook.changeNotebook(notebookId);
	});
	
	// 修改笔记本标题, blur后修改标题之
	/*
	enterBlur("#notebookList", "input#editNotebookTitle");
	$("#notebookList").on("blur", "input#editNotebookTitle", Notebook.doUpdateNotebookTitle);
	*/
	
	//-------------------
	// 右键菜单
	var notebookListMenu = {
		width: 150, 
		items: [
			{ text: "分享给好友", alias: 'shareToFriends', icon: "", faIcon: "fa-share-square-o", action: Notebook.listNotebookShareUserInfo},
			{ type: "splitLine" },
			{ text: "公开为博客", alias: 'set2Blog', icon: "", action: Notebook.setNotebook2Blog },
			{ text: "取消公开为博客", alias: 'unset2Blog', icon: "", action: Notebook.setNotebook2Blog }, // Unset
			{ type: "splitLine" },
			{ text: "添加子笔记本", icon: "", action: Notebook.addChildNotebook },
			{ text: "重命名", icon: "", action: Notebook.updateNotebookTitle },
			{ text: "删除", icon: "", alias: 'delete', faIcon: "fa-trash-o", action: Notebook.deleteNotebook }
		],
		onShow: applyrule,
    	onContextMenu: beforeContextMenu,
    	parent: "#notebookList ",
    	children: "li a"
	}
	
	function applyrule(menu) {
		var notebookId = $(this).attr("notebookId");
		var notebook = Notebook.cache[notebookId];
		if(!notebook) {
			return;
		}
		// disabled的items
		var items = [];
		// 是否已公开为blog
		if(!notebook.IsBlog) {
			items.push("unset2Blog");
		} else {
			items.push("set2Blog");
		}
		// 是否还有笔记
		if(Note.notebookHasNotes(notebookId)) {
			items.push("delete");
		}
        menu.applyrule({
        	name: "target2",
            disable: true,
            items: items
        });
	}
	// 哪个不能
	function beforeContextMenu() {
		var notebookId = $(this).attr("notebookId");
		return !Notebook.isTrashNotebookId(notebookId) && !Notebook.isAllNotebookId(notebookId);
	}
	$("#notebookList li a").contextmenu(notebookListMenu);
	
	//------------------
	// 添加notebook
	// 右键菜单
	var addNotebookContextmenu = {
		width: 150, 
		items: [
			{ text: "添加笔记本", icon: "", action: Notebook.addNotebook },
		],
    	parent: "#myNotebooks",
    	children: ""
	}
	$("#myNotebooks").contextmenu(addNotebookContextmenu);
	
	//---------------
	// 点击中部notebook nav
	$("#notebookNavForListNote").on("click", "li", function() {
		var notebookId = $(this).find("a").attr("notebookId");
		Notebook.changeNotebook(notebookId);
	});
	
	// 添加笔记本
	$("#addNotebookPlus").click(function(e) {
		e.stopPropagation();
		Notebook.addNotebook();
	});
});