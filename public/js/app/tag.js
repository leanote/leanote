// Tag

// 蓝色, 红色怎么存到数据库中? 直接存蓝色

Tag.classes = {
	"蓝色": "label label-blue",
	"红色": "label label-red",
	"绿色": "label label-green",
	"黄色": "label label-yellow",
	"blue": "label label-blue",
	"red": "label label-red",
	"green": "label label-green",
	"yellow": "label label-yellow"
}

// 数据库中统一存En
Tag.mapCn2En = {
	"蓝色": "blue",
	"红色": "red",
	"绿色": "green",
	"黄色": "yellow",

}
Tag.mapEn2Cn = {
	"blue": "蓝色",
	"red": "红色",
	"green": "绿色",
	"yellow": "黄色",
}

Tag.t = $("#tags");

// called by Note
Tag.getTags = function() {
	var tags = [];
	Tag.t.children().each(function(){
		var text = $(this).text();
		text = text.substring(0, text.length - 1); // 把X去掉
		text = Tag.mapCn2En[text] || text;
		tags.push(text);
	});
	// 需要去重吗? 正常情况下不会重复
	return tags;
}

// called by Note
Tag.clearTags = function() {
	Tag.t.html("");
}

// 设置tags
// called by Note
Tag.renderTags = function(tags) {
	Tag.t.html("");
	if(isEmpty(tags)) {
		return;
	}
	// TODO 重构, 这样不高效
	for(var i = 0; i < tags.length; ++i) {
		var tag = tags[i];
		Tag.appendTag(tag);
	}
}

// tag最初状态
function revertTagStatus() {
	$("#addTagTrigger").show();
	$("#addTagInput").hide();
	// hideTagList();
}

function hideTagList(event) {
	$("#tagDropdown").removeClass("open");
	if (event) {
		event.stopPropagation()
	}
}
function showTagList(event) {
	$("#tagDropdown").addClass("open");
	if (event) {
		event.stopPropagation()
	}
}

// 只读模式下显示tags
// called by Note
Tag.renderReadOnlyTags = function(tags) {
	// 先清空
	$("#noteReadTags").html("");
	if(isEmpty(tags)) {
		$("#noteReadTags").html("无标签");
	}
	
	var i = true;
	function getNextDefaultClasses() {
		if (i) {
			return "label label-default";
			i = false
		} else {
			i = true;
			return "label label-info";
		}
	}
	
	for(var i in tags) {
		var text = tags[i];
		text = Tag.mapEn2Cn[text] || text;
		var classes = Tag.classes[text];
		if(!classes) {
			classes = getNextDefaultClasses();
		}
		tag = tt('<span class="?">?</span>', classes, text);
		
		$("#noteReadTags").append(tag);
	}
}

// 添加tag
// tag = {classes:"label label-red", text:"红色"}
// tag = life
Tag.appendTag = function(tag) {
	var isColor = false;
	var classes, text;
	
	if (typeof tag == "object") {
		classes = tag.classes;
		text = tag.text;
		if(!text) {
			return;
		}
	} else {
		tag = $.trim(tag);
		text = tag;
		if(!text) {
			return;
		}
		var classes = Tag.classes[text];
		if(classes) {
			isColor = true;
		} else {
			classes = "label label-default";
		}
	}
	
	text = Tag.mapEn2Cn[text] || text;
	tag = tt('<span class="?">?<i title="删除">X</i></span>', classes, text);

	// 避免重复
	$("#tags").children().each(function() {
		if (isColor) {
			var tagHtml = $("<div></div>").append($(this).clone()).html();
			if (tagHtml == tag) {
				$(this).remove();
			}
		} else if (text + "X" == $(this).text()) {
			$(this).remove();
		}
	});

	$("#tags").append(tag);

	hideTagList();

	if (!isColor) {
		reRenderTags();
	}
}

// 为了颜色间隔, add, delete时调用
function reRenderTags() {
	var defautClasses = [ "label label-default", "label label-info" ];
	var i = 0;
	$("#tags").children().each(
		function() {
			var thisClasses = $(this).attr("class");
			if (thisClasses == "label label-default"
					|| thisClasses == "label label-info") {
				$(this).removeClass(thisClasses).addClass(
						defautClasses[i % 2]);
				i++;
			}
		});
}

//-----------
// 左侧nav en -> cn
Tag.renderTagNav = function(tags) {
	tags = tags || [];
	for(var i in tags) {
		var tag = tags[i];
		if(tag == "red" || tag == "blue" || tag == "yellow" || tag == "green") {
			continue;
		}
		var text = Tag.mapEn2Cn[tag] || tag;
		var classes = Tag.classes[tag] || "label label-default";
		$("#tagNav").append(tt('<li><a> <span class="?">?</span></li>', classes, text));
	}
}


// 事件
$(function() {
	// tag
	$("#addTagTrigger").click(function() {
		$(this).hide();
		$("#addTagInput").show().focus().val("");
	});
	
	$("#addTagInput").click(function(event) {
		showTagList(event);
	});
	
	$("#addTagInput").blur(function() {
		var val = $(this).val();
		if(val) {
			Tag.appendTag(val, true);
		}
		return;
		// 下面不能有, 有就有问题
		$("#addTagTrigger").show();
		$("#addTagInput").hide();
		// revertTagStatus();
	});
	$('#addTagInput').keydown(function(e) {
		if (e.keyCode == 13) {
			hideTagList();
			// 如果有值, 再生成, 没值直接隐藏
			if ($("#addTagInput").val()) {
				$(this).trigger("blur");
				$("#addTagTrigger").trigger("click");
			} else {
				$(this).trigger("blur");
			}
		}
	});
	// 点击下拉时也会触发input的blur事件
	$("#tagColor li").click(function(event) {
		var a;
		if($(this).attr("role")) {
			a = $(this).find("span");
		} else {
			a = $(this);
		}
		Tag.appendTag({
			classes : a.attr("class"),
			text : a.text()
		});
	});
	// 这是个问题, 为什么? 捕获不了事件?, input的blur造成
	/*
	$(".label").click(function(event) {
		var a = $(this);
		Tag.appendTag({
			classes : a.attr("class"),
			text : a.text()
		});
		// event.stopPropagation();
	});
	*/
	
	$("#tags").on("click", "i", function() {
		$(this).parent().remove();
		reRenderTags();
	});
	
	//-------------
	// nav 标签搜索
	function searchTag() {
		var tag = $.trim($(this).text());
		tag = Tag.mapCn2En[tag] || tag;
		
		// 学习changeNotebook
		
		// 1
		Note.curChangedSaveIt();
		
		// 2 先清空所有
		// 也会把curNoteId清空
		Note.clearAll();
		
		$("#tagSearch").html($(this).html()).show();
		showLoading();
		ajaxGet("/note/searchNoteByTags", {tags: [tag]}, function(notes) {
			hideLoading();
			if(notes) {
				// 和note搜索一样
				// 设空, 防止发生上述情况
				// Note.curNoteId = "";
				
				Note.renderNotes(notes);
				if(!isEmpty(notes)) {
					Note.changeNote(notes[0].NoteId);
				}
			}
		});
	}
	$("#myTag .folderBody").on("click", "li", searchTag);
	$("#minTagNav").on("click", "li", searchTag);
});