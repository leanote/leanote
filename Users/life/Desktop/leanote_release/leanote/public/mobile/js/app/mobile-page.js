//------------------
// mobile
//------------------
// 展示哪个notebooks, note, view, editor
LEA.curStatus = "notebooks";
LEA.preStatus = [];
LEA.toggles = ["notebooks", "notes", "view"]
function toggle(which, title, isBack) {
	if(!which || which == LEA.curStatus) {
		// 这个时候, 可以调出slider
		if(LEA.snapper.state().state=="left" ){
			LEA.snapper.close();
		} else {
			LEA.snapper.open('left');
		}
		return;
	}
	
	// isBack, 不要加入历史中
	if(!isBack) {
		LEA.preStatus.push(LEA.curStatus);
	}
	LEA.curStatus = which;
	
	for(var i in LEA.toggles) {
		var w = LEA.toggles[i];
		if(w != which && !$("#" + w).is(":hidden")) {
//			$("#" + w).fadeOut();
			$("#" + w).hide();
		}
	}
//	$("#" + which).show(100);
	$("#" + which).fadeIn();
	
	// 设置标题
	if(title) {
	    var maxTitleWidth = $("#" + which + " .content-controls").width()-50-$("#"+ which + " .btns").width();
	    $("#" + which + " .g-title").html(title).width(maxTitleWidth - 10);
	}
}// 当前状态是?

// notebooks -> notes -> view -> edit
function back() {
	toggle(LEA.preStatus.pop(), "", true);
}

$(function() {
	$(".back").click(function() {
		back();
	});
	/*
	window.onbeforeunload = function(e) {
		back();
		e.preventDefault();
	};
	*/
	$("#status").fadeOut(); // will first fade out the loading animation
	$("#preloader").delay(350).fadeOut("slow"); // will fade out the white DIV that covers the website.
	
	$('.close-nav, .close-sidebar, .sidebar-close').click(function(){
		snapper.close();
		return false;
	});
		
	var snapper = new Snap({
	  element: document.getElementById('content')
	});
	LEA.snapper = snapper;

	/*
	$('.deploy-sidebar').click(function(){
		if( snapper.state().state=="left" ){
			snapper.close();
		} else {
			snapper.open('left');
		}
		return false;
	});
	*/
	
	$(".nav-newest").click(function() {
		Notebook.ChangeNotebook("0");
		snapper.close();
	});
	$(".nav-myNotebooks").click(function() {
		toggle("notebooks");
		snapper.close();
	});
	$(".nav-logout").click(function() {
		location.href="/mobile/logout";
	});
})