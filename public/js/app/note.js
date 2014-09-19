// 1. notebook change
// notebook一改变, 当前的肯定要保存, ajax是异步的. 此时先清空所有note信息. -> 得到该notebook的notes, 显示出来, 并选中第一个!
// 在这期间定时器还会保存, curNoteId还没换, 所以会清空curNoteId的content!!!

// 2. note change, save cur, 立即curNoteId = ""!!

// 3. 什么时候设置curNoteId? 是ajax得到内容之后设置

// note
Note.curNoteId = "";

Note.interval = ""; // 定时器

Note.itemIsBlog = '<div class="item-blog"><i class="fa fa-bold" title="blog"></i></div><div class="item-setting"><i class="fa fa-cog" title="setting"></i></div>';
// for render
Note.itemTplNoImg = '<div href="#" class="item ?" noteId="?">'
Note.itemTplNoImg += Note.itemIsBlog +'<div class="item-desc" style="right: 0;"><p class="item-title">?</p><p class="item-text"><i class="fa fa-book"></i> <span class="note-notebook">?</span> <i class="fa fa-calendar"></i> <span class="updated-time">?</span> <br /><span class="desc">?</span></p></div></div>';

Note.itemTpl = '<div href="#" class="item ?" noteId="?"><div class="item-thumb" style=""><img src="?"/></div>'
Note.itemTpl +=Note.itemIsBlog + '<div class="item-desc" style=""><p class="item-title">?</p><p class="item-text"><i class="fa fa-book"></i> <span class="note-notebook">?</span> <i class="fa fa-calendar"></i> <span class="updated-time">?</span> <br /><span class="desc">?</span></p></div></div>';

// for new
Note.newItemTpl = '<div href="#" class="item item-active ?" fromUserId="?" noteId="?">'
Note.newItemTpl += Note.itemIsBlog + '<div class="item-desc" style="right: 0px;"><p class="item-title">?</p><p class="item-text"><i class="fa fa-book"></i> <span class="note-notebook">?</span> <i class="fa fa-calendar"></i> <span class="updated-time">?</span><br /><span class="desc">?</span></p></div></div>';

Note.noteItemListO = $("#noteItemList");

// notbeookId => {"updatedTime" => [noteId1, noteId2], "title" => [noteId1, noteId2...]} 排序方式分组
// 一旦某notebook改变了就清空, 重新排序之. (用js排)
Note.cacheByNotebookId = {all: {}};
Note.notebookIds = {}; // notebookId => true

Note.isReadOnly = false;
// 定时保存信息
Note.intervalTime = 600000; // 600s, 10mins
Note.startInterval = function() {
	Note.interval = setInterval(function() {
		log("自动保存开始...")
		changedNote = Note.curChangedSaveIt(false);
	}, Note.intervalTime); // 600s, 10mins
}
// 停止, 当切换note时
// 但过5000后自动启动
Note.stopInterval = function() {
	clearInterval(Note.interval);
	
	setTimeout(function() {
		Note.startInterval();
	}, Note.intervalTime);
}

// note = {NoteId, Desc, UserId,...}
Note.addNoteCache = function(note) {
	Note.cache[note.NoteId] = note;
	Note.clearCacheByNotebookId(note.NotebookId);
}
// content = {NoteId:, Content:}
// 还可以设置其它的值
Note.setNoteCache = function(content, clear) {
	if(!Note.cache[content.NoteId]) {
		 Note.cache[content.NoteId] = content;
	} else {
		$.extend(Note.cache[content.NoteId], content);
	}
	
	if(clear == undefined) {
		clear = true;
	}
	if(clear) {
		Note.clearCacheByNotebookId(content.NotebookId);
	}
}

// 得到当前的笔记
Note.getCurNote = function() {
	var self = this;
	if(self.curNoteId == "") {
		return null;
	}
	return self.cache[self.curNoteId];
}

// 每当有notebookId相应的note改变时都要重新清空之
// 并设置该notebookId有值
Note.clearCacheByNotebookId = function(notebookId) {
	if(notebookId) {
		Note.cacheByNotebookId[notebookId] = {};
		Note.cacheByNotebookId["all"] = {};
		Note.notebookIds[notebookId] = true;
	}
}

// notebook是否有notes
// called by Notebook
Note.notebookHasNotes = function(notebookId) {
	var notes = Note.getNotesByNotebookId(notebookId);
	return !isEmpty(notes);
}

// 得到notebook下的notes, 按什么排序 updatedTime?
Note.getNotesByNotebookId = function(notebookId, sortBy, isAsc) {
	if(!sortBy) {
		sortBy = "UpdatedTime";
	}
	if(isAsc == "undefined") {
		isAsc = false; // 默认是降序
	}
	
	if(!notebookId) {
		notebookId = "all";
	}
	
	if(!Note.cacheByNotebookId[notebookId]) {
		return [];
	}
	
	if(Note.cacheByNotebookId[notebookId][sortBy]) {
		return Note.cacheByNotebookId[notebookId][sortBy];
	} else {
	}
	
	// 从所有的notes中找到notebookId的, 并排序之
	var notes = [];
	var sortBys = [];
	for(var i in Note.cache) {
		if(!i) {
			continue;
		}
		var note = Note.cache[i];
		// 不要trash的not, 共享的也不要
		if(note.IsTrash || note.IsShared) {
			continue;
		}
		if(notebookId == "all" || note.NotebookId == notebookId) {
			notes.push(note);
		}
	}
	// 排序之
	notes.sort(function(a, b) {
		var t1 = a[sortBy];
		var t2 = b[sortBy];
		
		if(isAsc) {
			if(t1 < t2) {
				return -1;
			} else if (t1 > t2) {
				return 1;
			}	
		} else {
			if(t1 < t2) {
				return 1;
			} else if (t1 > t2) {
				return -1;
			}
		}
		return 0;
	});
	
	// 缓存之
	Note.cacheByNotebookId[notebookId][sortBy] = notes;
	return notes;
}

// called by Notebook
// render 所有notes, 和第一个note的content
Note.renderNotesAndFirstOneContent = function(ret) {
	// 错误的ret是一个Object
	if(!isArray(ret)) {
		return;
	}
	
	// note 导航
	Note.renderNotes(ret);
	// 渲染第一个
	if(!isEmpty(ret[0])) {
		Note.changeNote(ret[0].NoteId);
	} else {
	}

}

// 当前的note是否改变过了?
// 返回已改变的信息
// force bool true表示content比较是比较HTML, 否则比较text, 默认为true
// 定时保存用false
Note.curHasChanged = function(force) {
	if(force == undefined) {
		force = true;
	}
	var cacheNote = Note.cache[Note.curNoteId] || {};
	// 收集当前信息, 与cache比对
	var title = $("#noteTitle").val();
	var tags = Tag.getTags(); // TODO
	
	// 如果是markdown返回[content, preview]
	var contents = getEditorContent(cacheNote.IsMarkdown);
	var content, preview;
	var contentText;
	if (isArray(contents)) {
		content = contents[0];
		preview = contents[1];
		contentText = content;
		// preview可能没来得到及解析
		if (content && previewIsEmpty(preview)) {
			preview = Converter.makeHtml(content);
		}
		if(!content) {
			preview = "";
		}
		cacheNote.Preview = preview; // 仅仅缓存在前台
	} else {
		content = contents;
		try {
			contentText = $(content).text();
		} catch(e) {
		}
	}
	
	var hasChanged = {
		hasChanged: false, // 总的是否有改变
		IsNew: cacheNote.IsNew, // 是否是新添加的
		IsMarkdown: cacheNote.IsMarkdown, // 是否是markdown笔记
		FromUserId: cacheNote.FromUserId, // 是否是共享新建的
		NoteId: cacheNote.NoteId,
		NotebookId: cacheNote.NotebookId
	};
	
	if(hasChanged.IsNew) {
		$.extend(hasChanged, cacheNote);
	}
	
	if(cacheNote.Title != title) {
		hasChanged.hasChanged = true; // 本页使用用小写
		hasChanged.Title = title; // 要传到后台的用大写
		if(!hasChanged.Title) {
//			alert(1);
		}
	}
	
	if(!arrayEqual(cacheNote.Tags, tags)) {
		hasChanged.hasChanged = true;
		hasChanged.Tags = tags;
	}
	
	// 比较text, 因为note Nav会添加dom会导致content改变
	if((force && cacheNote.Content != content) || (!force && $(cacheNote.Content).text() != contentText)) {
		hasChanged.hasChanged = true;
		hasChanged.Content = content;
		
		// 从html中得到...
		var c = preview || content;
		
		hasChanged.Desc = Note.genDesc(c);
		hasChanged.ImgSrc = Note.getImgSrc(c);
		hasChanged.Abstract = Note.genAbstract(c);
		
	} else {
		log("text相同");
		log(cacheNote.Content == content);
	}
	
	hasChanged["UserId"] = cacheNote["UserId"] || "";
	
	return hasChanged;
}

// 由content生成desc
// 换行不要替换
Note.genDesc = function(content) {
	if(!content) {
		return "";
	}
	
	// 将</div>, </p>替换成\n
	var token = "ALEALE";
	content = content.replace(/<\/p>/g, token); 
	content = content.replace(/<\/div>/g, token);
	content = content.replace(/<\/?.+?>/g," ");
	
	pattern = new RegExp(token, "g");
	content = content.replace(pattern, "<br />");
	content = content.replace(/<br \/>( *)<br \/>/g, "<br />"); // 两个<br />之间可能有空白
	content = content.replace(/<br \/>( *)<br \/>/g, "<br />");
	
	// 去掉最开始的<br />或<p />
	content = trimLeft(content, " ");
	content = trimLeft(content, "<br />");
	content = trimLeft(content, "</p>");
	content = trimLeft(content, "</div>");
	
	if(content.length < 300) {
		return content;
	}
	return content.substring(0, 300);
}

// 得到摘要
Note.genAbstract = function(content, len) {
	if(len == undefined) {
		len = 1000;
	}
	if(content.length < len) {
		return content;
	}
	var isCode = false;
	var isHTML = false;
	var n = 0;
	var result = "";
	var maxLen = len;
	for(var i = 0; i < content.length; ++i) {
		var temp = content[i]
		if (temp == '<') {
			isCode = true
		} else if (temp == '&') {
			isHTML = true
		} else if (temp == '>' && isCode) {
			n = n - 1
			isCode = false
		} else if (temp == ';' && isHTML) {
			isHTML = false
		}
		if (!isCode && !isHTML) {
			n = n + 1
		}
		result += temp
		if (n >= maxLen) {
			break
		}
	}
	
	var d = document.createElement("div");
    d.innerHTML = result
    return d.innerHTML;
}

Note.getImgSrc = function(content) {
	if(!content) {
		return "";
	}
	var imgs = $(content).find("img");
	for(var i in imgs) {
		var src = imgs.eq(i).attr("src");
		if(src) {
			return src;
		}
	}
	return "";
}

// 如果当前的改变了, 就保存它
// 以后要定时调用
// force , 默认是true, 表强校验内容
// 定时保存传false
Note.curChangedSaveIt = function(force) {
	// 如果当前没有笔记, 不保存
	if(!Note.curNoteId || Note.isReadOnly) {
		return;
	}
	
	var hasChanged = Note.curHasChanged(force);
		
	// 把已改变的渲染到左边 item-list
	Note.renderChangedNote(hasChanged);
	
	if(hasChanged.hasChanged || hasChanged.IsNew) {
		delete hasChanged.hasChanged;
		
		// 先缓存, 把markdown的preview也缓存起来
		Note.setNoteCache(hasChanged, false);
		
		// 设置更新时间
		Note.setNoteCache({"NoteId": hasChanged.NoteId, "UpdatedTime": (new Date()).format("yyyy-MM-ddThh:mm:ss.S")}, false);
		
		// 保存之
		showMsg(getMsg("saving"));
		ajaxPost("/note/UpdateNoteOrContent", hasChanged, function(ret) {
			if(hasChanged.IsNew) {
				// 缓存之, 后台得到其它信息
				ret.IsNew = false;
				Note.setNoteCache(ret, false);
			}
			showMsg(getMsg("saveSuccess"), 1000);
		});
		
		return hasChanged;
	}
	return false;
}

// 样式
Note.selectTarget = function(target) {
	$(".item").removeClass("item-active");
	$(target).addClass("item-active");
}

// 改变note
// 可能改变的是share note
// 1. 保存之前的note
// 2. ajax得到现在的note
Note.changeNote = function(selectNoteId, isShare, needSaveChanged) {
	// -1 停止定时器
	Note.stopInterval();
	
	// 0
	var target = $(t('[noteId="?"]', selectNoteId))
	Note.selectTarget(target);
	
	// 1 之前的note, 判断是否已改变, 改变了就要保存之
	// 这里, 在搜索的时候总是保存, 搜索的话, 比较快, 肯定没有变化, 就不要执行该操作
	if(needSaveChanged == undefined) {
		needSaveChanged  = true;
	}
	if(needSaveChanged) {
		var changedNote = Note.curChangedSaveIt();
	}
	
	// 2. 设空, 防止在内容得到之前又发生保存
	Note.curNoteId = "";
	
	// 2 得到现在的
	// ajax之
	var cacheNote = Note.cache[selectNoteId];
	
	// 判断是否是共享notes
	if(!isShare) {
		if(cacheNote.Perm != undefined) {
			isShare = true;
		}
	}
	var hasPerm = !isShare || Share.hasUpdatePerm(selectNoteId); // 不是共享, 或者是共享但有权限
	
	// 不是手机浏览器且有权限
	if(!LEA.isMobile && hasPerm) {
		Note.hideReadOnly();
		Note.renderNote(cacheNote);
		
		// 这里要切换编辑器
		switchEditor(cacheNote.IsMarkdown)
		
	} else {
		Note.renderNoteReadOnly(cacheNote);
	}
	
	function setContent(ret) {
		Note.setNoteCache(ret, false);
		// 把其它信息也带上
		ret = Note.cache[selectNoteId]
		if(!LEA.isMobile && hasPerm) {
			Note.renderNoteContent(ret);
		} else {
			Note.renderNoteContentReadOnly(ret);
		}
		hideLoading();
	}
	
	if(cacheNote.Content) {
		setContent(cacheNote);
		return;
	}
	
	var url = "/note/GetNoteContent";
	var param = {noteId: selectNoteId};
	if(isShare) {
		url = "/share/GetShareNoteContent";
		param.sharedUserId = cacheNote.UserId // 谁的笔记
	}
	
	// 这里loading
	showLoading();
	ajaxGet(url, param, setContent);
}

// 渲染

// 更改信息到左侧
// 定时更改 当前正在编辑的信息到左侧导航
// 或change select. 之前的note, 已经改变了
Note.renderChangedNote = function(changedNote) {
	if(!changedNote) {
		return;
	}
	
	// 找到左侧相应的note
	var $leftNoteNav = $(t('[noteId="?"]', changedNote.NoteId));
	if(changedNote.Title) {
		$leftNoteNav.find(".item-title").html(changedNote.Title);
	}
	if(changedNote.Desc) {
		$leftNoteNav.find(".desc").html(changedNote.Desc);
	}
	if(changedNote.ImgSrc && !LEA.isMobile) {
		$thumb = $leftNoteNav.find(".item-thumb");
		// 有可能之前没有图片
		if($thumb.length > 0) {
			$thumb.find("img").attr("src", changedNote.ImgSrc);
		} else {
			$leftNoteNav.append(t('<div class="item-thumb" style=""><img src="?"></div>', changedNote.ImgSrc));
		}
		$leftNoteNav.find(".item-desc").removeAttr("style");
	} else if(changedNote.ImgSrc == "") {
		$leftNoteNav.find(".item-thumb").remove(); // 以前有, 现在没有了
		$leftNoteNav.find(".item-desc").css("right", 0);
	}
}

// 清空右侧note信息, 可能是共享的, 
// 此时需要清空只读的, 且切换到note edit模式下
Note.clearNoteInfo = function() {
	Note.curNoteId = "";
	Tag.clearTags();
	$("#noteTitle").val("");
	setEditorContent("");
	
	// markdown editor
	$("#wmd-input").val("");
	$("#wmd-preview").html("");
	
	// 只隐藏即可
	$("#noteRead").hide();
}
// 清除noteList导航
Note.clearNoteList = function() {
	Note.noteItemListO.html(""); // 清空
}

// 清空所有, 在转换notebook时使用
Note.clearAll = function() {
	// 当前的笔记清空掉
	Note.curNoteId = "";
	
	Note.clearNoteInfo();
	Note.clearNoteList();
}

// render到编辑器
// render note
Note.renderNote = function(note) {
	if(!note) {
		return;
	}
	// title
	$("#noteTitle").val(note.Title);
	
	// 当前正在编辑的
	// tags
	Tag.renderTags(note.Tags);
}

// render content
Note.renderNoteContent = function(content) {
	setEditorContent(content.Content, content.IsMarkdown, content.Preview);
	// 只有在renderNoteContent时才设置curNoteId
	Note.curNoteId = content.NoteId;
}

// 初始化时渲染最初的notes
/**
    <div id="noteItemList">
	  <!--
      <div href="#" class="item">
        <div class="item-thumb" style="">
          <img src="images/a.gif"/>
        </div>

        <div class="item-desc" style="">
            <p class="item-title">?</p>
            <p class="item-text">
            	?
            </p>
        </div>
      </div>
      -->
*/

Note.showEditorMask = function() {
	$("#editorMask").css("z-index", 10);
	// 要判断是否是垃圾筒
	if(Notebook.curNotebookIsTrashOrAll()) {
		$("#editorMaskBtns").hide();
		$("#editorMaskBtnsEmpty").show();
	} else {
		$("#editorMaskBtns").show();
		$("#editorMaskBtnsEmpty").hide();
	}
}
Note.hideEditorMask = function() {
	$("#editorMask").css("z-index", -10);
}

// 这里如果notes过多>100个将会很慢!!, 使用setTimeout来分解
Note.renderNotesC = 0;
Note.renderNotes = function(notes, forNewNote, isShared) {
	var renderNotesC = ++Note.renderNotesC;
	
	$("#noteItemList").slimScroll({ scrollTo: '0px', height: "100%", onlyScrollBar: true});
	
	if(!notes || typeof notes != "object" || notes.length <= 0) {
		// 如果没有, 那么是不是应该hide editor?
		if(!forNewNote) {
			Note.showEditorMask();
		}
		return;
	}
	Note.hideEditorMask();
	// 新建笔记时会先创建一个新笔记, 所以不能清空
	if(forNewNote == undefined) {
		forNewNote = false;
	}
	if(!forNewNote) {
		Note.noteItemListO.html(""); // 清空
	}
	
	// 20个一次
	var len = notes.length;
	var c = Math.ceil(len/20);
	
	Note._renderNotes(notes, forNewNote, isShared, 1);
	
	// 先设置缓存
	for(var i = 0; i < len; ++i) {
		var note = notes[i];
		// 不清空
		// 之前是addNoteCache, 如果是搜索出的, 会把内容都重置了
		Note.setNoteCache(note, false);
		
		// 如果是共享的笔记本, 缓存也放在Share下
		if(isShared) {
			Share.setCache(note);
		}
	}
	
	for(var i = 1; i < c; ++i) {
		setTimeout(
			(function(i) {
				// 防止还没渲染完就点击另一个notebook了
				return function() {
					if(renderNotesC == Note.renderNotesC) {
						Note._renderNotes(notes, forNewNote, isShared, i+1);
					}
				}
			})(i), i*2000);
	}
}
Note._renderNotes = function(notes, forNewNote, isShared, tang) { // 第几趟
	var baseClasses = "item-my";
	if(isShared) {
		baseClasses = "item-shared";
	}
	
	var len = notes.length;
	for(var i = (tang-1)*20; i < len && i < tang*20; ++i) {
		var classes = baseClasses;
		if(!forNewNote && i == 0) {
			classes += " item-active";
		}
		var note = notes[i];
		var tmp;
		if(note.ImgSrc && !LEA.isMobile) {
			tmp = t(Note.itemTpl, classes, note.NoteId, note.ImgSrc, note.Title, Notebook.getNotebookTitle(note.NotebookId), goNowToDatetime(note.UpdatedTime), note.Desc);
		} else {
			tmp = t(Note.itemTplNoImg, classes, note.NoteId, note.Title, Notebook.getNotebookTitle(note.NotebookId), goNowToDatetime(note.UpdatedTime), note.Desc);
		}
		if(!note.IsBlog) {
			tmp = $(tmp);
			tmp.find(".item-blog").hide();
		}
		Note.noteItemListO.append(tmp);
		
		/*
		// 共享的note也放在Note的cache一份
		if(isShared) {
			note.IsShared = true; // 注明是共享的
		}
		
		// 不清空
		// 之前是addNoteCache, 如果是搜索出的, 会把内容都重置了
		Note.setNoteCache(note, false);
		
		// 如果是共享的笔记本, 缓存也放在Share下
		if(isShared) {
			Share.setCache(note);
		}
		*/
	}
} 

// 新建一个笔记
// 要切换到当前的notebook下去新建笔记
// isShare时fromUserId才有用
// 3.8 add isMarkdown
Note.newNote = function(notebookId, isShare, fromUserId, isMarkdown) {
	// 切换编辑器
	switchEditor(isMarkdown);
	Note.hideEditorMask();
	
	// 防止从共享read only跳到添加
	Note.hideReadOnly();
	
	Note.stopInterval();
	// 保存当前的笔记
	Note.curChangedSaveIt();
	
	var note = {NoteId: getObjectId(), Title: "", Tags:[], Content:"", NotebookId: notebookId, IsNew: true, FromUserId: fromUserId, IsMarkdown: isMarkdown}; // 是新的
	// 添加到缓存中
	Note.addNoteCache(note);
	
	// 是否是为共享的notebook添加笔记, 如果是, 则还要记录fromUserId
	var newItem = "";
	
	var baseClasses = "item-my";
	if(isShare) {
		baseClasses = "item-shared";
	}
	
	var notebook = Notebook.getNotebook(notebookId);
	var notebookTitle = notebook ? notebook.Title : "";
	var curDate = getCurDate();
	if(isShare) {
		newItem = t(Note.newItemTpl, baseClasses, fromUserId, note.NoteId, note.Title, notebookTitle, curDate, "");
	} else {
		newItem = t(Note.newItemTpl, baseClasses, "", note.NoteId, note.Title, notebookTitle, curDate, "");
	}
	
	// notebook是否是Blog
	if(!notebook.IsBlog) {
		newItem = $(newItem);
		newItem.find(".item-blog").hide();
	}
	
	// 是否在当前notebook下, 不是则切换过去, 并得到该notebook下所有的notes, 追加到后面!
	if(!Notebook.isCurNotebook(notebookId)) {
		// 先清空所有
		Note.clearAll();
		
		// 插入到第一个位置
		Note.noteItemListO.prepend(newItem);
		
		// 改变为当前的notebookId
		// 会得到该notebookId的其它笔记
		if(!isShare) {
			Notebook.changeNotebookForNewNote(notebookId);
		} else {
			Share.changeNotebookForNewNote(notebookId);
		}
	} else {
		// 插入到第一个位置
		Note.noteItemListO.prepend(newItem);
	}
	
	Note.selectTarget($(t('[noteId="?"]', note.NoteId)));
	
	$("#noteTitle").focus();
	
	Note.renderNote(note);
	Note.renderNoteContent(note);
	Note.curNoteId = note.NoteId;
}

// 保存note ctrl + s
Note.saveNote = function(e) {
	var num = e.which ? e.which : e.keyCode;
	// 保存
    if((e.ctrlKey || e.metaKey) && num == 83 ) { // ctrl + s or command + s
    	Note.curChangedSaveIt();
    	e.preventDefault();
    	return false;
    } else {
    }
};

// 删除或移动笔记后, 渲染下一个或上一个
Note.changeToNext = function(target) {
	var $target = $(target);
	var next = $target.next();
	if(!next.length) {
		var prev = $target.prev();
		if(prev.length) {
			next = prev;
		} else {
			// 就它一个
			Note.showEditorMask();
			return;
		}
	}
	
	Note.changeNote(next.attr("noteId"));
}

// 删除笔记
// 1. 先隐藏, 成功后再移除DOM
// 2. ajax之 noteId
// Share.deleteSharedNote调用
Note.deleteNote = function(target, contextmenuItem, isShared) {
	// 如果删除的是已选中的, 赶紧设置curNoteId = null
	if($(target).hasClass("item-active")) {
		// -1 停止定时器
		Note.stopInterval();
		// 不保存
		Note.curNoteId = null;
		// 清空信息
		Note.clearNoteInfo();
	}
	
	noteId = $(target).attr("noteId");
	if(!noteId) {
		return;
	}
	// 1
	$(target).hide();
	
	// 2
	var note = Note.cache[noteId];
	var url = "/note/deleteNote"
	if(note.IsTrash) {
		url = "/note/deleteTrash";
	}
	
	ajaxGet(url, {noteId: noteId, userId: note.UserId, isShared: isShared}, function(ret) {
		if(ret) {
			Note.changeToNext(target);
			
			$(target).remove();
			
			// 删除缓存
			if(note) {
				Note.clearCacheByNotebookId(note.NotebookId)
				delete Note.cache[noteId]
			}
			
			showMsg("删除成功!", 500);
		} else {
			// 弹出信息 popup 不用点确认的
			$(target).show();
			showMsg("删除失败!", 2000);
		}
	});
}

// 显示共享信息
Note.listNoteShareUserInfo = function(target) {
	var noteId = $(target).attr("noteId");
	showDialogRemote("share/listNoteShareUserInfo", {noteId: noteId});
}
	
// 共享笔记
Note.shareNote = function(target) {
	var title = $(target).find(".item-title").text();
	showDialog("dialogShareNote", {title: "分享笔记给好友-" + title});
	
	setTimeout(function() {
		$("#friendsEmail").focus();
	}, 500);
	
	var noteId = $(target).attr("noteId");
	shareNoteOrNotebook(noteId, true);
}

// 历史记录
Note.listNoteContentHistories = function() {
	// 弹框
	$("#leanoteDialog #modalTitle").html(getMsg("history"));
	$content = $("#leanoteDialog .modal-body");
	$content.html("");
	$("#leanoteDialog .modal-footer").html('<button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>');
	options = {}
	options.show = true;
	$("#leanoteDialog").modal(options);
	
	ajaxGet("noteContentHistory/listHistories", {noteId: Note.curNoteId}, function(re) {
		if(!isArray(re)) {$content.html("无历史记录"); return}
		// 组装成一个tab
		var str = 'leanote会保存笔记的最近10份历史记录. <div id="historyList"><table class="table table-hover">';
		note = Note.cache[Note.curNoteId];
		var s = "div"
		if(note.IsMarkdown) {
			s = "pre";
		}
		for (i in re) {
			var content = re[i]
			content.Ab = Note.genAbstract(content.Content, 200);
			str += t('<tr><td seq="?"><? class="each-content">?</?> <div class="btns">时间: <span class="label label-default">?</span> <button class="btn btn-default all">展开</button> <button class="btn btn-primary back">还原</button></div></td></tr>', i, s, content.Ab, s, goNowToDatetime(content.UpdatedTime))
		}
		str += "</table></div>";
		$content.html(str);
		$("#historyList .all").click(function() {
			$p = $(this).parent().parent();
			var seq = $p.attr("seq");
			var $c = $p.find(".each-content");
			if($(this).text() == "展开") {
				$(this).text("折叠")
				$c.html(re[seq].Content);
			} else {
				$(this).text("展开")
				$c.html(re[seq].Ab);
			}
		});
		
		// 还原
		$("#historyList .back").click(function() {
			$p = $(this).parent().parent();
			var seq = $p.attr("seq");
			if(confirm("确定要从该版还原? 还原前leanote会备份当前版本到历史记录中.")) {
				// 保存当前版本
				Note.curChangedSaveIt();
				// 设置之
				note = Note.cache[Note.curNoteId];
				setEditorContent(re[seq].Content, note.IsMarkdown);
				//
				hideDialog();
			}
		});
		
	});
}

// 长微博
Note.html2Image = function(target) {
	var noteId = $(target).attr("noteId");
	showDialog("html2ImageDialog", {title: "发送长微博", postShow: function() {
		ajaxGet("/note/html2Image", {noteId: noteId}, function(ret) {
			if (typeof ret == "object" && ret.Ok) {
				$("#leanoteDialog .weibo span").html("生成成功, 右键图片保存到本地.")
				$("#leanoteDialog .weibo img").attr("src", ret.Id);
				$("#leanoteDialog .sendWeiboBtn").removeClass("disabled");
				$("#leanoteDialog .sendWeiboBtn").click(function() {
					var title = Note.cache[noteId].Title;
					var url = "http://service.weibo.com/share/share.php?title=" + title + " (" + UserInfo.Username + "分享. 来自leanote.com)";
					url += "&pic=" + UrlPrefix + ret.Id;
					window.open(url, "_blank");
				});
			} else {
				$("#leanoteDialog .weibo span").html("对不起, 我们出错了!")
			}
		});
	}});
}

//--------------
// read only

Note.showReadOnly = function() {
	Note.isReadOnly = true;
	$("#noteRead").show();
}
Note.hideReadOnly = function() {
	Note.isReadOnly = false;
	$("#noteRead").hide();
}
// read only
Note.renderNoteReadOnly = function(note) {
	Note.showReadOnly();
	$("#noteReadTitle").html(note.Title);
	
	Tag.renderReadOnlyTags(note.Tags);
	
	$("#noteReadCreatedTime").html(goNowToDatetime(note.CreatedTime));
	$("#noteReadUpdatedTime").html(goNowToDatetime(note.UpdatedTime));
}
Note.renderNoteContentReadOnly = function(note) {
	if(note.IsMarkdown) {
		$("#noteReadContent").html('<pre id="readOnlyMarkdown">' + note.Content + "</pre>");
	} else {
		$("#noteReadContent").html(note.Content);
	}
}

//---------------------------
// 搜索
// 有点小复杂, 因为速度过快会导致没加载完, 然后就保存上一个 => 致使标题没有
// 为什么会标题没有?
Note.lastSearch = null;
Note.lastKey = null; // 判断是否与上一个相等, 相等就不查询, 如果是等了很久再按enter?
Note.lastSearchTime = new Date();
Note.isOver2Seconds = false;
Note.isSameSearch = function(key) {
	// 判断时间是否超过了1秒, 超过了就认为是不同的
	var now = new Date();
	var duration = now.getTime() - Note.lastSearchTime.getTime();
	Note.isOver2Seconds = duration > 2000 ? true : false;
	if(!Note.lastKey || Note.lastKey != key || duration > 1000) {
		Note.lastKey = key;
		Note.lastSearchTime = now;
		return false;
	}
	
	if(key == Note.lastKey) {
		return true;
	}
	
	Note.lastSearchTime = now;
	Note.lastKey = key;
	return false;
}

Note.searchNote = function() {
	var val = $("#searchNoteInput").val();
	if(!val) {
		// 定位到all
		Notebook.changeNotebook("0");
		return;
	}
	// 判断是否与上一个是相同的搜索, 是则不搜索
	if(Note.isSameSearch(val)) {
		return;
	}
	
	// 之前有, 还有结束的
	if(Note.lastSearch) {
		Note.lastSearch.abort();
	}
	
	// 步骤与tag的搜索一样 
	// 1
	Note.curChangedSaveIt();
	
	// 2 先清空所有
	Note.clearAll();
	
	// 发送请求之
	// 先取消上一个
	showLoading();
	Note.lastSearch = $.post("/note/searchNote", {key: val}, function(notes) {
		hideLoading();
		if(notes) {
			// 成功后设为空
			Note.lastSearch = null;
			// renderNotes只是note列表加载, 右侧笔记详情还没加载
			// 这个时候, 定位第一个, 保存之前的,
			// 	如果: 第一次搜索, renderNotes OK, 还没等到changeNote时
			//		第二次搜索来到, Note.curChangedSaveIt();
			//		导致没有标题了
			// 不是这个原因, 下面的Note.changeNote会导致保存
			
			// 设空, 防止发生上述情况
			// Note.curNoteId = "";
			
			Note.renderNotes(notes);
			if(!isEmpty(notes)) {
				Note.changeNote(notes[0].NoteId, false/*, true || Note.isOver2Seconds*/); // isShare, needSaveChanged?, 超过2秒就要保存
			}
		} else {
			// abort的
		}
	});
	// Note.lastSearch.abort();
}

//----------
//设为blog/unset
Note.setNote2Blog = function(target) {
	var noteId = $(target).attr("noteId");
	var note = Note.cache[noteId];
	var isBlog = true;
	if(note.IsBlog != undefined) {
		isBlog = !note.IsBlog;
	}
	// 标志添加/去掉
	if(isBlog) {
		$(target).find(".item-blog").show();
	} else {
		$(target).find(".item-blog").hide();
	}
	ajaxPost("/blog/setNote2Blog", {noteId: noteId, isBlog: isBlog}, function(ret) {
		if(ret) {
			Note.setNoteCache({NoteId: noteId, IsBlog: isBlog}, false); // 不清空NotesByNotebookId缓存
		}
	});
}

// 设置notebook的blog状态
// 当修改notebook是否是blog时调用
Note.setAllNoteBlogStatus = function(notebookId, isBlog) {
	if(!notebookId) {
		return;
	}
	var notes = Note.getNotesByNotebookId(notebookId);
	if(!isArray(notes)) {
		return;
	}
	var len = notes.length;
	if(len == 0) {
		for(var i in Note.cache) {
			if(Note.cache[i].NotebookId == notebookId) {
				Note.cache[i].IsBlog = isBlog;
			}
		}
	} else {
		for(var i = 0; i < len; ++i) {
			notes[i].IsBlog = isBlog;
		}
	}
}

// 移动
Note.moveNote = function(target, data) {
	var noteId = $(target).attr("noteId");
	var note = Note.cache[noteId];
	var notebookId = data.notebookId;
	
	if(!note.IsTrash && note.NotebookId == notebookId) {
		return;
	}
	ajaxGet("/note/moveNote", {noteId: noteId, notebookId: notebookId}, function(ret) {
		if(ret && ret.NoteId) {
			if(note.IsTrash) {
				Note.changeToNext(target);
				$(target).remove();
				Note.clearCacheByNotebookId(notebookId);
			} else {
				// 不是trash, 移动, 那么判断是当前是否是all下
				// 不在all下, 就删除之
				// 如果当前是active, 那么clearNoteInfo之
				if(!Notebook.curActiveNotebookIsAll()) {
					Note.changeToNext(target);
					if($(target).hasClass("item-active")) {
						Note.clearNoteInfo();
					}
					$(target).remove();
				} else {
					// 不移动, 那么要改变其notebook title
					$(target).find(".note-notebook").html(Notebook.getNotebookTitle(notebookId));
				}
				
				// 重新清空cache 之前的和之后的
				Note.clearCacheByNotebookId(note.NotebookId);
				Note.clearCacheByNotebookId(notebookId);
			}
			
			// 改变缓存
			Note.setNoteCache(ret)
		}
	});
}

// 复制
// data是自动传来的, 是contextmenu数据 
Note.copyNote = function(target, data, isShared) {
	var noteId = $(target).attr("noteId");
	var note = Note.cache[noteId];
	var notebookId = data.notebookId;
	
	// trash不能复制, 不能复制给自己
	if(note.IsTrash || note.NotebookId == notebookId) {
		return;
	}
	
	var url = "/note/copyNote";
	var data = {noteId: noteId, notebookId: notebookId};
	if(isShared) {
		url = "/note/copySharedNote";
		data.fromUserId = note.UserId;
	}
	
	ajaxGet(url, data, function(ret) {
		if(ret && ret.NoteId) {
			// 重新清空cache 之后的
			Note.clearCacheByNotebookId(notebookId);
			// 改变缓存, 添加之
			Note.setNoteCache(ret)
		}
	});
}

// 这里速度不慢, 很快
Note.getContextNotebooks = function(notebooks) {
	var moves = [];
	var copys = [];
	var copys2 = [];
	for(var i in notebooks) {
		var notebook = notebooks[i];
		var move = {text: notebook.Title, notebookId: notebook.NotebookId, action: Note.moveNote}
		var copy = {text: notebook.Title, notebookId: notebook.NotebookId, action: Note.copyNote}
		var copy2 = {text: notebook.Title, notebookId: notebook.NotebookId, action: Share.copySharedNote}
		if(!isEmpty(notebook.Subs)) {
			var mc = Note.getContextNotebooks(notebook.Subs);
			move.items = mc[0];
			copy.items = mc[1];
			copy2.items = mc[2];
			move.type = "group";
			move.width = 150;
			copy.type = "group";
			copy.width = 150;
			copy2.type = "group";
			copy2.width = 150;
		}
		moves.push(move);
		copys.push(copy);
		copys2.push(copy2);
	}
	return [moves, copys, copys2];
}
// Notebook调用
Note.contextmenu = null;
Note.notebooksCopy = []; // share会用到
Note.initContextmenu = function() {
	var self = Note;
	if(Note.contextmenu) {
		Note.contextmenu.destroy();
	}
	// 得到可移动的notebook
	var notebooks = Notebook.everNotebooks;
	var mc = self.getContextNotebooks(notebooks);
	
	var notebooksMove = mc[0];
	var notebooksCopy = mc[1];
	self.notebooksCopy = mc[2];
	
	//---------------------
	// context menu
	//---------------------
	var noteListMenu = {
		width: 150, 
		items: [
			{ text: "分享给好友", alias: 'shareToFriends', icon: "", faIcon: "fa-share-square-o", action: Note.listNoteShareUserInfo},
			{ type: "splitLine" },
			{ text: "公开为博客", alias: 'set2Blog', icon: "", action: Note.setNote2Blog },
			{ text: "取消公开为博客", alias: 'unset2Blog', icon: "", action: Note.setNote2Blog },
			{ type: "splitLine" },
			// { text: "发送长微博", alias: 'html2Image', icon: "", action: Note.html2Image , width: 150, type: "group", items:[{text: "a"}]},
			// { type: "splitLine" },
			{ text: "删除", icon: "", faIcon: "fa-trash-o", action: Note.deleteNote },
			{ text: "移动", alias: "move", icon: "",
				type: "group", 
				width: 150, 
				items: notebooksMove
			},
			{ text: "复制", alias: "copy", icon: "",
				type: "group", 
				width: 150, 
				items: notebooksCopy
			}
		], 
		onShow: applyrule,
		onContextMenu: beforeContextMenu,
		
		parent: "#noteItemList",
		children: ".item-my",
	}
		
	function menuAction(target) {
		// $('#myModal').modal('show')
		showDialog("dialogUpdateNotebook", {title: "修改笔记本", postShow: function() {
		}});
	}
	function applyrule(menu) {
		var noteId = $(this).attr("noteId");
		var note = Note.cache[noteId];
		if(!note) {
			return;
		}
		// 要disable的items
		var items = [];
		
		// 如果是trash, 什么都不能做
		if(note.IsTrash) {
			items.push("shareToFriends");
			items.push("shareStatus");
			items.push("unset2Blog");
			items.push("set2Blog");
			items.push("copy");
		} else {
			// 是否已公开为blog
			if(!note.IsBlog) {
				items.push("unset2Blog");
			} else {
				items.push("set2Blog");
			}
			
			// 移动与复制不能是本notebook下
			var notebookTitle = Notebook.getNotebookTitle(note.NotebookId);
			items.push("move." + notebookTitle);
			items.push("copy." + notebookTitle);
		}

        menu.applyrule({
        	name: "target..",
            disable: true,
            items: items
        });		
	   
	}
	function beforeContextMenu() {
	    return this.id != "target3";
	}
	
	// 这里很慢!!
	Note.contextmenu = $("#noteItemList .item-my").contextmenu(noteListMenu);
}

//------------------- 事件
$(function() {
	//-----------------
	// for list nav
	$("#noteItemList").on("click", ".item", function(event) {
		event.stopPropagation();
		// 找到上级.item
		var parent = findParents(this, ".item");
		if(!parent) {
			return;
		}
		
		var noteId = parent.attr("noteId");
		if(!noteId) {
			return;
		}
		// 当前的和所选的是一个, 不改变
		if(Note.curNoteId == noteId) {
			return;
		}
		Note.changeNote(noteId);
	});
	
	//------------------
	// 新建笔记
	// 1. 直接点击新建 OR
	// 2. 点击nav for new note
	$("#newNoteBtn, #editorMask .note").click(function() {
		var notebookId = $("#curNotebookForNewNote").attr('notebookId');
		Note.newNote(notebookId);
	});
	$("#newNoteMarkdownBtn, #editorMask .markdown").click(function() {
		var notebookId = $("#curNotebookForNewNote").attr('notebookId');
		Note.newNote(notebookId, false, "", true);
	});
	$("#notebookNavForNewNote").on("click", "li div", function() {
		var notebookId = $(this).attr("notebookId");
		if($(this).hasClass("new-note-right")) {
			Note.newNote(notebookId, false, "", true);
		} else {
			Note.newNote(notebookId);
		}
	});
	$("#searchNotebookForAdd").click(function(e) {
		e.stopPropagation();
	});
	$("#searchNotebookForAdd").keyup(function() {
		var key = $(this).val();
		Notebook.searchNotebookForAddNote(key);
	});
	$("#searchNotebookForList").keyup(function() {
		var key = $(this).val();
		Notebook.searchNotebookForList(key);
	});
	
	//---------------------------
	// 搜索, 按enter才搜索
	/*
	$("#searchNoteInput").on("keyup", function(e) {
		Note.searchNote();
	});
	*/
	$("#searchNoteInput").on("keydown", function(e) {
		var theEvent = e; // window.event || arguments.callee.caller.arguments[0];
		if(theEvent.keyCode == 13 || theEvent.keyCode == 108) {
			theEvent.preventDefault();
			Note.searchNote();
			return false;
		}
	});
	
	//--------------------
	// Note.initContextmenu();
	
	//------------
	// 文档历史
	$("#contentHistory").click(function() {
		Note.listNoteContentHistories()
	});
	
	$("#saveBtn").click(function() {
		Note.curChangedSaveIt(true);
	});
	
	// blog
	$("#noteItemList").on("click", ".item-blog", function(e) {
		e.preventDefault();
		e.stopPropagation();
		// 得到ID
		var noteId = $(this).parent().attr('noteId');
		window.open("/blog/view/" + noteId);
	});
	
	// note setting
	$("#noteItemList").on("click", ".item-my .item-setting", function(e) {
		e.preventDefault();
		e.stopPropagation();
		var $p = $(this).parent();
		Note.contextmenu.showMenu(e, $p);
	});
});

// 定时器启动
Note.startInterval();