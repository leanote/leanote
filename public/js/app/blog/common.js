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
function goLogin(){ 
	var loginUrl = urlPrefix + '/login?from=' + encodeURI(location.href);
	location.href = loginUrl;
}
function goRegister() {
	var registerUrl = urlPrefix + '/register?from=' + encodeURI(location.href);
	location.href = registerUrl;
}
function needLogin() {
	if(typeof visitUserInfo == "undefined" || !visitUserInfo || !visitUserInfo.UserId) {
		// 弹框之
		var loginUrl = urlPrefix + '/login?from=' + encodeURI(location.href);
		var registerUrl = urlPrefix + '/register?from=' + encodeURI(location.href);
		var modal = BootstrapDialog.show({
	        title: "你还未登录",
	        message: '<div class="needLogin" style="border:none"><a href="' + loginUrl + '">立即登录</a>, 发表评论.<br />没有帐号? <a href="' + registerUrl +'">立即注册</a>',
	        nl2br: false
	   });
	   return true;
   }
   return false;
}
function scrollToTarget(t, fixed) {
	if(!fixed) {
		fixed = 0;
	}
	var $t = $(t)
	var targetOffset = $t.offset().top + fixed;
	$('html,body').animate({scrollTop: targetOffset}, 300);
}

var windowParam = 'width=700, height=580, top=180, left=320, toolbar=no, menubar=no, scrollbars=no, location=yes, resizable=no, status=no';
function getShareUrl(noteId) {
	return viewUrl + "/" + noteId;
}
function getShareTitle(title) {
	return encodeURI(title + " (来自leanote.com)");
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
	
//JavaScript函数：
var minute = 1000 * 60;
var hour = minute * 60;
var day = hour * 24;
var halfamonth = day * 15;
var month = day * 30;
// 2014-01-06T18:29:48.802+08:00

function goNowToDatetime(goNow) {
	if(!goNow) {
		return "";
	}
	return goNow.substr(0, 10) + " " + goNow.substr(11, 8);
}
function getDateDiff(dateTimeStamp){
	var now = new Date().getTime();
	var diffValue = now - dateTimeStamp;
	if(diffValue < 0){
		return "";
	}
	var monthC =diffValue/month;
	var weekC =diffValue/(7*day);
	var dayC =diffValue/day;
	var hourC =diffValue/hour;
	var minC =diffValue/minute;
	if(monthC>=1){
		 result=parseInt(monthC) + getMsg("monthsAgo");
	 }
	 else if(weekC>=1){
		 result=parseInt(weekC) + getMsg("weeksAgo");
	 }
	 else if(dayC>=1){
		 result=parseInt(dayC) + getMsg("daysAgo");
	 }
	 else if(hourC>=1){
		 result=parseInt(hourC) + getMsg("hoursAgo");
	 }
	 else if(minC>=1){
	 result=parseInt(minC) + getMsg("minutesAgo");
	 }else {
		 result=getMsg("justNow");
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