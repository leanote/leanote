// 返回是否是re.Ok == true
function reIsOk(re) {
	return re && typeof re == "object" && re.Ok;
}
function showAlert(id, msg, type, id2Focus) {
	$(id).html(msg).removeClass("alert-danger").removeClass("alert-success").removeClass("alert-warning").addClass("alert-" + type).show();
	if(id2Focus) {
		$(id2Focus).focus();
	}
}
function hideAlert(id, timeout) {
	if(timeout) {
		setTimeout(function() {
			$(id).hide();
		}, timeout);
	} else {
		$(id).hide();
	}
}
function ajaxGet(url, param, func) {
	$.get(url, param, func);
}

function ajaxPost(url, param, func) {
	$.post(url, param, func);
}
//------
// jsonp
function ajaxGetP(url, param, callback) {
	ajaxP("GET", url, param, callback)
}
function ajaxPostP(url, param, callback) {
	ajaxP("POST", url, param, callback)
}
function ajaxP(method, url, param, callback) {
	param = param || {};
	callback = callback || function() {};
	$.ajax({ 
        dataType: "jsonp",//跨域访问 dataType 必须是jsonp 类型。  
        url: url,  
        type: method, 
        data: param,
        jsonp: "callback",
		jsonpCallback: "jsonpCallback",
        success: callback
    });
}

//------------------
// 分享
var windowParam = 'width=700, height=580, top=180, left=320, toolbar=no, menubar=no, scrollbars=no, location=yes, resizable=no, status=no';
function getShareUrl() {
	return location.href;
}
function getShareTitle(title) {
	return encodeURI(title + " (from https://leanote.com)");
}
function shareSinaWeibo(noteId, title, pic) {
	var url = "http://service.weibo.com/share/share.php?title=" + getShareTitle(title) + "&url=" + getShareUrl(noteId);
	window.open(url, 'Share', windowParam);
}
function shareTencentWeibo(noteId, title, pic) {
	var _appkey = '801542571';
	var url = "http://share.v.t.qq.com/index.php?c=share&a=index&appkey=" + _appkey +"&title=" + getShareTitle(title) + "&url=" + getShareUrl(noteId) +"&pic=" + pic;
	window.open(url, 'Share', windowParam);
}
function shareQQ(noteId, title, pic) {
	var url = 'http://sns.qzone.qq.com/cgi-bin/qzshare/cgi_qzshare_onekey?url=' + getShareUrl(noteId) + '&title=' + title + '&pics=' + pic;
	window.open(url, 'Share', windowParam);
}
function shareRenRen(noteId, title, pic) {
	var url = 'http://widget.renren.com/dialog/share?resourceUrl=' + getShareUrl(noteId) +  '&srcUrl=' + getShareUrl(noteId) + '&title=' + getShareTitle(title) + '&pic=' + pic;
	window.open(url, 'Share', windowParam);
}

// https://twitter.com/intent/tweet?text=&pic=
function shareTwitter(noteId, title, pic) {
	var url = 'https://twitter.com/intent/tweet?text=' + getShareTitle(title) + '&pic=' + pic;
	window.open(url, 'Share', windowParam);
}
// http://www.facebook.com/sharer.php?u=<?php the_permalink();?>&t=<?php the_title(); ?>” 
function shareFacebook(noteId, title, pic) {
	var url = ' http://www.facebook.com/sharer.php?t=' + getShareTitle(title) + '&pic=' + pic;
	window.open(url, 'Share', windowParam);
}
	

// go的datetime转成datetime字符串
// 2014-01-06T18:29:48.802+08:00 => 2012-12-12 12:12:12 字符串
function goNowToDatetime(goNow) {
	if(!goNow) {
		return "";
	}
	return goNow.substr(0, 10) + " " + goNow.substr(11, 8);
}

// 转成Date()
// 2014-01-06T18:29:48.802+08:00 => new Date(year, month, day, hour, minute, second)
function goNowToDate(goNow) {
	if(!goNow) {
		return "";
	}
	var str = goNow.substr(0, 10) + " " + goNow.substr(11, 8);
	// 2012-12-12 12:12:12
	var tempStrs = str.split(" ");
	var dateStrs = tempStrs[0].split("-");
	var year = parseInt(dateStrs[0], 10);
	var month = parseInt(dateStrs[1], 10) - 1;
	var day = parseInt(dateStrs[2], 10);
	var timeStrs = tempStrs[1].split(":");
	var hour = parseInt(timeStrs [0], 10);
	var minute = parseInt(timeStrs[1], 10) - 1;
	var second = parseInt(timeStrs[2], 10);
	var date = new Date(year, month, day, hour, minute, second);
	return date;
}

// 距离现在有多久, 评论时间间隔
var diff = {
	minute : 1000 * 60,
	hour : 1000 * 60 * 60,
	day : 1000 * 60 * 60 * 24,
	halfamonth : 1000 * 60 * 60 * 24 * 15,
	month : 1000 * 60 * 60 * 24 * 30
}
function getDateDiff(dateTimeStamp) {
    var now = new Date().getTime();
    var diffValue = now - dateTimeStamp;
    if (diffValue < 0) {
        return "";
    }
    var monthC = diffValue / diff.month;
    var weekC = diffValue / (7 * diff.day);
    var dayC = diffValue / diff.day;
    var hourC = diffValue / diff.hour;
    var minC = parseInt(diffValue / diff.minute);
    if (monthC >= 1) {
        result = parseInt(monthC) + " month ago";
    } else if (weekC >= 1) {
        result = parseInt(weekC) + " weeks ago";
    } else if (dayC >= 1) {
        result = parseInt(dayC) + " days ago";
    } else if (hourC >= 1) {
        result = parseInt(hourC) + " hours ago";
    } else if (minC > 1) {
        result = minC + " minutes ago";
    } else {
        result = "Just now";
    }
    return result;
}

function weixin() {
	var local=window.location.href;
	var title = $.trim($(".title").text());
	var desc = $.trim($("#desc").text());
	var imgUrl = $("#content img").eq(0).attr('src');
	window.shareData = { 
	   "imgUrl": imgUrl, 
		"timeLineLink":local,
		"sendFriendLink": local,
		"weiboLink":local,
		"tTitle": title,
		"tContent": desc,
		"fTitle": title,
		"fContent": desc,
		"wContent": desc 
	};
	document.addEventListener('WeixinJSBridgeReady', function onBridgeReady() {
		// 发送给好友
		WeixinJSBridge.on('menu:share:appmessage', function (argv) {
			WeixinJSBridge.invoke('sendAppMessage', { 
				"img_url": window.shareData.imgUrl,
				"img_width": "200",
				"link": window.shareData.sendFriendLink,
				"desc": window.shareData.fContent,
				"title": window.shareData.fTitle
			}, function (res) {
				hs_guide('none');
				_report('send_msg', res.err_msg);
			})
		});
	
		// 分享到朋友圈
		WeixinJSBridge.on('menu:share:timeline', function (argv) {
			WeixinJSBridge.invoke('shareTimeline', {
				"img_url": window.shareData.imgUrl,
				"img_width": "200",
				"link": window.shareData.timeLineLink,
				"desc": window.shareData.tContent,
				"title": window.shareData.tTitle
			}, function (res) {
				hs_guide('none');
				_report('timeline', res.err_msg);
			});
		});
	
		// 分享到微博
		WeixinJSBridge.on('menu:share:weibo', function (argv) {
			WeixinJSBridge.invoke('shareWeibo', {
				"content": window.shareData.wContent,
				"url": window.shareData.weiboLink,
			}, function (res) {
				hs_guide('none');
				_report('weibo', res.err_msg);
			});
		});
	}, false);
}


function scrollTo(self, tagName, text) {
	var iframe = $("#content");
	var target = iframe.find(tagName + ":contains(" + text + ")");
	
	// 找到是第几个
	// 在nav是第几个
	var navs = $('#blogNavContent [data-a="' + tagName + '-' + encodeURI(text) + '"]');
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
		if(LEA.isMobile) {
			top -= 50;
		}
		var nowTop = $(document).scrollTop();
		// 用$("body").scrllTop(10)没反应 firefox下
		$('html,body').animate({scrollTop: top}, 200);
		return;
	}
}
function genNav() {
	var $con = $("#content");
	var html = $con.html();
	// 构造一棵树
	// {"h1-title":{h2-title:{}}}
	var tree = [];//[{title: "xx", children:[{}]}, {title:"xx2"}];
	var hs = $con.find("h1,h2,h3,h4,h5,h6").toArray();
	var titles = '<ul>';
	for(var i = 0; i < hs.length; ++i) {
		var text = $(hs[i]).text(); 
		var tagName = hs[i].tagName.toLowerCase();
		// scrollTo在page.js中定义
		titles += '<li class="nav-' + tagName + '"><a data-a="' + tagName + '-' + encodeURI(text)+'" onclick="scrollTo(this, \'' + tagName + '\', \'' + text + '\')">' + text + '</a></li>';
	}
	titles += "</ul>";
	$("#blogNavContent").html(titles);
	if(!hs.length) {
		$("#blogNavContent").html("无");
		return false;
	}
	return true;
}

function initNav() {
	var hasNav = genNav();
	if(!hasNav) {
		return;
	}
	
	var $title = $(".title");
	var titlePos = $title.offset();
	var top = titlePos.top + 10;// - $title.height();
	// 手机下不要与标题在同一高度
	if(LEA.isMobile){ 
		top += 30;
	}
	if(top < 0) {
		top = 10;
	}

	var left = $title.width() + titlePos.left - 100;
	$("#blogNav").css("top", top).css("left", left);
	$("#blogNav").show();
	
	$("#blogNavNav").click(function() {
		var $o = $("#blogNavContent");
		if($o.is(":hidden")) {
			$o.show();
		} else {
			$o.hide();
		}
	});
	
	var $d = $(document);
	function reNav() {
	    var vtop = $d.scrollTop();
	    if(vtop <= top) {
			$("#blogNav").css("top", top-vtop);
	    } else {
	    	// 差距很磊了
	    	if(LEA.isMobile) {
				$("#blogNav").css("top", 50);
			} else {
				$("#blogNav").css("top", 10);
			}
	    }
	}
	reNav();
	$(window).scroll(reNav);
}

//-----------
// 分享与评论

// 得到登录的链接
function goLogin(){ 
	var loginUrl = siteUrl + '/login?from=' + encodeURI(location.href);
	location.href = loginUrl;
}
// 得到注册链接
function goRegister() {
	var registerUrl = siteUrl + '/register?from=' + encodeURI(location.href);
	location.href = registerUrl;
}
// 是否需要登录, 评论时用到
function needLogin() {
	if(typeof visitUserInfo == "undefined" || !visitUserInfo || !visitUserInfo.UserId) {
		// 弹框之
		var loginUrl = siteUrl + '/login?from=' + encodeURI(location.href);
		var registerUrl = siteUrl + '/register?from=' + encodeURI(location.href);
		try {
			var modal = BootstrapDialog.show({
		        title: "Please sign in first",
		        message: '<div class="needLogin" style="border:none"><a href="' + loginUrl + '">Sign in</a> to to leave a comment.<br />No Leanote account? <a href="' + registerUrl +'">Sign up now</a>',
		        nl2br: false
		   });
	   } catch(e) {}
	   
	   return true;
   }
   return false;
}

// 定位评论
function scrollToTarget(t, fixed) {
	if(!fixed) {
		fixed = 0;
	}
	var $t = $(t)
	var targetOffset = $t.offset().top + fixed;
	$('html,body').animate({scrollTop: targetOffset}, 300);
}

// 得到访问者信息, 因为自定义域名的原因, 需要用jsonp来获取
function getCurVisitUserInfo() {
	
}
// 增加阅读次数
function incReadNum(noteId) {
	//if(!$.cookie(noteId)) {
	//	$.cookie(noteId, 1);
		ajaxGet(getCurHostUrl() + "/blog/incReadNum", {noteId: noteId});
	//}
}
function getCurHostUrl() {
	return "//" + location.host;
}
function getLeanoteUrl() {
	return siteUrl || "http://leanote.com";
}
// 得到博客统计信息
function getPostStat(noteId, callback) {
	ajaxGet(getCurHostUrl() + "/blog/getPostStat", {noteId: noteId}, callback);
}
// 得到赞
function getLikes(noteId, callback) {
	ajaxGetP(getLeanoteUrl() + "/blog/getLikes", {noteId: noteId}, callback);
}
// 得到赞和评论
function getLikesAndComments(noteId, callback) {
	ajaxGetP(getLeanoteUrl() + "/blog/getLikesAndComments", {noteId: noteId}, callback);
}
// 得到评论
function getComments(noteId, page, callback) {
	ajaxGetP(getLeanoteUrl() + "/blog/getComments", {noteId: noteId, page: page}, callback);
}
// 点赞
function likePost(noteId, callback) {
	ajaxPostP(getLeanoteUrl() + "/blog/likePost", {noteId: noteId}, callback)
}
// 提交评论
function commentPost(noteId, commentId, content, callback) {
	var data = {noteId: self.noteId, toCommentId: commentId, content: content};
	ajaxPostP(getLeanoteUrl() + "/blog/commentPost", data, callback);
}
// 删除评论
function deleteComment(noteId, commentId, callback) {
	ajaxPostP(getLeanoteUrl() + "/blog/deleteComment", {noteId: noteId, commentId: commentId}, callback);
}
// 点赞评论
function likeComment(commentId, callback) {
	ajaxPostP(getLeanoteUrl() + "/blog/likeComment", {commentId: commentId}, callback);
}

// 分享与评论结束
//------------

//------------
// 手机端
var LEA = {isMobile: false};
function isMobile() {
	var u = navigator.userAgent;
	LEA.isMobile = false;
	LEA.isMobile = /Mobile|Android|iPhone/i.test(u);
	if(!LEA.isMobile && $(document).width() <= 600){ 
		LEA.isMobile = true
	}
}
$(function() {
	isMobile();
});
