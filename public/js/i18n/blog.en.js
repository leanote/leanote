var MSG = {"a":"a","aboutMe":"About Me","author":"Author","baseInfoSet":"Base info","blogClass":"Category","blogDesc":"Description","blogLogo":"Logo","blogLogoTips":"Upload image to replace blog title","blogName":"Title","blogNav":"Blog nav","blogNavs":"Navs","blogSet":"Blog configuration","cancel":"Cancel","chooseComment":"Comment System","chooseReason":"请选择举报理由","comment":"Comment","commentSet":"Comment","comments":"Comments","community":"Community","confirm":"Confirm","confirmDeleteComment":"Are you sure?","createdTime":"Created at","daysAgo":"days ago","delete":"Delete","disqusHelp":"Please input your Disqus Id","domain":"Custom domain","domainSet":"Domain","elegant":"Elegant","error":"Error","fullBlog":"Full blog","home":"Home","hoursAgo":"hours ago","justNow":"Just now","latestPosts":"Latest posts","like":"Like","minutesAgo":"minutes ago","monthsAgo":"months ago","more":"More...","moreShare":"More","navFixed":"Nav fixed at left side","needHelp":"Need help?","next":"Next","noBlog":"No blog","noTag":"No tag","none":"None","openComment":"Open comment?","other":"Other","previous":"Previous","qqZone":"QQ Zone","quickLinks":"Quick links","renren":"Renren","reply":"Reply","report":"Report","reportBlog?":"举报该博客?","reportComment?":"举报该评论?","reportReason":"Reason","reportReason1":"不友善内容","reportReason2":"广告等垃圾信息","reportReason3":"违法违规内容","reportReason4":"不宜公开讨论的政治内容","reportSuccess":"举报成功, 我们处理后会通知作者, 感谢您的监督","saveSuccess":"Save success","scanQRCode":"Open weichat and scan the QR code","signIn":"Sign In","signUp":"Sign Up","sinaWeibo":"Weibo","subDomain":"Sub domain","submitComment":"Submit","tencentWeibo":"Tencent Weibo","theme":"Theme","themeSet":"Theme","unlike":"Unlike","updatedTime":"Updated at","viewers":"Viewers","weeksAgo":"weeks ago","weixin":"Weichat"};
function getMsg(key, data) {
	var msg = MSG[key]
	if(msg) {
		if(data) {
			if(!isArray(data)) {
				data = [data];
			}
			for(var i = 0; i < data.length; ++i) {
				msg = msg.replace("%s", data[i]);
			}
		}
		return msg;
	}
	return key;
}